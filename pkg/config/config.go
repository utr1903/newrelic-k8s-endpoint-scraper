package config

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Newrelic struct {
		LogLevel string `yaml:"logLevel"`
	} `yaml:"newrelic"`
	Endpoints []struct {
		Type string `yaml:"type"`
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
	} `yaml:"endpoints"`
}

func New() (
	*Config,
	error,
) {

	filename, _ := filepath.Abs("./config.yaml")
	configFile, err := ioutil.ReadFile(filename)
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
