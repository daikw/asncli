package asana

import "errors"

var ErrUnauthorized = errors.New("unauthorized")

type errorResponse struct {
	Errors []struct {
		Message   string `json:"message"`
		Help      string `json:"help"`
		ErrorCode string `json:"error_code"`
	} `json:"errors"`
}
