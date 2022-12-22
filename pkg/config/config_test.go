package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func Test_NewRelicLicenseKeyIsNotDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return ""
	}
	assert.Panics(t, func() { parseNewRelicLicenseKey() })
}

func Test_NewRelicLicenseKeyIsDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return "LICENSE_KEY"
	}
	assert.NotPanics(t, func() { parseNewRelicLicenseKey() })
}

func Test_NewRelicAccountIdIsNotDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return ""
	}
	assert.Panics(t, func() { setNewRelicEventsEndpoint("") })
}

func Test_NewRelicEndpointIsEu(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	nrAccountId := "ACCOUNT_ID"
	getEnv = func(string) string {
		return nrAccountId
	}
	assert.Equal(t, "https://insights-collector.eu01.nr-data.net/v1/accounts/"+nrAccountId+"/events", setNewRelicEventsEndpoint("eu_LICENSE_KEY"))
}

func Test_NewRelicEndpointIsUs(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	nrAccountId := "ACCOUNT_ID"
	getEnv = func(string) string {
		return nrAccountId
	}
	assert.Equal(t, "https://insights-collector.nr-data.net/v1/accounts/"+nrAccountId+"/events", setNewRelicEventsEndpoint("LICENSE_KEY"))
}

func Test_ConfigFilePathIsNotDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return ""
	}
	assert.Panics(t, func() { parseConfigFile() })
}

func Test_ConfigFileIsNotDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return "CONFIG_PATH"
	}
	assert.Panics(t, func() { parseConfigFile() })
}

func Test_ConfigFileHasInvalidYamlFormat(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return "CONFIG_PATH"
	}

	readFileMock := readFile
	defer func() {
		readFile = readFileMock
	}()

	readFile = func(string) ([]byte, error) {
		return []byte("false config"), nil
	}

	assert.Panics(t, func() { parseConfigFile() })
}

func Test_NoEndpointIsDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return "CONFIG_PATH"
	}

	readFileMock := readFile
	defer func() {
		readFile = readFileMock
	}()

	readFile = func(string) ([]byte, error) {
		cfg := &Config{
			Newrelic: &NewRelicInput{
				LogLevel:       "ERROR",
				EventsEndpoint: "",
				LicenseKey:     "",
			},
			Logger: nil,
		}

		bytes, err := yaml.Marshal(cfg)
		if err != nil {
			t.Log(err)
		}

		return bytes, nil
	}

	assert.Panics(t, func() { parseConfigFile() })
}

func Test_EndpointTypeIsNotDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return "CONFIG_PATH"
	}

	readFileMock := readFile
	defer func() {
		readFile = readFileMock
	}()

	readFile = func(string) ([]byte, error) {
		cfg := &Config{
			Newrelic: &NewRelicInput{
				LogLevel:       "ERROR",
				EventsEndpoint: "",
				LicenseKey:     "",
			},
			Logger: nil,
			Endpoints: []Endpoint{
				{
					Name: "Name",
					URL:  "URL",
				},
			},
		}

		bytes, err := yaml.Marshal(cfg)
		if err != nil {
			t.Log(err)
		}

		return bytes, nil
	}

	assert.Panics(t, func() { parseConfigFile() })
}

func Test_EndpointNameIsNotDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return "CONFIG_PATH"
	}

	readFileMock := readFile
	defer func() {
		readFile = readFileMock
	}()

	readFile = func(string) ([]byte, error) {
		cfg := &Config{
			Newrelic: &NewRelicInput{
				LogLevel:       "ERROR",
				EventsEndpoint: "",
				LicenseKey:     "",
			},
			Logger: nil,
			Endpoints: []Endpoint{
				{
					Type: "kvp",
					URL:  "URL",
				},
			},
		}

		bytes, err := yaml.Marshal(cfg)
		if err != nil {
			t.Log(err)
		}

		return bytes, nil
	}

	assert.Panics(t, func() { parseConfigFile() })
}

func Test_EndpointUrlIsNotDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return "CONFIG_PATH"
	}

	readFileMock := readFile
	defer func() {
		readFile = readFileMock
	}()

	readFile = func(string) ([]byte, error) {
		cfg := &Config{
			Newrelic: &NewRelicInput{
				LogLevel:       "ERROR",
				EventsEndpoint: "",
				LicenseKey:     "",
			},
			Logger: nil,
			Endpoints: []Endpoint{
				{
					Type: "kvp",
					Name: "Name",
				},
			},
		}

		bytes, err := yaml.Marshal(cfg)
		if err != nil {
			t.Log(err)
		}

		return bytes, nil
	}

	assert.Panics(t, func() { parseConfigFile() })
}

func Test_EndpointTypeIsNotSupported(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return "CONFIG_PATH"
	}

	readFileMock := readFile
	defer func() {
		readFile = readFileMock
	}()

	readFile = func(string) ([]byte, error) {
		cfg := &Config{
			Newrelic: &NewRelicInput{
				LogLevel:       "ERROR",
				EventsEndpoint: "",
				LicenseKey:     "",
			},
			Logger: nil,
			Endpoints: []Endpoint{
				{
					Type: "yaml",
					Name: "Name",
					URL:  "URL",
				},
			},
		}

		bytes, err := yaml.Marshal(cfg)
		if err != nil {
			t.Log(err)
		}

		return bytes, nil
	}

	assert.Panics(t, func() { parseConfigFile() })
}

func Test_ConfigFileIsValid(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return "CONFIG_PATH"
	}

	readFileMock := readFile
	defer func() {
		readFile = readFileMock
	}()

	readFile = func(string) ([]byte, error) {
		cfg := &Config{
			Newrelic: &NewRelicInput{
				LogLevel:       "ERROR",
				EventsEndpoint: "",
				LicenseKey:     "",
			},
			Logger: nil,
			Endpoints: []Endpoint{
				{
					Type: "kvp",
					Name: "Name",
					URL:  "URL",
				},
			},
		}

		bytes, err := yaml.Marshal(cfg)
		if err != nil {
			t.Log(err)
		}

		return bytes, nil
	}

	assert.NotPanics(t, func() { NewConfig() })
}
