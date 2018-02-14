package core

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type SubConfig struct {
	Port    string
	PortSSL string
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

func (c *Config) Validate() {
	if c.Mode != "DEV" && c.Mode != "PROD" {
		panic("`Mode` must be \"DEV\" or \"PROD\"")
	}
	if c.getSubConfig().Port == "" {
		panic("`Port` must be set")
	}
	if c.getSubConfig().PortSSL == "" {
		panic("`PortSSL` must be set")
	}
	if _, err := strconv.Atoi(c.getSubConfig().Port); err != nil {
		panic("`Port` must be a number")
	} else {
		log.Printf("Port: %v\n", c.getSubConfig().Port)
	}
	if _, err := strconv.Atoi(c.getSubConfig().PortSSL); err != nil {
		panic("`PortSSL` must be a number")
	} else {
		log.Printf("PortSSL: %v\n", c.getSubConfig().PortSSL)
	}
}

func (c *Config) getSubConfig() SubConfig {
	if c.Mode == "PROD" {
		return c.Prod
	}
	return c.Dev
}
