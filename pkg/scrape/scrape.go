package scraper

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/logging"
)

// Object which is responsible for scraping
type EndpointScraper struct {
	config *config.Config
	client *http.Client
	evs    *config.EndpointValues
}

var readResponseBody = func(
	body io.ReadCloser,
) (
	[]byte,
	error,
) {
	return ioutil.ReadAll(body)
}

// Creates new scraper for endpoints
func NewScraper(
	cfg *config.Config,
) *EndpointScraper {

	// Create HTTP client
	client := http.Client{Timeout: time.Duration(30 * time.Second)}

	evs := config.NewEndpointValues()

	cfg.Logger.Log(logrus.DebugLevel, "Scraper is succesfully initialized.")

	return &EndpointScraper{
		config: cfg,
		client: &client,
		evs:    evs,
	}
}

// Scrape endpoints
func (s *EndpointScraper) Run() *config.EndpointValues {

	// Loop & scrape all endpoints
	s.config.Logger.Log(logrus.DebugLevel, "Looping over the endpoints to scrape...")
	for _, endpoint := range s.config.Endpoints {

		s.config.Logger.LogWithFields(logrus.DebugLevel, "Scraping endpoint...",
			map[string]string{
				"endpointType": endpoint.Type,
				"endpointName": endpoint.Name,
				"endpointUrl":  endpoint.URL,
			})

		// Create HTTP request
		req, err := http.NewRequest(http.MethodGet, endpoint.URL, nil)
		if err != nil {
			s.config.Logger.LogWithFields(logrus.ErrorLevel, logging.SCRAPE__HTTP_REQUEST_COULD_NOT_BE_CREATED,
				map[string]string{
					"endpointType": endpoint.Type,
					"endpointName": endpoint.Name,
					"endpointUrl":  endpoint.URL,
					"error":        err.Error(),
				})
			continue
		}

		// Perform HTTP request
		res, err := s.client.Do(req)
		if err != nil {
			s.config.Logger.LogWithFields(logrus.ErrorLevel, logging.SCRAPE__HTTP_REQUEST_HAS_FAILED,
				map[string]string{
					"endpointType": endpoint.Type,
					"endpointName": endpoint.Name,
					"endpointUrl":  endpoint.URL,
					"error":        err.Error(),
				})
			continue
		}
		defer res.Body.Close()

		// Check if call was successful
		if res.StatusCode != http.StatusOK {
			s.config.Logger.LogWithFields(logrus.ErrorLevel, logging.SCRAPE__ENDPOINT_RETURNED_NOT_OK_STATUS,
				map[string]string{
					"endpointType": endpoint.Type,
					"endpointName": endpoint.Name,
					"endpointUrl":  endpoint.URL,
				})
			continue
		}

		// Extract response body
		body, err := readResponseBody(res.Body)
		if err != nil {
			s.config.Logger.LogWithFields(logrus.ErrorLevel, logging.SCRAPE__RESPONSE_BODY_COULD_NOT_BE_PARSED,
				map[string]string{
					"endpointType": endpoint.Type,
					"endpointName": endpoint.Name,
					"endpointUrl":  endpoint.URL,
					"error":        err.Error(),
				})
			continue
		}

		// Parse response body
		switch endpoint.Type {
		case "kvp":
			s.parse(&KvpParser{}, endpoint, body)
		}
	}

	return s.evs
}

func (s *EndpointScraper) parse(
	p Parser,
	endpoint config.Endpoint,
	data []byte,
) {
	s.evs.AddEndpointValues(endpoint, p.Run(data))

	s.config.Logger.LogWithFields(logrus.DebugLevel, "Endpoint values are parsed.",
		map[string]string{
			"endpointType": endpoint.Type,
			"endpointName": endpoint.Name,
			"endpointUrl":  endpoint.URL,
		})
}
