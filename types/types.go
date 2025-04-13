package types

import (
	// "math/rand"
	"net/http"
	"time"
)

type CreateEntryRequest struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// ApiFunc Decorator pattern to 'wrap' http.HandlerFunc
type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type ApiMessage struct {
	Message string `json:"message"`
}

type Entry struct {
	ID        int       `json:"ID"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewEntry(name, status string) *Entry {
	return &Entry{
		// ID: rand.Intn(10000),
		Name:   name,
		Status: status,
		UpdatedAt: time.Now().UTC(),
	}
}

type ListMessage struct {
	Length int     `json:"length"`
	List   []Entry `json:"entry"`
}
