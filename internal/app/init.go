package app

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
)

const espFileName = ".espFile"
const espFilePath = "./" + espFileName + ".yaml"

type Config struct {
	Loaded      bool
	Prefixes    map[string]string
	V           *viper.Viper
	EspFile     espFile
}

type espFile struct {
	orgName string
	appName string
	envs    []string
	path    string
}

// New creates a new config object
func New() Config {
	return Config{ EspFile: espFile{}}
}

// CheckForConfigFile is just a simple function to see if the espFile exists
func (e *espFile) CheckForConfigFile() bool {
	if _, err := os.Stat(e.path); os.IsNotExist(err) {
		return false
	}
	return true
}

// WriteConfig saves the app config to a file
func (e *espFile) WriteConfig(ci configInput) error {
	cfg := espFile{
		orgName: ci.OrgName,
		appName: ci.AppName,
		envs: regexp.MustCompile(", *").Split(ci.Envs, -1),
	}

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(espFilePath, out, 0660); err != nil {
		return err
	}

	return nil
}

// Init asks the user for properties to put in the espFile
func (c *Config) Init()  {
	ci := configInput{}
	if err := survey.Ask(qs, &ci); err != nil {
		fmt.Printf("There was an error: %s", err)
	}
	//fmt.Printf("Answers: %s", ci)
	err := c.EspFile.WriteConfig(ci)
	if err != nil {
		fmt.Printf("Error writing file: %s", err)
	}
}

func (c *Config) setPrefixes() {

}

func (c *Config) readEspFile()  {
	c.V = viper.New()
	c.V.SetConfigName(espFileName)
	c.V.AddConfigPath(".")
	err := c.V.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading the %s", espFilePath)
	}
	c.Loaded = true
	//c.Envs = c.V.GetStringSlice("envs")
}
