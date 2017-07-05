package core

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type SubConfig struct {
	Addr string
}

type Config struct {
	Mode     string
	Dev      SubConfig
	Prod     SubConfig
	RootPath string
	Site     map[string]interface{}
}

func (c *Config) ParseFile(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}
	configRaw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal([]byte(configRaw), c)
	c.Mode = strings.ToUpper(c.Mode)
	return c
}
