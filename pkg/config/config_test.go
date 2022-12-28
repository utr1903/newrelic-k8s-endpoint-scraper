package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/logging"
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

	licenseKey, err := parseNewRelicLicenseKey()
	assert.Equal(t, "", licenseKey)
	assert.NotNil(t, err)
	assert.Equal(t, logging.CONFIG__LICENSE_KEY_IS_NOT_PROVIDED, err.Error())
}

func Test_NewRelicLicenseKeyIsDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	licenseKeyExpected := "LICENSE_KEY"
	getEnv = func(string) string {
		return licenseKeyExpected
	}

	licenseKeyActual, err := parseNewRelicLicenseKey()
	assert.Equal(t, licenseKeyExpected, licenseKeyActual)
	assert.Nil(t, err)
}

func Test_NewRelicAccountIdIsNotDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	licenseKey := ""
	accountId := ""
	getEnv = func(string) string {
		return accountId
	}

	eventsEndpoint, err := setNewRelicEventsEndpoint(licenseKey)
	assert.Equal(t, "", eventsEndpoint)
	assert.NotNil(t, err)
	assert.Equal(t, logging.CONFIG__ACCOUNT_ID_IS_NOT_PROVIDED, err.Error())
}

func Test_NewRelicEndpointIsEu(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	licenseKey := "eu_LICENSE_KEY"
	nrAccountId := "ACCOUNT_ID"
	getEnv = func(string) string {
		return nrAccountId
	}

	eventsEndpoint, err := setNewRelicEventsEndpoint(licenseKey)
	assert.Nil(t, err)
	assert.Equal(t,
		"https://insights-collector.eu01.nr-data.net/v1/accounts/"+nrAccountId+"/events",
		eventsEndpoint,
	)
}

func Test_NewRelicEndpointIsUs(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	licenseKey := "us_LICENSE_KEY"
	nrAccountId := "ACCOUNT_ID"
	getEnv = func(string) string {
		return nrAccountId
	}

	eventsEndpoint, err := setNewRelicEventsEndpoint(licenseKey)
	assert.Nil(t, err)
	assert.Equal(t,
		"https://insights-collector.nr-data.net/v1/accounts/"+nrAccountId+"/events",
		eventsEndpoint,
	)
}

func Test_ConfigFilePathIsNotDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return ""
	}

	cfg, err := parseConfigFile()
	assert.Nil(t, cfg)
	assert.NotNil(t, err)
	assert.Equal(t, logging.CONFIG__CONFIG_PATH_IS_NOT_DEFINED, err.Error())
}

func Test_ConfigFileIsNotDefined(t *testing.T) {
	getEnvMock := getEnv
	defer func() {
		getEnv = getEnvMock
	}()

	getEnv = func(string) string {
		return "CONFIG_PATH"
	}

	cfg, err := parseConfigFile()
	assert.Nil(t, cfg)
	assert.NotNil(t, err)
	assert.Equal(t, logging.CONFIG__CONFIG_FILE_COULD_NOT_BE_READ, err.Error())
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

	cfg, err := parseConfigFile()
	assert.Nil(t, cfg)
	assert.NotNil(t, err)
	assert.Equal(t, logging.CONFIG__CONFIG_FILE_COULD_NOT_BE_PARSED_INTO_YAML, err.Error())
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

	cfg, err := parseConfigFile()
	assert.Nil(t, cfg)
	assert.NotNil(t, err)
	assert.Equal(t, logging.CONFIG__NO_ENDPOINT_IS_DEFINED, err.Error())
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

	cfg, err := parseConfigFile()
	assert.Nil(t, cfg)
	assert.NotNil(t, err)
	assert.Equal(t, logging.CONFIG__ENDPOINT_INFO_IS_MISSING, err.Error())
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

	cfg, err := parseConfigFile()
	assert.Nil(t, cfg)
	assert.NotNil(t, err)
	assert.Equal(t, logging.CONFIG__ENDPOINT_INFO_IS_MISSING, err.Error())
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

	cfg, err := parseConfigFile()
	assert.Nil(t, cfg)
	assert.NotNil(t, err)
	assert.Equal(t, logging.CONFIG__ENDPOINT_INFO_IS_MISSING, err.Error())
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

	cfg, err := parseConfigFile()
	assert.Nil(t, cfg)
	assert.NotNil(t, err)
	assert.Equal(t, logging.CONFIG__ENDPOINT_TYPE_IS_NOT_SUPPORTED, err.Error())
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
