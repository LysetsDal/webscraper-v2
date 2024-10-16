package types

import "net/http"

// ApiFunc Decorator pattern to 'wrap' http.HandlerFunc
type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type ApiMessage struct {
	Message string `json:"message"`
}
