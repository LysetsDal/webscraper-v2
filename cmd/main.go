package main

import (
	"fmt"
	"github.com/LysetsDal/webscraper-v2/cmd/api"
)

func main() {
	server := api.NewWebScraper(":3000")
	fmt.Printf("Msg: %s", api.Test())
	server.Run()
}
