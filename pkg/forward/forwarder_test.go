package forward

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/logging"
)

func Test_NewRelicEventsAreCreated(t *testing.T) {
	endpointInfoMock := createEndpointInfoMock()
	cfg := createConfig("", endpointInfoMock)
	evs := createEndpointValues(cfg, endpointInfoMock)

	forwarder := NewForwarder(cfg, evs)
	nrEvents := forwarder.createNewRelicEvents()

	for counter, nrEvent := range nrEvents {
		switch counter {
		case 0:
			assert.Equal(t, "ep1Url", nrEvent["endpointUrl"])
			assert.Equal(t, "MyEndpoint"+"ep1Url", nrEvent["eventType"])
			assert.Equal(t, "kvp", nrEvent["endpointType"])
			assert.Equal(t, "v1", nrEvent["k1"])
			assert.Equal(t, "v2", nrEvent["k2"])
		case 1:
			assert.Equal(t, "ep2Url", nrEvent["endpointUrl"])
			assert.Equal(t, "MyEndpoint"+"ep2Url", nrEvent["eventType"])
			assert.Equal(t, "kvp", nrEvent["endpointType"])
			assert.Equal(t, "v3", nrEvent["k3"])
			assert.Equal(t, "v4", nrEvent["k4"])
		}
	}
}

func Test_HttpRequestCouldNotBeCreated(t *testing.T) {
	endpointInfoMock := createEndpointInfoMock()
	cfg := createConfig("::", endpointInfoMock)
	evs := createEndpointValues(cfg, endpointInfoMock)

	forwarder := NewForwarder(cfg, evs)
	err := forwarder.Run()

	assert.NotNil(t, err)
	assert.Equal(t, logging.FORWARD__HTTP_REQUEST_COULD_NOT_BE_CREATED, err.Error())
}

func Test_HttpRequestHasFailed(t *testing.T) {
	endpointInfoMock := createEndpointInfoMock()
	cfg := createConfig("", endpointInfoMock)
	evs := createEndpointValues(cfg, endpointInfoMock)

	forwarder := NewForwarder(cfg, evs)
	err := forwarder.Run()

	assert.NotNil(t, err)
	assert.Equal(t, logging.FORWARD__HTTP_REQUEST_HAS_FAILED, err.Error())
}

func Test_NewRelicReturnsNotOkResponse(t *testing.T) {
	newrelicEventServerMock := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
	defer newrelicEventServerMock.Close()

	endpointInfoMock := createEndpointInfoMock()
	cfg := createConfig(newrelicEventServerMock.URL, endpointInfoMock)
	evs := createEndpointValues(cfg, endpointInfoMock)

	forwarder := NewForwarder(cfg, evs)
	err := forwarder.Run()

	assert.NotNil(t, err)
	assert.Equal(t, logging.FORWARD__NEW_RELIC_RETURNED_NOT_OK_STATUS, err.Error())
}

func Test_EventsAreSent(t *testing.T) {
	newrelicEventServerMock := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
	defer newrelicEventServerMock.Close()

	endpointInfoMock := createEndpointInfoMock()
	cfg := createConfig(newrelicEventServerMock.URL, endpointInfoMock)
	evs := createEndpointValues(cfg, endpointInfoMock)

	forwarder := NewForwarder(cfg, evs)
	err := forwarder.Run()
	assert.Nil(t, err)
}

func createEndpointValues(
	cfg *config.Config,
	endpointInfoMock map[string](map[string]string),
) *config.EndpointValues {
	evs := config.NewEndpointValues()
	for _, endpoint := range cfg.Endpoints {
		evs.AddEndpointValues(endpoint, endpointInfoMock[endpoint.URL])
	}
	return evs
}

func createEndpointInfoMock() map[string](map[string]string) {
	return map[string](map[string]string){
		"ep1Url": map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
		"ep2Url": map[string]string{
			"k3": "v3",
			"k4": "v4",
		},
	}
}

func createConfig(
	newrelicEventsUrl string,
	endpointInfo map[string](map[string]string),
) *config.Config {
	logLevel := "ERROR"
	eps := []config.Endpoint{}
	for url := range endpointInfo {
		eps = append(eps, config.Endpoint{
			Type: "kvp",
			Name: "MyEndpoint" + url,
			URL:  url,
		})
	}
	return &config.Config{
		Newrelic: &config.NewRelicInput{
			LogLevel:       logLevel,
			EventsEndpoint: newrelicEventsUrl,
			LicenseKey:     "",
		},
		Logger:    logging.NewLogger(logLevel),
		Endpoints: eps,
	}
}
