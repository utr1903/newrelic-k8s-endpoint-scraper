package config

import (
	"fmt"
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

type NewRelicInput struct {
	LogLevel       string `default:"ERROR" yaml:"logLevel"`
	EventsEndpoint string
	LicenseKey     string
}

type Config struct {
	Newrelic  *NewRelicInput `yaml:"newrelic"`
	Endpoints []Endpoint     `yaml:"endpoints"`
	Logger    *Logger
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

func NewConfig() *Config {

	// Parse config file
	cfg := parseConfigFile()

	// Create logger
	cfg.Logger = NewLogger(cfg.Newrelic.LogLevel)

	// Parse New Relic license key
	cfg.Newrelic.LicenseKey = parseNewRelicLicenseKey()
	cfg.Newrelic.EventsEndpoint = setNewRelicEventsEndpoint(cfg.Newrelic.LicenseKey)

	cfg.Logger.Log(logrus.DebugLevel, "Config file is succesfully created.")
	return cfg
}

func parseConfigFile() *Config {

	// Get & check config path
	configPath := getEnv("CONFIG_PATH")
	if configPath == "" {
		msg := "Config path is not defined!"
		fmt.Println(msg)
		panic(msg)
	}

	// Read config file
	configFile, err := readFile(configPath)
	if err != nil {
		fmt.Println("Config file could not be read!")
		panic(err)
	}

	// Parse config file
	var cfg Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		fmt.Println("Config file could not be parsed into yaml format!")
		panic(err)
	}

	// Check if endpoints are defined correctly
	checkEndpoints(&cfg)

	return &cfg
}

func parseNewRelicLicenseKey() string {
	nrLicenseKey := getEnv("NEW_RELIC_LICENSE_KEY")
	if nrLicenseKey == "" {
		msg := "License key is not provided! Define config.data.newrelic.licenseKey in your Helm deployment."
		fmt.Println(msg)
		panic(msg)
	}

	return nrLicenseKey
}

func setNewRelicEventsEndpoint(
	licenseKey string,
) string {

	nrAccountId := getEnv("NEW_RELIC_ACCOUNT_ID")
	if nrAccountId == "" {
		msg := "Account ID not provided! Define config.data.newrelic.accountId in your Helm deployment."
		fmt.Println(msg)
		panic(msg)
	}

	if licenseKey[0:2] == "eu" {
		return "https://insights-collector.eu01.nr-data.net/v1/accounts/" + nrAccountId + "/events"
	} else {
		return "https://insights-collector.nr-data.net/v1/accounts/" + nrAccountId + "/events"
	}
}

func checkEndpoints(
	cfg *Config,
) {
	if cfg.Endpoints == nil || len(cfg.Endpoints) == 0 {
		msg := "No endpoint is defined!"
		cfg.Logger.Log(logrus.ErrorLevel, msg)
		panic(msg)
	}

	for _, endpoint := range cfg.Endpoints {
		if endpoint.Type == "" || endpoint.Name == "" || endpoint.URL == "" {
			msg := "Check your endpoint definitions! Type, Name and URL must be defined!"
			cfg.Logger.Log(logrus.ErrorLevel, msg)
			panic(msg)
		}
	}
}
