package model

import (
	"fmt"
	"net/http"
	"net/url"
)

type GitLabError struct {
	Body     []byte
	Response *http.Response
	Message  string
}

func (e GitLabError) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}
