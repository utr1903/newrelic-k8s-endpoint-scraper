package forward

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
)

const NEW_RELIC_CUSTOM_EVENT_NAME = "K8sCustomEndpointScrapeSample"

type Forwarder struct {
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

func NewForwarder(
	cfg *config.Config,
	evs *config.EndpointValues,
) *Forwarder {

	// Create HTTP client
	client := http.Client{Timeout: time.Duration(30 * time.Second)}

	cfg.Logger.Log(logrus.DebugLevel, "Endpoint values are parsed.")

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

	f.config.Logger.Log(logrus.DebugLevel, "Creating New Relic events...")

	endpoints := f.evs.GetEndpoints()

	// Initialize to be sent New Relic events
	nrEvents := make([]map[string]string, 0, len(endpoints))

	for _, endpoint := range endpoints {

		// All of the events are to be stored under "K8sCustomEndpointScrapeSample"
		nrEvent := map[string]string{
			"eventType":           NEW_RELIC_CUSTOM_EVENT_NAME,
			"custom.endpointType": endpoint.Type,
			"custom.endpointName": endpoint.Name,
			"custom.endpointUrl":  endpoint.URL,
		}

		for endpointKey, endpointValue := range f.evs.GetEndpointValues(endpoint) {
			nrEvent[endpointKey] = endpointValue
		}
		nrEvents = append(nrEvents, nrEvent)
	}

	f.config.Logger.Log(logrus.DebugLevel, "New Relic events are created successfully.")
	return nrEvents
}

func (f *Forwarder) sendToNewRelic(
	nrEvents []map[string]string,
) {

	// Create payload
	f.config.Logger.Log(logrus.DebugLevel, "Creating payload...")
	json, err := json.Marshal(nrEvents)
	if err != nil {
		f.config.Logger.LogWithFields(logrus.ErrorLevel, "Payload could not be created.",
			map[string]string{
				"error": err.Error(),
			})
		return
	}
	payload := bytes.NewReader(json)

	// Create HTTP request
	f.config.Logger.Log(logrus.DebugLevel, "Creating HTTP request...")
	req, err := http.NewRequest(http.MethodPost, f.config.Newrelic.EventsEndpoint, payload)
	if err != nil {
		f.config.Logger.LogWithFields(logrus.ErrorLevel, "HTTP request could not be created.",
			map[string]string{
				"error": err.Error(),
			})
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Api-Key", f.config.Newrelic.LicenseKey)

	// Perform HTTP request
	f.config.Logger.Log(logrus.DebugLevel, "Performing HTTP request...")
	res, err := f.client.Do(req)
	if err != nil {
		msg := "HTTP request has failed."
		f.config.Logger.LogWithFields(logrus.ErrorLevel, msg,
			map[string]string{
				"error": err.Error(),
			})
		panic(msg)
	}
	defer res.Body.Close()

	// Check if call was successful
	if res.StatusCode == http.StatusOK {
		f.config.Logger.Log(logrus.DebugLevel, "New Relic events are forwarded successfully.")
	} else {
		msg := "HTTP request has returned not OK status."
		f.config.Logger.Log(logrus.ErrorLevel, msg)
		panic(msg)
	}

	// // Extract response body
	// body, err := readResponseBody(res.Body)
	// if err != nil {
	// 	f.config.Logger.LogWithFields(logrus.ErrorLevel, "Response body could not be parsed.",
	// 		map[string]string{
	// 			"error": err.Error(),
	// 		})
	// 	return
	// }
}
