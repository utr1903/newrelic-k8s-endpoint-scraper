package main

import (
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
	forwarder "github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/forward"
	scraper "github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/scrape"
)

func main() {

	// Parse and create config
	config := config.NewConfig()

	// Scrape endpoints
	scraper := scraper.NewScraper(config)
	scraper.Run()

	// Forward endpoint values
	forwarder := forwarder.New(config, scraper.GetEndpointValues())
	forwarder.Run()
}
