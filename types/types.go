package types

import (
	"golang.org/x/net/html"
	"net/http"
)

// ApiFunc Decorator pattern to 'wrap' http.HandlerFunc
type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type ApiMessage struct {
	Message string `json:"message"`
}

type HtmlResponse struct {
	Body html.Node `json:"Response"`
}

type RequestData struct {
	URL string `json:"url"`
}
