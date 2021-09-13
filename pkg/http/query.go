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
