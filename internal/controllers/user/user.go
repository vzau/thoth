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

	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	"github.com/vzau/common/utils"
	"github.com/vzau/thoth/internal/server/response"
	"github.com/vzau/thoth/internal/vatusa"
	"github.com/vzau/thoth/pkg/database"
	"github.com/vzau/thoth/pkg/discord"
	userPkg "github.com/vzau/thoth/pkg/user"
	dbTypes "github.com/vzau/types/database"
	yaml "gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var log = log4g.Category("controllers/user")

func GetUsers(c *gin.Context) {
	users := []dbTypes.User{}
	if err := database.DB.Preload(clause.Associations).Find(&users).Error; err != nil {
		log.Error("Error fetching users: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error fetching users")
		return
	}

	response.Respond(c, http.StatusOK, struct {
		Users []dbTypes.User `json:"users"`
	}{users})
}

func GetUser(c *gin.Context) {
	cid, err := strconv.ParseUint(c.Param("cid"), 10, 32)
	if err != nil {
		response.RespondMessage(c, http.StatusBadRequest, "Invalid user id")
		return
	}

	user := findUserOrAbort(c, cid)

	response.Respond(c, http.StatusOK, struct {
		User dbTypes.User `json:"user"`
	}{User: *user})
}

func PatchUser(c *gin.Context) {
	cid, err := strconv.ParseUint(c.Param("cid"), 10, 32)
	if err != nil {
		response.RespondMessage(c, http.StatusBadRequest, "Invalid user id")
		return
	}

	user := findUserOrAbort(c, cid)
	data := UserPatch{}
	if err := c.ShouldBind(&data); err != nil {
		response.RespondMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	changer := c.MustGet("x-user").(*dbTypes.User)
	patchNotes, _ := yaml.Marshal(data)

	if !userPkg.HasRoles(cid, []string{"ATM", "DATM", "WM"}) {
		if data.OperatingInitials != "" || data.Status != "" {
			go discord.SendToDiscordf(utils.Getenv("DISCORD_WEBHOOK_AUDIT", ""), "%s %s attempted to edit %s %s's profile, but cannot change OI or Status so was denied:\n%s", changer.FirstName, changer.LastName, user.FirstName, user.LastName, string(patchNotes))
			response.RespondMessage(c, http.StatusForbidden, "You do not have permission to edit those field(s)")
			return
		}
	}

	// Remove from VATUSA, but only in prod env
	// Run in goroutine as VATUSA's API is slow
	go func() {
		if data.Status != user.Status && data.Status == "removed" {
			if user.ControllerType == "home" {
				err := vatusa.DeleteController(cid)
				if err != nil {
					log.Error("Error deleting controller: %s", err.Error())
					discord.SendToDiscordf(utils.Getenv("DISCORD_WEBHOOK_AUDIT", ""), "%s %s attempted to remove %s %s from visitng roster, but failed to remove from VATUSA: %s", changer.FirstName, changer.LastName, user.FirstName, user.LastName, err.Error())
					return
				}
				discord.SendToDiscordf(utils.Getenv("DISCORD_WEBHOOK_AUDIT", ""), "%s %s removed %s %s (%d) from home roster.", changer.FirstName, changer.LastName, user.FirstName, user.LastName, cid)
			} else if user.ControllerType == "visitor" {
				err := vatusa.DeleteVisitor(cid)
				if err != nil {
					log.Error("Error deleting visitor: %s", err.Error())
					discord.SendToDiscordf(utils.Getenv("DISCORD_WEBHOOK_AUDIT", ""), "%s %s attempted to remove %s %s from roster, but failed to remove from VATUSA: %s", changer.FirstName, changer.LastName, user.FirstName, user.LastName, err.Error())
					return
				}
				discord.SendToDiscordf(utils.Getenv("DISCORD_WEBHOOK_AUDIT", ""), "%s %s removed %s %s (%d) from visitor roster.", changer.FirstName, changer.LastName, user.FirstName, user.LastName, cid)
			}
		}
	}()

	if data.OperatingInitials != "" {
		user.OperatingInitials = data.OperatingInitials
	}
	if data.Status != "" {
		user.Status = data.Status
	}
	if data.Delivery != "" {
		user.Delivery = data.Delivery
	}
	if data.Ground != "" {
		user.Ground = data.Ground
	}
	if data.Local != "" {
		user.Local = data.Local
	}
	if data.Approach != "" {
		user.Approach = data.Approach
	}
	if data.Enroute != "" {
		user.Enroute = data.Enroute
	}

	go discord.SendToDiscordf(utils.Getenv("DISCORD_WEBHOOK_AUDIT", ""), "%s %s updated %s %s's profile:\n%s", changer.FirstName, changer.LastName, user.FirstName, user.LastName, string(patchNotes))

	if err = database.DB.Save(user).Error; err != nil {
		log.Error("Error updating user: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error updating user")
		return
	}

	response.Respond(c, http.StatusOK, struct {
		User dbTypes.User `json:"user"`
	}{User: *user})
}

func DeleteUser(c *gin.Context) {
	cid, err := strconv.ParseUint(c.Param("cid"), 10, 32)
	if err != nil {
		response.RespondMessage(c, http.StatusBadRequest, "Invalid user id")
		return
	}

	user := findUserOrAbort(c, cid)
	if user.Status != "removed" {
		response.RespondMessage(c, http.StatusUnprocessableEntity, "You cannot delete a user that is not removed")
		return
	}
	if err = database.DB.Delete(user).Error; err != nil {
		log.Error("Error deleting user: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error deleting user")
		return
	}

	response.RespondBlank(c, http.StatusNoContent)
}

func findUserOrAbort(c *gin.Context, cid uint64) *dbTypes.User {
	user := dbTypes.User{}
	if err := database.DB.Preload(clause.Associations).First(&user, cid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.RespondMessage(c, http.StatusNotFound, "User not found")
			c.Abort()
			return nil
		}
		log.Error("Error fetching user: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error fetching user")
		c.Abort()
		return nil
	}

	return &user
}
