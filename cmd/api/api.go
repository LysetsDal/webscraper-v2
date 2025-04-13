package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/LysetsDal/webscraper-v2/cmd/storage"
	"github.com/LysetsDal/webscraper-v2/service/scraper"
	util "github.com/LysetsDal/webscraper-v2/utils"
	"github.com/gorilla/mux"
)

const localhost string = "127.0.0.1:"

type WebScraper struct {
	Name          string
	ListenAddr    string
	StartTime     time.Time
	ScraperClient http.Client
	store         storage.PostgresStore
}

type ScraperData struct {
	Name          string `json:"Name"`
	ListenAddr    string `json:"ListenAddr"`
	StartTime     string `json:"Time"`
	ScraperClient string `json:"ScraperClient"`
}

func connectServer(port string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("tcp", localhost+port)
			},
		},
	}
}

func NewWebScraper(listenAddr string, store *storage.PostgresStore) *WebScraper {
	return &WebScraper{
		Name:          "WebScraper-v2",
		ListenAddr:    listenAddr,
		StartTime:     time.Now(),
		ScraperClient: *connectServer(listenAddr),
		store:         *store,
	}
}

func (s *WebScraper) Run(targetURL string) {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v2").Subrouter()

	scraperHandler := scraper.NewHandler(s.ScraperClient, &s.store, targetURL)
	scraperHandler.RegisterRoutes(subrouter) //<- register

	subrouter.HandleFunc("/", util.MakeHttpHandleFunc(s.HomeHandler))

	// Run server in go func so it doesn't block
	go func() {}()
	if err := http.ListenAndServe(s.ListenAddr, router); err != nil {
		return
	}

	// Graceful shutdown on ctrl + c:
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// This Blocks until signal received from c
	QUIT := <-c
	fmt.Println("Got Signal: ", QUIT)
}

func (s *WebScraper) DisplayStartMsg() string {
	return fmt.Sprintf("Started web-scraper service on port: %s", s.ListenAddr)
}

func (s *WebScraper) HomeHandler(w http.ResponseWriter, _ *http.Request) error {

	data := ScraperData{
		Name:          s.Name,
		ListenAddr:    s.ListenAddr,
		StartTime:     s.StartTime.String(),
		ScraperClient: localhost,
	}

	return util.WriteJson(w, http.StatusOK, data)
}
