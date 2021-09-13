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

package jwt

import (
	"strconv"

	"github.com/lestrrat-go/jwx/jwt"
	"github.com/vzau/thoth/pkg/jwk"
)

func GetCIDFromToken(token string, jwkUrl string) (uint64, error) {
	keyset, err := jwk.FetchKeySet(jwkUrl)
	if err != nil {
		return 0, err
	}

	parsedToken, err := jwt.Parse([]byte(token), jwt.WithKeySet(*keyset), jwt.WithValidate(true))
	if err != nil {
		return 0, err
	}

	cid, err := strconv.ParseUint(parsedToken.Subject(), 10, 32)
	if err != nil {
		return 0, err
	}

	return cid, nil
}
