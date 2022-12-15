package main

import (
	"fmt"

	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	fmt.Print(cfg.Endpoints[0].Name)
	fmt.Print(cfg.Endpoints[0].Type)
	fmt.Print(cfg.Endpoints[0].URL)
}
