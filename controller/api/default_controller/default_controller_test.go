package default_controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
)

func newJsonRequest(route string, object any, httpType string) *http.Request {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("could not create json of object: %s", object)
	}

	req := httptest.NewRequest(httpType, route, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func newPostJsonRequest(route string, object any) *http.Request {
	return newJsonRequest(route, object, "POST")
}

func newPutJsonRequest(route string, object any) *http.Request {
	return newJsonRequest(route, object, "PUT")
}
