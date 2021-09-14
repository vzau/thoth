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
	"time"

	"github.com/dhawton/log4g"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Make change that updates cookie value (obfuscation), and also updates maxage
func UpdateCookie(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("u", time.Now().String())
	err := session.Save()
	if err != nil {
		log4g.Category("middleware/cookie").Error("Error saving cookie %s", err.Error())
	}
	c.Next()
}
