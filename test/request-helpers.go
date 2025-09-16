package test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
)

var routerHandler http.Handler

func executeAuthenticatedJsonRequest(t *testing.T, verb string, path string, body []byte, token string, headers map[string]string) *httptest.ResponseRecorder {
	return executeAuthenticatedRequest(t, verb, path, body, token, "application/json", headers)
}

func executeUnauthenticatedJsonRequest(t *testing.T, verb string, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	return executeUnauthenticatedRequest(t, verb, path, body, "application/json", headers)
}

func executeAuthenticatedRequest(t *testing.T, verb string, path string, body []byte, token string, requestType string, headers map[string]string) *httptest.ResponseRecorder {
	req := buildRequest(t, verb, path, body, headers)

	req.Header.Set("Content-Type", requestType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return executeRequest(req)
}

func executeUnauthenticatedRequest(t *testing.T, verb string, path string, body []byte, requestType string, headers map[string]string) *httptest.ResponseRecorder {
	req := buildRequest(t, verb, path, body, headers)
	req.Header.Set("Content-Type", requestType)

	return executeRequest(req)
}

func buildRequest(t *testing.T, method string, url string, body []byte, headers map[string]string) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	assert.NoError(t, err)

	if len(headers) != 0 {
		for headerName, value := range headers {
			req.Header.Set(headerName, value)
		}
	}

	return req
}

func executeJsonApiRequest(t *testing.T, method, url string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	assert.NoError(t, err)

	if len(headers) != 0 {
		for headerName, value := range headers {
			req.Header.Set(headerName, value)
		}
	}
	req.Header.Set("Content-Type", jsonapi.MediaType)

	return executeRequest(req)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	if routerHandler == nil {
		routerHandler = common.Services.Router.BuildHandler()
	}

	routerHandler.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func checkResponse(t *testing.T, expectedStatusCode int, expectedResponse string, response *httptest.ResponseRecorder, formats ...interface{}) {
	ja := jsonassert.New(t)
	checkResponseCode(t, expectedStatusCode, response.Code)

	receivedResponse := response.Body.String()

	if receivedResponse == "" {
		assert.Equal(t, expectedResponse, receivedResponse)
		return
	}
	fmt.Println(receivedResponse)
	if formats != nil {
		ja.Assertf(receivedResponse, expectedResponse, formats)
	} else {
		ja.Assert(receivedResponse, expectedResponse)
	}
}

func EmptyHeaders() map[string]string {
	return map[string]string{}
}

func AssertAskResponse(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected response %v. Got %v\n", expected, actual)
	}
}
