package live

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vzau/thoth/internal/server/response"
	"github.com/vzau/thoth/pkg/database"
	dbTypes "github.com/vzau/types/database"
)

func GetFlights(c *gin.Context) {
	var flights []dbTypes.Flights
	database.DB.Where("facility = ?", c.Param("fac")).Find(&flights)

	response.Respond(c, http.StatusOK, struct {
		Flights []dbTypes.Flights `json:"flights"`
	}{flights})
}

func GetControllers(c *gin.Context) {
	var controllers []dbTypes.OnlineControllers
	database.DB.Where("facility = ?", c.Param("fac")).Find(&controllers)

	response.Respond(c, http.StatusOK, struct {
		Controllers []dbTypes.OnlineControllers `json:"controllers"`
	}{controllers})
}
