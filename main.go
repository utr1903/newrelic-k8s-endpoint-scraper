package main

import (
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
	forwarder "github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/forward"
	scraper "github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/scrape"
)

func main() {

	// Parse and create config
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// Scrape endpoints
	scraper := scraper.NewScraper(config)
	scraper.Run()

	// Forward endpoint values
	forwarder := forwarder.NewForwarder(config, scraper.GetEndpointValues())
	forwarder.Run()
}
