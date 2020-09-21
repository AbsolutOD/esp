package app

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
)


// WriteConfig saves the app config to a file
func (c *Config) WriteConfig(ci *configInput) error {
	c.OrgName = ci.OrgName
	c.AppName = ci.AppName
	c.Envs = regexp.MustCompile(", *").Split(ci.Envs, -1)
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(c.Path, out, 0660); err != nil {
		return err
	}
	return nil
}

// Init asks the user for properties to put in the espFile
func (c *Config) InitQuestions()  {
	ci := configInput{}
	if err := survey.Ask(qs, &ci); err != nil {
		fmt.Printf("There was an error: %s", err)
	}
	//fmt.Printf("Answers: %s", ci)
	err := c.WriteConfig(&ci)
	if err != nil {
		fmt.Printf("Error writing file: %s", err)
	}
}

func (c *Config) setPrefixes() {

}
