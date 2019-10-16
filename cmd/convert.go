/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

type EspVar struct {
	Name  string
	Value string
	Path  string
}

type TaskEnvVars struct {
	Environment  []TaskEnvVar    `json:"environment"`
	Secrets      []TaskSecretVar `json:"secrets"`
}

type TaskEnvVar struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

type TaskSecretVar struct {
	Name  string `json:"name"`
	Path  string `json:"valueFrom"`
	Value string
}

func readTaskJson(p string) TaskEnvVars {
	jsonfile, err := os.Open(p)
	if err != nil {
		fmt.Printf("%v %v", aurora.Red("Couldn't open: "), p)
	}
	//fmt.Printf("Successfully Opened %s\n", p)
	defer jsonfile.Close()

	byteValue, _ := ioutil.ReadAll(jsonfile)

	var EnvSecretsVars TaskEnvVars

	json.Unmarshal(byteValue, &EnvSecretsVars)

	fmt.Println(EnvSecretsVars)

	return EnvSecretsVars
}

func buildEspVars(ev TaskEnvVars) map[string]EspVar {
	var espvars = make(map[string]EspVar)

	for _, taskvar := range ev.Environment {
		espvar := EspVar{
			Name: taskvar.Name,
			Value: taskvar.Value,
		}
		espvars[taskvar.Name] = espvar
		fmt.Printf("%s: \"%s\"\n", taskvar.Name, taskvar.Value)
	}

	for _, taskvar := range ev.Secrets {
		espvar := EspVar{
			Name: taskvar.Name,
			Path: taskvar.Path,
		}
		espvars[taskvar.Name] = espvar
		fmt.Printf("%s: \"%s\"\n", taskvar.Name, taskvar.Path)
	}

	return espvars
}

func writeJsonFile(ev map[string]EspVar)  {
	file, err := json.MarshalIndent(ev, "", "  ")
	if err != nil {
		fmt.Println("Error encoding json.")
	}
	err = ioutil.WriteFile("espfile.json", file, 0644)
	if err != nil {
		fmt.Println("Error trying to write file.")
	}
}

// convertCmd represents the export command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert the environment and secrets keys to ESP format",
	Long: `If you remove all of the other top level keys from an AWS ECS task definition json object 
and just reformat so the "environment" and "secret" keys are at the top level.  Then this script
will reformat it, so that the ESP tool can read it.
`,
	Run: func(cmd *cobra.Command, args []string) {
		//ec := client.New("us-east-1")
		fmt.Println("export called")
		jsonvars := readTaskJson(args[0])
		espvars := buildEspVars(jsonvars)
		writeJsonFile(espvars)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
