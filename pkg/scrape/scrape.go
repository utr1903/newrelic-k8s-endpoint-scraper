package scraper

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
)

// Object which is responsible for scraping
type EndpointScraper struct {
	config *config.Config
	client *http.Client
	evs    *config.EndpointValues
}

// Creates new scraper for endpoints
func NewScraper(
	cfg *config.Config,
) *EndpointScraper {

	// Create HTTP client
	client := http.Client{Timeout: time.Duration(30 * time.Second)}

	evs := config.NewEndpointValues()

	return &EndpointScraper{
		config: cfg,
		client: &client,
		evs:    evs,
	}
}

// Scrape endpoints
func (s *EndpointScraper) Run() {

	// Loop & scrape all endpoints
	for _, endpoint := range s.config.Endpoints {

		// Create HTTP request
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

		// Parse response body
		switch endpoint.Type {
		case "kvp":
			s.parse(&KvpParser{}, endpoint, body)
		}
	}
}

func (s *EndpointScraper) parse(
	p Parser,
	endpoint config.Endpoint,
	data []byte,
) {
	s.evs.AddEndpointValues(endpoint, p.Run(data))
}

func (s *EndpointScraper) GetEndpointValues() *config.EndpointValues {
	return s.evs
}
