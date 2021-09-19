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
	"github.com/vzau/thoth/pkg/database"
	dbTypes "github.com/vzau/types/database"
	"gorm.io/gorm/clause"
)

var StaffRoles = []string{"ATM", "DATM", "TA", "EC", "FE", "WM"}
var StaffRolesTraining = []string{"ATM", "DATM", "TA", "EC", "FE", "WM", "INS", "MTR"}

func GetUser(cid uint64) (*dbTypes.User, error) {
	user := &dbTypes.User{}
	if err := database.DB.Where("cid = ?", cid).Preload(clause.Associations).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetRoles(cid uint64) ([]*dbTypes.Role, error) {
	user, err := GetUser(cid)
	if err != nil {
		return nil, err
	}

	return user.Roles, nil
}

func HasRoles(cid uint64, requiredRoles []string) bool {
	roles, err := GetRoles(cid)
	if err != nil {
		return false
	}

	for idx := range roles {
		for _, requiredRole := range requiredRoles {
			if roles[idx].Name == requiredRole {
				return true
			}
		}
	}

	return false
}

func HasRolesWithUser(user *dbTypes.User, requiredRoles []string) bool {
	for idx := range user.Roles {
		for _, requiredRole := range requiredRoles {
			if user.Roles[idx].Name == requiredRole {
				return true
			}
		}
	}
	return false
}
