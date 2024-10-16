package scraper

import (
	"fmt"
	. "github.com/LysetsDal/webscraper-v2/types"
	. "github.com/LysetsDal/webscraper-v2/utils"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Handler struct {
	Scraper http.Client
}

func NewHandler(client http.Client) *Handler {
	return &Handler{
		Scraper: client,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/test", MakeHttpHandleFunc(h.handleGetScrapeBody))
}

func (h *Handler) handleGetScrapeBody(w http.ResponseWriter, r *http.Request) error {
	url := "https://findbolig.nu/da-dk/udlejere"

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	body := string(bytes)

	fmt.Printf("Scraped URL: %s\n", url)

	WriteJson(w, http.StatusOK, ApiMessage{Message: body})

	return nil
}
