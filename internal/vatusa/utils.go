package vatusa

import (
	"fmt"
	"net/url"

	"github.com/vzau/common/utils"
)

func generateUrl(path string, args ...interface{}) string {
	url, err := url.Parse(BASE_URL + fmt.Sprintf(path, args...))
	if err != nil {
		log.Error("Error parsing URL, may not be valid: %s %s", BASE_URL+fmt.Sprintf(path, args...), err.Error())
		return ""
	}

	if utils.Getenv("APP_ENV", "dev") == "prod" {
		return url.String()
	} else {
		q := url.Query()
		q.Set("test", "")
		url.RawQuery = q.Encode()
		return url.String()
	}
}

func facilityPath(subpath string, args ...interface{}) string {
	return generateUrl("/facility/"+utils.Getenv("FACILITY_ID", "ZAU")+"/"+subpath, args...)
}
