package common

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	HTTPStatusOk            = 200
	HTTPStatusBadRequest    = 400
	HTTPStatusNotFound      = 404
	HTTPStatusInternalError = 500
)

var transport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}
var httpClient = &http.Client{
	Timeout:   10 * time.Second,
	Transport: transport,
}

// MakeHTTPRequest is using the http.Client
// to do a HTTP request and return the response's body
func MakeHTTPRequest(
	method string,
	URL string,
	reqBody *[]byte,
	headers *map[string]string,
) (resBody []byte, err error) {
	var req *http.Request

	if reqBody != nil {
		buf := bytes.NewBuffer(*reqBody)
		req, err = http.NewRequest(method, URL, buf)
	} else {
		req, err = http.NewRequest(method, URL, nil)
	}
	if err != nil {
		return
	}

	if headers != nil {
		for key, val := range *headers {
			req.Header.Set(key, val)
		}
	}

	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := httpClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return
	}

	resBody, err = ioutil.ReadAll(res.Body)

	return
}

func WriteHTTPResponse(
	w http.ResponseWriter,
	statusCode int,
	body []byte,
	headers *map[string]string,
) {
	if headers != nil {
		for key, value := range *headers {
			w.Header().Set(key, value)
		}
	}

	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(statusCode)
	w.Write(body)
}
