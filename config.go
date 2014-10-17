package bur

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type configService struct {
	Addr   string
	Params map[string]string
}

type configLog struct {
	File, Level string
}

type Config struct {
	Proxy map[string]configService
	Auth  configService
	Debug bool
	State configService
	Log   configLog
}

func LoadConfig(filename string) (config *Config, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &config)
	return
}
