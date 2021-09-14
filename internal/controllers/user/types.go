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

type UserPatch struct {
	OperatingInitials string `json:"operatingInitials" form:"operatingInitials" yaml:"operatingInitials" binding:"max=2"`
	Status            string `json:"status" form:"status" yaml:"status" binding:"oneof='' none active inactive leave"`
	Delivery          string `json:"delivery" form:"delivery" yaml:"delivery" binding:"oneof='' none minor major-solo certified"`
	Ground            string `json:"ground" form:"ground" yaml:"ground" binding:"oneof='' none minor major-solo certified"`
	Local             string `json:"local" form:"local" yaml:"local" binding:"oneof='' none minor major-solo certified"`
	Approach          string `json:"approach" form:"approach" yaml:"approach" binding:"oneof='' none minor major-solo certified"`
	Enroute           string `json:"enroute" form:"enroute" yaml:"enroute" binding:"oneof='' none major-solo certified"`
}

type UserPatchTraining struct {
	Delivery string `json:"delivery" form:"delivery" yaml:"delivery" binding:"oneof=none minor major-solo certified"`
	Ground   string `json:"ground" form:"ground" yaml:"ground" binding:"oneof=none minor major-solo certified"`
	Local    string `json:"local" form:"local" yaml:"local" binding:"oneof=none minor major-solo certified"`
	Approach string `json:"approach" form:"approach" yaml:"approach" binding:"oneof=none minor major-solo certified"`
	Enroute  string `json:"enroute" form:"enroute" yaml:"enroute" binding:"oneof=none major-solo certified"`
}
