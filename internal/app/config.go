package app

import (
	"fmt"
	"github.com/pinpt/esp/internal/utils"
	jww "github.com/spf13/jwalterweatherman"
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

// GetAppPath returns the application's base path from the config
func (c Config) GetAppPath() string {
	path := fmt.Sprintf("/%s/%s/%s/", c.OrgName, c.Env, c.AppName)
	jww.INFO.Printf("rendered path: %s", path)
	return path
}

// GetAppParamPath builds the application's param path from the config file
func (c Config) GetAppParamPath(p string) string {
	path := fmt.Sprintf("/%s/%s/%s/%s", c.OrgName, c.Env, c.AppName, p)
	jww.INFO.Printf("rendered param path: %s", path)
	return path
}

func (c *Config) CheckEnv(env string) bool {
	for _, e := range c.Envs {
		if env == e {
			jww.INFO.Printf("found ENV: %s", e)
			return true
		}
	}
	return false
}
