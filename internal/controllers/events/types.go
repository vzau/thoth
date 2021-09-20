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

package events

import (
	"time"

	dbTypes "github.com/vzau/types/database"
)

type EventDTO struct {
	ID          uint                    `json:"id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Start       time.Time               `json:"start"`
	End         time.Time               `json:"end"`
	Banner      string                  `json:"banner"`
	Positions   []dbTypes.EventPosition `json:"positions"`
	SignUps     []dbTypes.EventSignUp   `json:"signups"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}
