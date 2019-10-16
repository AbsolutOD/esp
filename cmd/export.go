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
	"github.com/pinpt/esp/pkg/ssmparam"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

type EspFile struct {
	Static []ConfigMapVar `json:"static"`
	Ssm    []ConfigMapVar `json:"ssm"`
}

type ConfigMapVar struct {
	Name string    `json:"Name"`
	Value string   `json:"Value,omitempty"`
	SsmPath string `json:"Path,omitempty"`
}

func readEspJson(p string) EspFile {
	jsonfile, err := os.Open(p)
	if err != nil {
		fmt.Printf("%v %v", aurora.Red("Couldn't open: "), p)
	}
	//fmt.Printf("Successfully Opened %s\n", p)
	defer jsonfile.Close()

	byteValue, _ := ioutil.ReadAll(jsonfile)

	var espfile EspFile

	err = json.Unmarshal(byteValue, &espfile)
	if err != nil {
		fmt.Println("Error parsing json.")
		fmt.Println(err)
	}
	//fmt.Println(espfile)

	return espfile
}

//  TODO: don't need this at the moment
//func buildVarsFromTaskEnvVars()  {
//}

// TODO: need to only pass 10 paths at a time.
/*func getMultipleParamsFromSsm(ec client.EspConfig, s []ConfigMapVar) []ConfigMapVar {
	var paths []*string
	var pIndex = make(map[string]int)
	var configMap []ConfigMapVar
	for i, cmv := range s {
		pIndex[cmv.SsmPath] = i
		paths = append(paths, &cmv.SsmPath)
	}
	//fmt.Println(paths)
	params := ssmparam.GetMany(ec, true, paths)
	//fmt.Println(params)

	for _, param := range params {
		fmt.Println(*param.Name)
		i := pIndex[*param.Name]
		s[i].Value = *param.Value
	}
	return configMap
}*/

func getParamsFromSsm(ec client.EspConfig, s []ConfigMapVar) []ConfigMapVar {
	var ssmparams []ConfigMapVar
	var pIndex = make(map[string]int)
	for i, cmv := range s {
		pIndex[cmv.SsmPath] = i
		param := ssmparam.GetOne(ec, true, cmv.SsmPath)
		cmv.Value = *param.Value
		ssmparams = append(ssmparams, cmv)
	}

	return ssmparams
}

func printVars(cm *[]string) {
	for _, cmstr := range *cm {
		fmt.Println(cmstr)
	}
}

func addConfigMapVars(cm []string, s []ConfigMapVar) []string {

	for _, cmv := range s {
		cmvs := fmt.Sprintf("%s: %s", cmv.Name, cmv.Value)
		cm = append(cm, cmvs)
	}
	return cm
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
		var configMap []string
		ec := client.New("us-east-1")
		//fmt.Println("export called")
		espfile := readEspJson(args[0])
		configMap = addConfigMapVars(configMap, espfile.Static)
		ssmparams := getParamsFromSsm(ec, espfile.Ssm)
		//fmt.Println(ssmvars)
		configMap = addConfigMapVars(configMap, ssmparams)
		fmt.Println(configMap)
		printVars(&configMap)
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
