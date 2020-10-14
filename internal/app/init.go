package app

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
)

type espFile struct {
	Backend     string
	OrgName     string
	OrgPrefix   string
	AppName     string
	Envs        []string
}

// WriteConfig saves the app config to a file
func (c *Config) WriteConfig() error {
	espFile := c.createEspFile()
	out, err := yaml.Marshal(espFile)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(c.Path, out, 0660); err != nil {
		return err
	}
	return nil
}

// UpdateWithInput updates the esp config struct with answers from the survey questions
func (c *Config) UpdateWithInput(ci configInput) {
	c.OrgName = ci.OrgName
	c.OrgPrefix = ci.OrgPrefix
	c.AppName = ci.AppName
	c.Backend = ci.Backend
	c.Envs = regexp.MustCompile(", *").Split(ci.Envs, -1)
}

// createEspFile takes configInput and returns an espFile struct
func (c Config) createEspFile() espFile {
	espFile := espFile{
		Backend: c.Backend,
		OrgName: c.OrgName,
		OrgPrefix: c.OrgPrefix,
		AppName: c.AppName,
		Envs: c.Envs,
	}
	return espFile
}
// Init asks the user for properties to put in the espFile
func (c *Config) InitQuestions() {
	ci := configInput{}
	if err := survey.Ask(qs, &ci); err != nil {
		fmt.Printf("There was an error: %s", err)
	}
	c.UpdateWithInput(ci)

	err := c.WriteConfig()
	if err != nil {
		fmt.Printf("Error writing file: %s", err)
	}
}

func (c *Config) setPrefixes() {

}
