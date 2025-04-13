package main

import (
	"fmt"
	"log"

	"github.com/LysetsDal/webscraper-v2/cmd/api"
	"github.com/LysetsDal/webscraper-v2/cmd/storage"
	conf "github.com/LysetsDal/webscraper-v2/config"
)

func main() {
	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewWebScraper(":3030", store)
	fmt.Println(server.DisplayStartMsg())

	// server.Run(conf.TARGET_URL)
	server.Run(conf.OFFLINE_URL)
}
