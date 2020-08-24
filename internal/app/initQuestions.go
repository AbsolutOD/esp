package app

import (
	"github.com/AlecAivazis/survey/v2"
)

type configInput struct {
	OrgName string
	AppName string
	Envs    string
}

/*func (ei *envInput) WriteAnswer(name string, value interface{}) error {
	ei.name = name
	ei.value = strings.Split(value.(string), "\n")
	//fmt.Printf("ENVS: %s", ei.value)
	return errors.New("Couldn't parse value.")
}*/

// the questions to ask
var qs = []*survey.Question{
	{
		Name:     "orgName",
		Prompt:   &survey.Input{Message: "What is your Org's name?"},
		Validate: survey.Required,
	},
	{
		Name: "appName",
		Prompt:   &survey.Input{Message: "What is the name of your app?"},
		Validate: survey.Required,
	},
	{
		Name: "Envs",
		Prompt: &survey.Input{
			Message: "What environments do you have? Separate with  `,`.",
			Help: "Separate envs with a ','.",
		},
		Validate: survey.Required,
	},
}
