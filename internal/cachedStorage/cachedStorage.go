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

package cachedStorage

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/vzau/thoth/pkg/cache"
	"github.com/vzau/thoth/pkg/storage"
	dbTypes "github.com/vzau/types/database"
)

func GetCachedFile(file dbTypes.File) (*bytes.Reader, error) {
	cacheData, exists := cache.Cache.Get(fmt.Sprintf("cdn/%d", file.ID))
	if !exists {
		data, err := storage.GetFileStream(file.Filename)
		if err != nil {
			return nil, err
		}

		defer data.Close()

		buff, err := ioutil.ReadAll(data)
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(buff)
		cache.Cache.Set(fmt.Sprintf("cdn/%d", file.ID), FileCache{
			File: file,
			Data: buff,
		}, 0)
		return reader, nil
	} else {
		data := FileCache{}
		err := data.UnmarshalBinary([]byte(cacheData.(string)))
		if err != nil {
			cache.Cache.Delete(fmt.Sprintf("cdn/%d", file.ID))
			return GetCachedFile(file)
		}

		return bytes.NewReader(data.Data), nil
	}
}
