package app

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	IsEspFile   bool
	OrgName     string
	AppName     string
	Envs        []string
	Path        string
}

// New creates a new config object
func New(isEspFile bool) Config {
	cfg := Config{
		IsEspFile: isEspFile,
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("unable to decode into struct, %v\n", err)
	}
	return cfg
}

func (c *Config) CheckEnv(env string) bool {
	for _, e := range c.Envs {
		if env == e {
			return true
		}
	}
	return false
}
