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
