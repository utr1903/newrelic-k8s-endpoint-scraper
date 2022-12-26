package main

import (
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
	forwarder "github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/forward"
	scraper "github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/scrape"
)

func main() {

	// Parse and create cfg
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// Scrape endpoints
	scraper := scraper.NewScraper(cfg)
	evs := scraper.Run()

	// Forward endpoint values
	forwarder := forwarder.NewForwarder(cfg, evs)
	forwarder.Run()
}
