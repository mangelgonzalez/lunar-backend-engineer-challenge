package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func ClientIp(req *http.Request) string {
	ipAddress := req.RemoteAddr
	fwdAddress := req.Header.Get("X-Forwarded-For")
	if fwdAddress != "" {
		ipAddress = fwdAddress

		ips := strings.Split(fwdAddress, ", ")
		if len(ips) > 1 {
			ipAddress = ips[0]
		}
	}

	return ipAddress
}

func CloneRequest(r *http.Request) *http.Request {
	var bodyBytes []byte
	newRequest := *r.WithContext(r.Context())

	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	newRequest.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return &newRequest
}

func AllParamsRequest(r *http.Request) (map[string]interface{}, error) {
	allParams := make(map[string]interface{}, 0)

	newRequest := CloneRequest(r)

	allParams, err := ConvertRequestToBodyMap(newRequest)
	if err != nil {
		return allParams, err
	}

	queryParams := QueryParams(newRequest)

	for k, v := range queryParams {
		allParams[k] = v
	}

	return allParams, nil
}

func ConvertRequestToBodyMap(r *http.Request) (map[string]interface{}, error) {
	requestBody := make(map[string]interface{}, 0)
	var err error
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return requestBody, err
	}
	if err = json.Unmarshal(b, &requestBody); err != nil {
		return requestBody, err
	}

	return requestBody, err
}

func QueryParams(r *http.Request) map[string]interface{} {
	query := make(map[string]interface{})

	for k, v := range r.URL.Query() {
		query[k] = v[0]
	}

	return query
}
