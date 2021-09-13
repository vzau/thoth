/*
ZAU Thoth API
Copyright (C) 2021 Daniel A. Hawton (daniel@hawton.org)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dhawton/log4g"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/vzau/common/utils"
	"github.com/vzau/thoth/internal/server/response"
	dbTypes "github.com/vzau/types/database"
	"golang.org/x/oauth2"
)

var oauthConfig *oauth2.Config
var log = log4g.Category("controllers/oauth")

type UserResult struct {
	UserData dbTypes.User `json:"user"`
	Message  string       `json:"message"`
	err      error
}

type RoleResult struct {
	Roles []string
	err   error
}

func Init() {
	oauthConfig = &oauth2.Config{
		ClientID:     utils.Getenv("OAUTH_CLIENT_ID", ""),
		ClientSecret: utils.Getenv("OAUTH_CLIENT_SECRET", ""),
		RedirectURL:  utils.Getenv("OAUTH_REDIRECT_URI", ""),
		Scopes:       []string{"all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("%s/oauth/authorize", utils.Getenv("OAUTH_SSO_BASE_URL", "https://auth.chicagoartcc.org")),
			TokenURL:  fmt.Sprintf("%s/oauth/token", utils.Getenv("OAUTH_SSO_BASE_URL", "https://auth.chicagoartcc.org")),
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}

func GetLogin(c *gin.Context) {
	session := sessions.Default(c)
	code, _ := gonanoid.Generate("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", 32)
	session.Set("state", code)
	session.Save()

	url := oauthConfig.AuthCodeURL(code)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GetCallback(c *gin.Context) {
	session := sessions.Default(c)
	state := session.Get("state").(string)
	session.Delete("state")
	session.Save()

	if queryState, _ := c.GetQuery("state"); queryState != state {
		log.Error("Received bad state at callback")
		response.HandleError(c, "Invalid State")
	}

	code, _ := c.GetQuery("code")
	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Error("Error exchanging code for token", err)
		response.HandleError(c, "Error exchanging code for token")
	}

	userchan := make(chan UserResult)
	go GetUserData(token.AccessToken, userchan)
	userResult := <-userchan

	if userResult.err != nil {
		log.Error("Error getting user data", userResult.err)
		response.HandleError(c, "Error getting user data")
	}

	session.Set("user", userResult.UserData)
	session.Set("cid", userResult.UserData.CID)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, utils.Getenv("LOGIN_REDIRECT", "https://www.chicagoartcc.org"))
}

func GetUserData(token string, result chan UserResult) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/info", utils.Getenv("OAUTH_SSO_BASE_URL", "https://auth.chicagoartcc.org")), bytes.NewBuffer(nil))
	if err != nil {
		log.Error("Error creating new request", err)
		result <- UserResult{err: err}
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for key, val := range via[0].Header {
			req.Header[key] = val
		}

		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error getting user data", err)
		result <- UserResult{err: err}
		return
	}
	defer resp.Body.Close()

	userdata := UserResult{}
	data, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(data, &userdata); err != nil {
		log.Error("Error unmarshalling user data (%s)", err, string(data))
		result <- UserResult{err: err}
		return
	}
	result <- userdata
}
