package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Endpoint struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Config struct {
	Newrelic struct {
		LogLevel string `yaml:"logLevel"`
	} `yaml:"newrelic"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

func NewConfig() (
	*Config,
	error,
) {

	configPath := os.Getenv("CONFIG_PATH")
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var cfg Config

	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
