package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

type EspVar struct {
	Name  string
	Value string
	Path  string
}

type TaskEnvVars struct {
	Environment []TaskEnvVar    `json:"environment"`
	Secrets     []TaskSecretVar `json:"secrets"`
}

type TaskEnvVar struct {
	Name  string `json:"name"`
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
		panic(fmt.Sprintf("%v %v", aurora.Red("Couldn't open: "), p))
	}
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
			Name:  taskvar.Name,
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

func writeJsonFile(ev map[string]EspVar) {
	file, err := json.MarshalIndent(ev, "", "  ")
	if err != nil {
		panic("Error encoding json.")
	}
	err = ioutil.WriteFile("espfile.json", file, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error trying to write file. %w", err))
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
		jsonvars := readTaskJson(args[0])
		espvars := buildEspVars(jsonvars)
		writeJsonFile(espvars)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

}
