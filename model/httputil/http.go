package httputil

type HTTPError struct {
	Error   string `json:"error"`
	Success bool   `json:"success" example:"false"`
}
