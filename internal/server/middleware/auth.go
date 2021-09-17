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

package middleware

import (
	"github.com/dhawton/log4g"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/vzau/thoth/internal/server/response"
	"github.com/vzau/thoth/pkg/database"
	dbTypes "github.com/vzau/types/database"
	"gorm.io/gorm/clause"
)

var log = log4g.Category("middleware/auth")

func Auth(c *gin.Context) {
	session := sessions.Default(c)
	cidSession := session.Get("cid")

	if cidSession == nil {
		c.Set("x-cid", 0)
		c.Set("x-user", nil)
		c.Next()
		return
	}

	cid := cidSession.(uint)

	user := &dbTypes.User{}
	if err := database.DB.Where(&dbTypes.User{CID: cid}).Preload(clause.Associations).First(&user).Error; err != nil {
		log.Warning("User not found: %d", cid)
		c.Set("x-cid", 0)
		c.Set("x-user", &user)
		response.RespondError(c, 401, "Unauthorized")
		c.Abort()
		return
	}

	c.Set("x-cid", user.CID)
	c.Set("x-user", user)
	c.Next()
}

func NotGuest(c *gin.Context) {
	if c.GetUint("x-cid") == 0 {
		response.RespondError(c, 401, "Unauthorized")
		c.Abort()
		return
	}

	c.Next()
}

func HasRoles(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		NotGuest(c)

		user := c.MustGet("x-user").(*dbTypes.User)
		for _, v := range roles {
			if userHasRole(user, v) {
				c.Next()
				return
			}
		}
		response.RespondError(c, 403, "Forbidden")
		c.Abort()
	}
}

func HasRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		NotGuest(c)

		user := c.MustGet("x-user").(*dbTypes.User)

		if !userHasRole(user, role) {
			response.RespondError(c, 403, "Forbidden")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func userHasRole(user *dbTypes.User, role string) bool {
	for _, v := range user.Roles {
		if v.Name == role {
			return true
		}
	}

	return false
}

func IsStaff(c *gin.Context) {
	NotGuest(c)

	user := c.MustGet("x-user").(*dbTypes.User)
	staffRoles := []string{"ATM", "DATM", "TA", "EC", "FE", "WM"}

	for _, v := range staffRoles {
		if userHasRole(user, v) {
			c.Next()
			return
		}
	}

	response.RespondError(c, 403, "Forbidden")
	c.Abort()
}
