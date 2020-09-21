package app

import (
	"fmt"
	"github.com/pinpt/esp/internal/utils"
	"github.com/spf13/viper"
)

type Config struct {
	IsEspProject   bool
	Backend     string
	OrgName     string
	AppName     string
	Envs        []string
	Env         string
	Path        string
}

// New creates a new config object
func New(isEspFile bool) Config {
	cfg := Config{
		IsEspProject: isEspFile,
		Path: utils.GetCwd() + "/.espFile.yaml",
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
