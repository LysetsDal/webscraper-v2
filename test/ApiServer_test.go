package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/html"

	"github.com/LysetsDal/webscraper-v2/service/scraper"
	. "github.com/LysetsDal/webscraper-v2/types"
	"github.com/LysetsDal/webscraper-v2/utils"
	"github.com/gorilla/mux"
)

const html_data = `
	<html>
		<body>
			<table>
				<tr><td>Fonden Vestergården</td><td>Åben</td></tr>
				<tr><td>Frederiksberg Boligfond</td><td>Lukket</td></tr>
			</table>
		</body>
	</html>`

func setupTestServer(targetHtml []byte) *httptest.Server {
	// Simulate target HTML page
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(targetHtml)
	}))

	// Build router just like in WebScraper.Run()
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v2").Subrouter()

	client := http.Client{}
	scraperHandler := scraper.NewHandler(client, nil, target.URL)
	scraperHandler.RegisterRoutes(subrouter)

	// Include HomeHandler if needed
	subrouter.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		utils.WriteJson(w, http.StatusOK, "home")
	})

	return httptest.NewServer(router)
}

func TestFindElementTagEndpoint(t *testing.T) {
	// Arrange
	ts := setupTestServer([]byte(html_data))
	defer ts.Close()

	// Act
	result, err := http.Get(ts.URL + "/api/v2/find/element/tag/body")
	if err != nil {
		t.Fatalf("failed to GET endpoint: %v", err)
	}
	defer result.Body.Close()

	doc, err := html.Parse(result.Body)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	// Assert
	if result.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", result.StatusCode)
	}

	body := utils.ExtractText(doc)

	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Fonden") || !strings.Contains(bodyStr, "Lukket") {
		t.Errorf("unexpected response body: %s", bodyStr)
	}
}

func TestGetTableEntries(t *testing.T) {
	// Arrange
	var expected_length int = 2

	ts := setupTestServer([]byte(html_data))
	defer ts.Close()

	// Act
	result, err := http.Get(ts.URL + "/api/v2/find/table/entries")
	if err != nil {
		t.Fatalf("failed to GET endpoint: %v", err)
	}
	defer result.Body.Close()

	var msg ListMessage
	err = utils.ReadJson(result.Body, &msg)
	if err != nil {
		t.Fatalf("Error reading JSON: %v", err)
	}

	// Assert
	if result.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", result.StatusCode)
	}

	if msg.Length != expected_length {
		t.Errorf("Result was incorrect. Expected: %v Got: %v", expected_length, msg.Length)
	}

	var e0 Entry = msg.List[0]
	var e1 Entry = msg.List[1]

	if e0.Name != "Fonden Vestergården" || e0.Status != "Åben" {
		t.Errorf("Result was incorrect. Expected: (%v, %v)", e0.Name, e0.Status)
	}

	if e1.Name != "Frederiksberg Boligfond" || e1.Status != "Lukket" {
		t.Errorf("Result was incorrect. Expected: (%v, %v)", e1.Name, e1.Status)
	}

}
