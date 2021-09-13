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

package jwk

import (
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/vzau/thoth/pkg/http"
)

func FetchKeySet(url string) (*jwk.Set, error) {
	certs, err := http.Query("GET", url, http.Options{})
	if err != nil {
		return nil, err
	}

	keyset, err := jwk.Parse([]byte(certs))
	if err != nil {
		return nil, err
	}

	pubkeyset, _ := jwk.PublicSetOf(keyset)
	return &pubkeyset, nil
}
