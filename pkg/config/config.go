package config

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Endpoint struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Config struct {
	Newrelic struct {
		LogLevel       string `yaml:"logLevel"`
		EventsEndpoint string
		LicenseKey     string
	} `yaml:"newrelic"`
	Endpoints []Endpoint `yaml:"endpoints"`
	Logger    *Logger
}

func NewConfig() *Config {

	// Parse config file
	cfg := parseConfigFile()

	// Create logger
	cfg.Logger = NewLogger(cfg.Newrelic.LogLevel)

	// Parse New Relic license key
	cfg.Newrelic.LicenseKey = parseNewRelicLicenseKey()
	cfg.Newrelic.EventsEndpoint = setNewRelicEventsEndpoint(cfg.Newrelic.LicenseKey)

	cfg.Logger.Log(logrus.DebugLevel, "Config file is succesfully created.")
	return &cfg
}

func parseConfigFile() Config {

	// Get & check config path
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("Config path is empty!")
	}

	// Read config file
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	// Parse config file
	var cfg Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func parseNewRelicLicenseKey() string {
	nrLicenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if nrLicenseKey == "" {
		panic("License key is not provided! Define config.data.newrelic.licenseKey in your Helm deployment.")
	}

	return nrLicenseKey
}

func setNewRelicEventsEndpoint(
	licenseKey string,
) string {

	nrAccountId := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	if nrAccountId == "" {
		panic("Account ID not provided! Define config.data.newrelic.accountId in your Helm deployment.")
	}

	if licenseKey[0:2] == "eu" {
		return "https://insights-collector.eu01.nr-data.net/v1/accounts/" + nrAccountId + "/events"
	} else {
		return "https://insights-collector.nr-data.net/v1/accounts/" + nrAccountId + "/events"
	}
}
