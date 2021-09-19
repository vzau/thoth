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

package cdn

import (
	"encoding/json"
	"time"

	dbTypes "github.com/vzau/types/database"
)

type CDN struct {
	Name        string `json:"name" form:"name" yaml:"name"`
	Category    string `json:"category" form:"category" yaml:"category"`
	Description string `json:"description" form:"description" yaml:"description"`
}

type FileDTO struct {
	Name        string    `json:"name" form:"name" yaml:"name"`
	Category    string    `json:"category" form:"category" yaml:"category"`
	Description string    `json:"description" form:"description" yaml:"description"`
	URL         string    `json:"url" form:"url" yaml:"url"`
	CreatedAt   time.Time `json:"created_at" form:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" form:"updated_at" yaml:"updated_at"`
}

type FileCache struct {
	File dbTypes.File `json:"file" form:"file" yaml:"file"`
	Data []byte       `json:"data" form:"data" yaml:"data"`
}

func (f FileCache) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *FileCache) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, f)
}
