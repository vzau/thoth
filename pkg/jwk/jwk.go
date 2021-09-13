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
