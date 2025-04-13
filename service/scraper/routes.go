package scraper

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LysetsDal/webscraper-v2/cmd/storage"
	T "github.com/LysetsDal/webscraper-v2/types"
	util "github.com/LysetsDal/webscraper-v2/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	Scraper   http.Client
	TargetURL string
	Store     *storage.PostgresStore
}

func NewHandler(client http.Client, store *storage.PostgresStore, targetURL string) *Handler {
	return &Handler{
		Scraper:   client,
		TargetURL: targetURL,
		Store:     store,
	}
}

func (h *Handler) UpdateTargetUrl(newURL string) string {
	h.TargetURL = newURL
	return fmt.Sprintf("Updated target URL to: %s", newURL)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/find/element/tag/{element}", util.MakeHttpHandleFunc(h.handleGetElementTag))
	router.HandleFunc("/find/elements/tag/{element1}&{element2}", util.MakeHttpHandleFunc(h.handleGetElementsTags))
	router.HandleFunc("/find/table/entries", util.MakeHttpHandleFunc(h.handleGetTableEntries))
	// TEST :
	router.HandleFunc("/insert/table/entries", util.MakeHttpHandleFunc(h.handleUpdateTableEntries))
}

// ====================================================================
// ========================= DATABASE =================================
// ====================================================================

func (h *Handler) handleGetTableEntries(w http.ResponseWriter, r *http.Request) error {
	resp, _ := http.Get(util.HttpPrefix(h.TargetURL))

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	// Parse response body
	doc := util.GetHTMLDoc(resp)
	table, _ := util.GetHTMLElement(doc, "table")

	entries, err := util.ParseOpenList(table)
	if err != nil {
		fmt.Println("List parse error: ", err)
		return err
	}

	// server-side logging
	fmt.Printf("[handleGetTableEntries] Scraped URL: %s\n", h.TargetURL)

	// Response to client
	return util.WriteJson(w, http.StatusOK, T.ListMessage{Length: len(entries), List: entries}) 
}

func (h *Handler) handleUpdateTableEntries(w http.ResponseWriter, r *http.Request) error {
	resp, _ := http.Get(util.HttpPrefix(h.TargetURL))

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	// Parse response body
	doc := util.GetHTMLDoc(resp)

	// Extract entries
	table, _ := util.GetHTMLElement(doc, "table")
	entries, err := util.ParseOpenList(table)
	if err != nil {
		fmt.Println("List parse error: ", err)
		return err
	}

	// server-side logging
	fmt.Printf("[handleUpdateTableEntries] Scraped URL: %s\n", h.TargetURL)

	// Persist data
	if err = h.Store.InsertEntries(entries); err != nil {
		return err
	}

	// Response to client
	return util.WriteJson(w, http.StatusOK, T.ListMessage{Length: len(entries), List: entries})
}

// ====================================================================
// ====================== HTML ELEMENTS ===============================
// ====================================================================

func (h *Handler) handleGetElementTag(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	resp, _ := http.Get(util.HttpPrefix(h.TargetURL))

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	// Parse response body
	doc := util.GetHTMLDoc(resp)
	body, _ := util.GetHTMLBody(doc)

	tags, err := util.GetHTMLElement(body, pathVars["element"])
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err)
	}

	// server-side logging
	fmt.Printf("[handleGetElementTag] Scraped URL: %s\n", h.TargetURL)

	// Response to client
	return util.WriteJson(w, http.StatusOK, T.ApiMessage{Message: util.PrettyPrint(tags)})
}

func (h *Handler) handleGetElementsTags(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	resp, _ := http.Get(util.HttpPrefix(h.TargetURL))

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	// Parse response body
	doc := util.GetHTMLDoc(resp)
	body, _ := util.GetHTMLBody(doc)

	tags, err := util.GetNestedHTMLElement(body, pathVars["element1"], pathVars["element2"])
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err)
	}

	// server-side logging
	fmt.Printf("[handleGetElementsTags] Scraped URL: %s\n", h.TargetURL)

	// Response to client
	return util.WriteJson(w, http.StatusOK, T.ApiMessage{Message: util.PrettyPrint(tags)})
}
