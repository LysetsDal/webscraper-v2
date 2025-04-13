package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	T "github.com/LysetsDal/webscraper-v2/types"
)

func ReadJson(body io.Reader, decodeTo any) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(decodeTo)
}

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

// WriteJson Write Json with standard header.
func WriteJson(w http.ResponseWriter, status int, encodeVal any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(encodeVal)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	_ = WriteJson(w, status, map[string]string{"error": err.Error()})
}

func MakeHttpHandleFunc(f T.ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			err := WriteJson(w, http.StatusBadRequest, T.ApiError{Error: err.Error()})
			if err != nil {
				return
			}
		}
	}
}

// prepends
func HttpPrefix(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}

