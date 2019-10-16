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
	"github.com/pinpt/esp/pkg/client"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

type ConfigMap struct {
	Data []EspVar
}

type ConfigMapVar struct {
	Name string
	Value string
	SsmPath string
}

func readEspJson(p string) TaskEnvVars {
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

func buildVarsFromTaskEnvVars()  {

}

func getSecretsFromSsm(ec client.EspConfig, s []TaskSecretVar) []string {
	var paths []*string
	var pIndex map[string]int
	for i, path := range s {
		pIndex[path.Name] = i
		paths = append(paths, &path.Path)
	}
	/*params := ssmparam.GetMany(ec, true, paths)

	for _, param := range params {

	}*/
	return []string{}
}

func printVars(ev TaskEnvVars) {
	for _, envvar := range ev.Environment {
		fmt.Printf("%s: \"%s\"\n", envvar.Name, envvar.Value)
	}

	for _, envvar := range ev.Secrets {

		fmt.Printf("%s: \"%s\"\n", envvar.Name, envvar.Path)
	}
}

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export SSM Parameters in a format for a ConfigMap",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//ec := client.New("us-east-1")
		fmt.Println("export called")
		envvars := readEspJson(args[0])
		//secretVars := getSecretsFromSsm(envvars.Secrets)
		printVars(envvars)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
