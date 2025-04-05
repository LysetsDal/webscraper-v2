package main

import (
	"fmt"

	"github.com/LysetsDal/webscraper-v2/cmd/api"
	. "github.com/LysetsDal/webscraper-v2/config"
)

func main() {
	server := api.NewWebScraper(":3030")
	fmt.Println(server.DisplayStartMsg())
	server.Run(TARGET_URL)

}
