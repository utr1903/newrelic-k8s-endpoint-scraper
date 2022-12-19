package main

import (
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/scraper"
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
