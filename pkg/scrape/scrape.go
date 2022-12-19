package scraper

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
)

// Object to store all values of all endpoints
type EndpointValues struct {
	// To avoid multi-thread read/write into the map
	mux *sync.RWMutex

	// Map to store all values according to endpoints
	// -> Key: endpoint itself
	// -> Val: attributes which the endpoint has exposed
	Values map[config.Endpoint](map[string]string)
}

// Object which is responsible for scraping
type EndpointScraper struct {
	config         *config.Config
	client         *http.Client
	EndpointValues EndpointValues
}

// Creates new scraper for endpoints
func New(
	cfg *config.Config,
) *EndpointScraper {

	// Create HTTP client
	client := http.Client{Timeout: time.Duration(30 * time.Second)}

	evs := EndpointValues{
		mux:    &sync.RWMutex{},
		Values: make(map[config.Endpoint](map[string]string)),
	}

	return &EndpointScraper{
		config:         cfg,
		client:         &client,
		EndpointValues: evs,
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
	s.addEndpointValues(endpoint, p.Run(data))
}

func (s *EndpointScraper) addEndpointValues(
	endpoint config.Endpoint,
	values map[string]string,
) {
	s.EndpointValues.mux.Lock()
	s.EndpointValues.Values[endpoint] = values
	s.EndpointValues.mux.Unlock()
}

func (s *EndpointScraper) getEndpointValues(
	endpoint config.Endpoint,
) map[string]string {
	s.EndpointValues.mux.RLock()
	values := s.EndpointValues.Values[endpoint]
	s.EndpointValues.mux.RUnlock()
	return values
}
