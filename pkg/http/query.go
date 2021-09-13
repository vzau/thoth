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

package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dhawton/log4g"
)

var log = log4g.Category("http/query")

type Options struct {
	Headers  map[string]string
	PostData interface{}
}

func Query(method string, url string, options Options) (string, error) {
	body, _ := json.Marshal(options.PostData)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	for k, v := range options.Headers {
		req.Header.Set(k, v)
	}

	if options.PostData != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for key, val := range via[0].Header {
			req.Header[key] = val
		}

		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	return string(data), nil
}
