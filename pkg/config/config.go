package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v2"

	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/logging"
)

type Endpoint struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type NewRelicInput struct {
	LogLevel       string `default:"ERROR" yaml:"logLevel"`
	EventsEndpoint string
	LicenseKey     string
}

type Config struct {
	Newrelic  *NewRelicInput `yaml:"newrelic"`
	Endpoints []Endpoint     `yaml:"endpoints"`
	Logger    *logging.Logger
}

var getEnv = func(
	name string,
) string {
	return os.Getenv(name)
}

var readFile = func(
	path string,
) (
	[]byte,
	error,
) {
	return ioutil.ReadFile(path)
}

func NewConfig() (
	*Config,
	error,
) {

	// Parse config file
	cfg, err := parseConfigFile()
	if err != nil {
		return nil, err
	}

	// Parse New Relic license key
	licenseKey, err := parseNewRelicLicenseKey()
	if err != nil {
		return nil, err
	}
	cfg.Newrelic.LicenseKey = licenseKey

	// Set New Relic events endpoint
	eventsEndpoint, err := setNewRelicEventsEndpoint(cfg.Newrelic.LicenseKey)
	if err != nil {
		return nil, err
	}
	cfg.Newrelic.EventsEndpoint = eventsEndpoint

	cfg.Logger.Log(logrus.DebugLevel, "Config file is succesfully created.")
	return cfg, nil
}

func parseConfigFile() (
	*Config,
	error,
) {

	// Get & check config path
	configPath := getEnv("CONFIG_PATH")
	if configPath == "" {
		fmt.Println(logging.CONFIG__CONFIG_PATH_IS_NOT_DEFINED)
		return nil, errors.New(logging.CONFIG__CONFIG_PATH_IS_NOT_DEFINED)
	}

	// Read config file
	configFile, err := readFile(configPath)
	if err != nil {
		fmt.Println(logging.CONFIG__CONFIG_FILE_COULD_NOT_BE_READ)
		return nil, errors.New(logging.CONFIG__CONFIG_FILE_COULD_NOT_BE_READ)
	}

	// Parse config file
	var cfg Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		fmt.Println(logging.CONFIG__CONFIG_FILE_COULD_NOT_BE_PARSED_INTO_YAML)
		return nil, errors.New(logging.CONFIG__CONFIG_FILE_COULD_NOT_BE_PARSED_INTO_YAML)
	}

	// Create logger
	cfg.Logger = logging.NewLogger(cfg.Newrelic.LogLevel)

	// Check if endpoints are defined correctly
	err = checkEndpoints(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseNewRelicLicenseKey() (
	string,
	error,
) {
	nrLicenseKey := getEnv("NEW_RELIC_LICENSE_KEY")
	if nrLicenseKey == "" {
		fmt.Println(logging.CONFIG__LICENSE_KEY_IS_NOT_PROVIDED)
		return "", errors.New(logging.CONFIG__LICENSE_KEY_IS_NOT_PROVIDED)
	}

	return nrLicenseKey, nil
}

func setNewRelicEventsEndpoint(
	licenseKey string,
) (
	string,
	error,
) {

	nrAccountId := getEnv("NEW_RELIC_ACCOUNT_ID")
	if nrAccountId == "" {
		fmt.Println(logging.CONFIG__ACCOUNT_ID_IS_NOT_PROVIDED)
		return "", errors.New(logging.CONFIG__ACCOUNT_ID_IS_NOT_PROVIDED)
	}

	if licenseKey[0:2] == "eu" {
		return "https://insights-collector.eu01.nr-data.net/v1/accounts/" + nrAccountId + "/events", nil
	} else {
		return "https://insights-collector.nr-data.net/v1/accounts/" + nrAccountId + "/events", nil
	}
}

func checkEndpoints(
	cfg *Config,
) error {
	if cfg.Endpoints == nil || len(cfg.Endpoints) == 0 {
		cfg.Logger.Log(logrus.ErrorLevel, logging.CONFIG__NO_ENDPOINT_IS_DEFINED)
		return errors.New(logging.CONFIG__NO_ENDPOINT_IS_DEFINED)
	}

	for _, endpoint := range cfg.Endpoints {
		if endpoint.Type == "" || endpoint.Name == "" || endpoint.URL == "" {
			cfg.Logger.Log(logrus.ErrorLevel, logging.CONFIG__ENDPOINT_INFO_IS_MISSING)
			return errors.New(logging.CONFIG__ENDPOINT_INFO_IS_MISSING)
		}

		if endpoint.Type != "kvp" {
			cfg.Logger.Log(logrus.ErrorLevel, logging.CONFIG__ENDPOINT_TYPE_IS_NOT_SUPPORTED)
			return errors.New(logging.CONFIG__ENDPOINT_TYPE_IS_NOT_SUPPORTED)
		}
	}

	return nil
}
