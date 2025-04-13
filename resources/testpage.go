package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	conf "github.com/LysetsDal/webscraper-v2/config"
	util "github.com/LysetsDal/webscraper-v2/utils"
	"github.com/gorilla/mux"
)

func serveStaticPage(w http.ResponseWriter, r *http.Request) error {
	data, err := os.ReadFile(conf.OFFLINE_TEST_PAGE_PATH) // returns []byte
	if err != nil {
		fmt.Printf("Error writing data: %v\n", err)
		return err
	}

	fmt.Printf("Served static page to: %v\n", r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err = w.Write(data)
	if err != nil {
		fmt.Printf("Error writing data: %v\n", err)
	}

	return err
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", util.MakeHttpHandleFunc(serveStaticPage))
	fmt.Printf("Started test server on port: %s\n", conf.OFFLINE_PORT)
	http.ListenAndServe(":"+conf.OFFLINE_PORT, router)

	// Graceful shutdown on ctrl + c:
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// This Blocks until signal received from c
	QUIT := <-c
	fmt.Println("Got Signal: ", QUIT)
}
