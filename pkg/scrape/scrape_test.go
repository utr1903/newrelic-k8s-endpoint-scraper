package scraper

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/logging"
)

func Test_EndpointReturnsNotOkResponse(t *testing.T) {
	endpointServerMock := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
	defer endpointServerMock.Close()

	cfg := createConfig([]string{
		endpointServerMock.URL,
	})
	scraper := NewScraper(cfg)
	evs := scraper.Run()

	assert.Equal(t, 0, len(evs.Values))
}

func Test_EndpointsAreScrapedSuccessfully(t *testing.T) {
	endpointServerMock1 := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var b bytes.Buffer
			b.WriteString("k1:v1" + "\n")
			b.WriteString("k2:v2" + "\n")

			w.WriteHeader(http.StatusOK)
			w.Write(b.Bytes())
		}))
	defer endpointServerMock1.Close()

	endpointServerMock2 := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var b bytes.Buffer
			b.WriteString("k3:v3" + "\n")
			b.WriteString("k4:v4" + "\n")

			w.WriteHeader(http.StatusOK)
			w.Write(b.Bytes())
		}))
	defer endpointServerMock2.Close()

	cfg := createConfig([]string{
		endpointServerMock1.URL,
		endpointServerMock2.URL,
	})
	scraper := NewScraper(cfg)
	evs := scraper.Run()

	assert.Equal(t, 2, len(evs.Values))

	for endpoint, values := range scraper.evs.Values {
		if endpoint.URL == endpointServerMock1.URL {
			assert.Equal(t, "v1", values["k1"])
			assert.Equal(t, "v2", values["k2"])
		}
		if endpoint.URL == endpointServerMock2.URL {
			assert.Equal(t, "v3", values["k3"])
			assert.Equal(t, "v4", values["k4"])
		}
	}
}

func createConfig(
	endpointUrls []string,
) *config.Config {
	logLevel := "ERROR"
	eps := []config.Endpoint{}
	for i, ep := range endpointUrls {
		eps = append(eps, config.Endpoint{
			Type: "kvp",
			Name: "MyEndpoint" + strconv.FormatInt(int64(i), 10),
			URL:  ep,
		})
	}
	return &config.Config{
		Newrelic: &config.NewRelicInput{
			LogLevel:       logLevel,
			EventsEndpoint: "",
			LicenseKey:     "",
		},
		Logger:    logging.NewLogger(logLevel),
		Endpoints: eps,
	}
}
