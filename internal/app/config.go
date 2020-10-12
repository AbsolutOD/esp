package app

import (
	"fmt"
	"github.com/pinpt/esp/internal/utils"
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

func (c Config) GetAppPath() string {
	path := fmt.Sprintf("/%s/%s/%s/", c.OrgName, c.Env, c.AppName)
	return path
}

func (c *Config) CheckEnv(env string) bool {
	for _, e := range c.Envs {
		if env == e {
			return true
		}
	}
	return false
}
