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
	OrgPrefix   string
	AppName     string
	Envs        []string
	Env         string
	Path        string
	Filename    string
}

// New creates a new config object
func New(isEspFile bool) *Config {
	cfg := &Config{
		IsEspProject: isEspFile,
		Filename: ".espFile",
		Path: utils.GetCwd(),
	}
	return cfg
}

func (c *Config) ReadEspFile()  {
	if err := viper.Unmarshal(c); err != nil {
		fmt.Printf("unable to decode into struct, %v\n", err)
	}
}

func (c *Config) CheckEnv(env string) bool {
	for _, e := range c.Envs {
		if env == e {
			return true
		}
	}
	return false
}
