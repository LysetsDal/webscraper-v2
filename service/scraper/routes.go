package scraper

import (
	"fmt"
	. "github.com/LysetsDal/webscraper-v2/types"
	. "github.com/LysetsDal/webscraper-v2/utils"
	"github.com/gorilla/mux"
	"golang.org/x/net/html"
	"log"
	"net/http"
)

type Handler struct {
	Scraper   http.Client
	TargetURL string
}

func NewHandler(client http.Client, targetURL string) *Handler {
	return &Handler{
		Scraper:   client,
		TargetURL: targetURL,
	}
}

func (h *Handler) UpdateTargetUrl(newURL string) string {
	h.TargetURL = newURL
	return fmt.Sprintf("Updated target URL to: %s", newURL)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/find/element/tag/{element}", MakeHttpHandleFunc(h.handleGetElementTag))
	router.HandleFunc("/find/elements/tag/{element1}&{element2}", MakeHttpHandleFunc(h.handleGetElementsTags))
	router.HandleFunc("/find/element/tag/{element}", MakeHttpHandleFunc(h.handleGetElementTagFromUrl))
}

func (h *Handler) handleGetElementTag(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	resp, _ := http.Get(h.TargetURL)

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	// Parse response body
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	body, _ := GetHTMLBody(doc)

	tags, err := getHTMLElement(body, pathVars["element"])
	if err != nil {
		WriteError(w, http.StatusNotFound, err)
	}

	// server-side logging
	fmt.Printf("Scraped URL: %s\n", h.TargetURL)

	// Response to client
	err = WriteJson(w, http.StatusOK, ApiMessage{Message: PrettyPrint(tags)})
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return nil
}

func (h *Handler) handleGetElementsTags(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	resp, _ := http.Get(h.TargetURL)

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	// Parse response body
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	body, _ := GetHTMLBody(doc)

	tags, err := getNestedHTMLElement(body, pathVars["element1"], pathVars["element2"])
	if err != nil {
		WriteError(w, http.StatusNotFound, err)
	}

	// server-side logging
	fmt.Printf("Scraped URL: %s\n", h.TargetURL)

	// Response to client
	err = WriteJson(w, http.StatusOK, ApiMessage{Message: PrettyPrint(tags)})
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return nil
}

func (h *Handler) handleGetElementTagFromUrl(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("error parsing form data"))
	}

	URL := r.FormValue("url")
	if URL == "" {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("missing 'url' in form data"))
		return fmt.Errorf("missing 'url' in form data")
	}

	h.UpdateTargetUrl(URL)

	return h.handleGetElementTag(w, r)
}

func getHTMLElement(doc *html.Node, elementTag string) (*html.Node, error) {
	var element *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == elementTag {
			element = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(doc)
	if element != nil {
		return element, nil
	}

	return nil, fmt.Errorf("missing <%s> in the node tree", elementTag)
}

func getNestedHTMLElement(doc *html.Node, outerTag string, innerTag string) (*html.Node, error) {

	// Find the Outer element (outerTag)
	outerElement, err := getHTMLElement(doc, outerTag)
	if err != nil {
		return nil, fmt.Errorf("missing guard element <%s> in the node tree", outerTag)
	}

	// Search for the inner element (innerTag) within outer element (outerTag)
	innerElement, err := getHTMLElement(outerElement, innerTag)
	if err != nil {
		return nil, fmt.Errorf("missing nested element <%s> inside <%s>", innerTag, outerTag)
	}

	return innerElement, nil
}
