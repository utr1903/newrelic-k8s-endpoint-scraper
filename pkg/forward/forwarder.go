package forward

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
)

const NEW_RELIC_CUSTOM_EVENT_NAME = "K8sCustomEndpointScrapeSample"

type Forwarder struct {
	config *config.Config
	client *http.Client
	evs    *config.EndpointValues
}

func New(
	cfg *config.Config,
	evs *config.EndpointValues,
) *Forwarder {

	// Create HTTP client
	client := http.Client{Timeout: time.Duration(30 * time.Second)}

	return &Forwarder{
		config: cfg,
		client: &client,
		evs:    evs,
	}
}

func (f *Forwarder) Run() {

	// Create New Relic events
	nrEvents := f.createNewRelicEvents()

	// Flush data to New Relic
	f.sendToNewRelic(nrEvents)
}

func (f *Forwarder) createNewRelicEvents() []map[string]string {

	endpoints := f.evs.GetEndpoints()

	// Initialize to be sent New Relic events
	nrEvents := make([]map[string]string, len(endpoints))

	for _, endpoint := range endpoints {

		// All of the events are to be stored under "K8sCustomEndpointScrapeSample"
		nrEvent := map[string]string{
			"custom.eventType":    NEW_RELIC_CUSTOM_EVENT_NAME,
			"custom.endpointType": endpoint.Type,
			"custom.endpointName": endpoint.Name,
			"custom.endpointUrl":  endpoint.URL,
		}

		for endpointKey, endpointValue := range f.evs.GetEndpointValues(endpoint) {
			nrEvent[endpointKey] = endpointValue
		}
		nrEvents = append(nrEvents, nrEvent)
	}

	return nrEvents
}

func (f *Forwarder) sendToNewRelic(
	nrEvents []map[string]string,
) {

	// Create payload
	json, err := json.Marshal(nrEvents)
	if err != nil {
		return
	}
	payload := bytes.NewReader(json)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, f.config.Newrelic.EventsEndpoint, payload)
	if err != nil {
		return
	}

	// Perform HTTP request
	res, err := f.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
}
