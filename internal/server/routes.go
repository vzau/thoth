package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vzau/thoth/internal/controllers/auth"
	"github.com/vzau/thoth/internal/controllers/live"
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
			//authGroup.GET("/refresh", )
			//authGroup.GET("/logout",)
		}
	}

	engine.GET("/test", func(c *gin.Context) {
		response.HandleError(c, "Error Test")
	})

	engine.GET("/ping", func(c *gin.Context) {
		response.Respond(c, http.StatusOK, struct{ Message string }{"PONG"})
	})
}
