package vatusa

import (
	"github.com/dhawton/log4g"
	"github.com/vzau/common/utils"
	"github.com/vzau/thoth/pkg/http"
)

var BASE_URL = "https://api.vatusa.net/v2"

var log = log4g.Category("vatusa")

func DeleteVisitor(cid uint64) error {
	_, err := http.Query("DELETE", facilityPath("roster/manageVisitor/%d", cid), http.Options{
		Headers: map[string]string{
			"Authorization": "APIKey " + utils.Getenv("VATUSA_API_KEY", ""),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func DeleteController(cid uint64) error {
	_, err := http.Query("DELETE", facilityPath("roster/%d", cid), http.Options{
		Headers: map[string]string{
			"Authorization": "APIKey " + utils.Getenv("VATUSA_API_KEY", ""),
		},
	})

	if err != nil {
		return err
	}

	return nil
}
