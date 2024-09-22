package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
)

// NewJsonRequest creates a new http request with the given object as json body.
func NewJsonRequest(route string, object any, httpType string) *http.Request {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("could not create json of object: %s", object)
	}

	req := httptest.NewRequest(httpType, route, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	return req
}

// NewPostJsonRequest creates a new http POST request with the given object as json body.
func NewPostJsonRequest(route string, object any) *http.Request {
	return NewJsonRequest(route, object, "POST")
}

// NewPutJsonRequest creates a new http PUT request with the given object as json body.
func NewPutJsonRequest(route string, object any) *http.Request {
	return NewJsonRequest(route, object, "PUT")
}
