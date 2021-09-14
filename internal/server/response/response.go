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

package response

import (
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
)

type R struct {
	XMLName xml.Name    `xml:"response" json:"-" yaml:"-"`
	Status  string      `xml:"status" json:"status"`
	Data    interface{} `xml:"data" json:"data"`
}

func RespondMessage(c *gin.Context, status int, message string) {
	Respond(c, status, struct {
		Message string `json:"message"`
	}{message})
}

func RespondError(c *gin.Context, status int, message string) {
	Respond(c, status, struct {
		Message string `json:"message"`
	}{message})
}

func Respond(c *gin.Context, status int, data interface{}) {
	ret := R{}
	ret.Status = http.StatusText(status)
	ret.Data = data

	// Use this to allow client to specify what format, but default to JSON
	if c.GetHeader("Accept") == "text/x-yaml" || c.GetHeader("Accept") == "application/x-yaml" || c.GetHeader("Accept") == "application/yaml" {
		c.YAML(status, ret)
	} else if c.GetHeader("Accept") == "application/xml" {
		c.XML(status, ret)
	} else {
		c.JSON(status, ret)
	}
}

func HandleError(c *gin.Context, message string) {
	c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": message})
	c.Abort()
}
