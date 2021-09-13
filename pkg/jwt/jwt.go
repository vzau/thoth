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
