package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
)

func NewJSONRequest(route string, object any, httpType string) *http.Request {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("could not create json of object: %s", object)
	}

	req := httptest.NewRequest(httpType, route, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func NewPostJSONRequest(route string, object any) *http.Request {
	return NewJSONRequest(route, object, "POST")
}

func NewPutJSONRequest(route string, object any) *http.Request {
	return NewJSONRequest(route, object, "PUT")
}
