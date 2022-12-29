package main

import (
	"fmt"

	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
	forwarder "github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/forward"
	scraper "github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/scrape"
)

func main() {

	// Parse and create config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// Scrape endpoints
	scraper := scraper.NewScraper(cfg)
	evs := scraper.Run()

	// Forward endpoint values to New Relic
	forwarder := forwarder.NewForwarder(cfg, evs)
	err = forwarder.Run()
	if err != nil {
		panic(err)
	}

	// Send the app logs to New Relic
	err = cfg.Logger.Flush()
	if err != nil {
		fmt.Println(err)
	}
}
