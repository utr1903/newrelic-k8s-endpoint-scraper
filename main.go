package main

import (
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
	scraper "github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/scrape"
)

func main() {

	// Parse and create config
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	// Scrape endpoints
	scraper := scraper.New(config)
	scraper.Run()
}
