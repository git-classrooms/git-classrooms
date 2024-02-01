package default_controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
)

func newPostJsonRequest(route string, object any) *http.Request {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("could not create json of object: %s", object)
	}

	req := httptest.NewRequest("POST", route, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	return req
}
