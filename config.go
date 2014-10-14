package bur

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type serverConfig struct {
	Addr   string
	Params map[string]string
}

type config struct {
	Proxy  map[string]serverConfig
	Auth   string
	Debug  bool
	Access string
}

var _config config

func ReadConfig(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &_config); err != nil {
		return err
	}
	return nil
}
