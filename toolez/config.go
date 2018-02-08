package toolez

import "log"

type Config struct {
	HostNames []string
}

func (c *Config) Validate() {
	for _, h := range c.HostNames {
		if h == "" {
			panic("`HostNames` must not contain empty string")
		}
	}
	log.Printf("HostNames: %v", c.HostNames)
}
