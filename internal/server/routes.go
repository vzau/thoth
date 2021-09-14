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

package server

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/vzau/thoth/internal/controllers/auth"
	"github.com/vzau/thoth/internal/controllers/live"
	"github.com/vzau/thoth/internal/server/middleware"
	"github.com/vzau/thoth/internal/server/response"
)

func SetupRoutes(engine *gin.Engine) {
	v1 := engine.Group("/v1")
	{
		liveGroup := v1.Group("/live")
		{
			liveGroup.GET("/flights/:fac", live.GetFlights)
			liveGroup.GET("/controllers/:fac", live.GetControllers)
		}

		authGroup := v1.Group("/auth")
		{
			authGroup.GET("/login", auth.GetLogin)
			authGroup.GET("/callback", auth.GetCallback)
			authGroup.GET("/logout", auth.GetLogout)
		}

		authorized := v1.Group("/")
		authorized.Use(middleware.NotGuest)
		{
			authorized.GET("/auth/info", auth.GetInfo)
		}
	}

	engine.GET("/test", func(c *gin.Context) {
		s := sessions.Default(c)
		response.Respond(c, http.StatusOK, gin.H{"user": s.Get("user"), "cid": s.Get("cid")})
	})

	engine.GET("/ping", func(c *gin.Context) {
		response.Respond(c, http.StatusOK, struct{ Message string }{"PONG"})
	})
}
