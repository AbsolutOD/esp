package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/pinpt/esp/internal/client"
	"github.com/pinpt/esp/internal/ssmparam"
	"github.com/spf13/cobra"
)

type EspFile struct {
	Static []ConfigMapVar `json:"static"`
	Ssm    []ConfigMapVar `json:"ssm"`
}

type ConfigMapVar struct {
	Name    string `json:"Name"`
	Value   string `json:"Value,omitempty"`
	SsmPath string `json:"Path,omitempty"`
}

func readEspJson(p string) EspFile {
	jsonfile, err := os.Open(p)
	if err != nil {
		panic(fmt.Sprintf("%v %v", aurora.Red("Couldn't open: "), p))
	}
	//fmt.Printf("Successfully Opened %s\n", p)
	defer jsonfile.Close()

	byteValue, _ := ioutil.ReadAll(jsonfile)

	var espfile EspFile

	err = json.Unmarshal(byteValue, &espfile)
	if err != nil {
		panic(fmt.Sprintf("Error parsing json. %w", err))
	}
	return espfile
}

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
	Run: func(cmd *cobra.Command, args []string) {
		var configMap []string
		region, _ := cmd.Flags().GetString("region")
		ec := client.New(region)
		espfile := readEspJson(args[0])
		configMap = addConfigMapVars(configMap, espfile.Static)
		ssmparams := getParamsFromSsm(ec, espfile.Ssm)
		configMap = addConfigMapVars(configMap, ssmparams)
		printVars(&configMap)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
