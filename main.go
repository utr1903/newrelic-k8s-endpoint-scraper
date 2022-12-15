package main

import (
	"fmt"

	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/scraper"
)

func main() {

	config, err := config.New()
	if err != nil {
		panic(err)
	}

	fmt.Print(config.Endpoints[0].Name)
	fmt.Print(config.Endpoints[0].Type)
	fmt.Print(config.Endpoints[0].URL)

	scraper := scraper.New(config)
	scraper.Run()
}
