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

package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vzau/thoth/internal/server/response"
	"github.com/vzau/thoth/pkg/database"
	"github.com/vzau/thoth/pkg/user"
	dbTypes "github.com/vzau/types/database"
)

func GetRoles(c *gin.Context) {
	cid, err := strconv.ParseUint(c.Param("cid"), 10, 64)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	roles, err := user.GetRoles(cid)
	if err != nil {
		log.Error("Error fetching user %d roles: %s", cid, err.Error())
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, struct {
		Roles []*dbTypes.Role `json:"roles"`
	}{Roles: roles})
}

func PostRole(c *gin.Context) {
	cid, err := strconv.ParseUint(c.Param("cid"), 10, 64)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	var role dbTypes.Role
	if user.HasRoles(cid, []string{c.Param("role")}) {
		response.RespondError(c, http.StatusConflict, "Conflict")
		return
	}

	if err := database.DB.Where("name = ?", c.Param("role")).First(&role).Error; err != nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	if err := database.DB.Exec("INSERT INTO user_roles(user_c_id, role_id) VALUES(?, ?)", cid, role.ID).Error; err != nil {
		log.Error("Could not add to user_roles: %d %d: %s", cid, role.ID, err.Error())
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.RespondBlank(c, http.StatusCreated)
}

func DeleteRole(c *gin.Context) {
	cid, err := strconv.ParseUint(c.Param("cid"), 10, 64)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	if !user.HasRoles(cid, []string{c.Param("role")}) {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	if err := database.DB.Exec("DELETE FROM user_roles WHERE user_c_id = ? AND role_id = (SELECT id FROM roles WHERE name = ?)", cid, c.Param("role")).Error; err != nil {
		log.Error("Error deleting role %s from %d: %s", c.Param("role"), cid, err.Error())
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.RespondBlank(c, http.StatusNoContent)
}
