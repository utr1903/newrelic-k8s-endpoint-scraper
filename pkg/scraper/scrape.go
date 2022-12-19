package scraper

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
)

type EndpointScraper struct {
	config *config.Config
	client *http.Client

	mux            *sync.RWMutex
	EndpointValues map[config.Endpoint](map[string]string)
}

// Creates new scraper for endpoints
func New(
	cfg *config.Config,
) *EndpointScraper {

	// Create HTTP client
	client := http.Client{Timeout: time.Duration(30 * time.Second)}

	return &EndpointScraper{
		config:         cfg,
		client:         &client,
		mux:            &sync.RWMutex{},
		EndpointValues: make(map[config.Endpoint](map[string]string)),
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

	s.mux.Lock()
	s.addEndpointValues(endpoint, p.Run(data))
	s.mux.Unlock()
}

func (s *EndpointScraper) addEndpointValues(
	endpoint config.Endpoint,
	values map[string]string,
) {
	s.mux.Lock()
	s.EndpointValues[endpoint] = values
	s.mux.Unlock()
}

func (s *EndpointScraper) getEndpointValues(
	endpoint config.Endpoint,
) map[string]string {
	s.mux.RLock()
	values := s.EndpointValues[endpoint]
	s.mux.RUnlock()
	return values
}
