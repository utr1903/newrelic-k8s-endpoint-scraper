package scraper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
)

type EndpointScraper struct {
	config *config.Config
	client *http.Client
}

func New(
	config *config.Config,
) *EndpointScraper {

	// Create HTTP client
	client := http.Client{Timeout: time.Duration(30 * time.Second)}

	return &EndpointScraper{
		config: config,
		client: &client,
	}
}

func (s *EndpointScraper) Run() {

	// Scrape endpoints
	s.scrape()
}

func (s *EndpointScraper) scrape() {

	// Loop & scrape all endpoints
	for _, endpoint := range s.config.Endpoints {
		req, err := http.NewRequest(http.MethodGet, endpoint.URL, nil)
		if err != nil {
			return
		}

		// Perform HTTP request
		res, err := s.client.Do(req)
		if err != nil {
			return
		}
		defer res.Body.Close()

		// Extract response body
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}

		fmt.Println(string(body))
	}
}
