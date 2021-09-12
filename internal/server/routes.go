package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vzau/thoth/internal/controllers/live"
)

func SetupRoutes(engine *gin.Engine) {
	v1 := engine.Group("/v1")
	{
		liveGroup := v1.Group("/live")
		{
			liveGroup.GET("/flights/:fac", live.GetFlights)
			liveGroup.GET("/controllers/:fac", live.GetControllers)
		}
	}

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "PONG"})
	})
}
