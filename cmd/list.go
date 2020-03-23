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
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/logrusorgru/aurora"
	"github.com/pinpt/esp/internal/client"
	"github.com/pinpt/esp/internal/errors"

	"github.com/spf13/cobra"
)

func getParamsByPath(ec client.EspConfig, d bool, path string) []*ssm.Parameter {
	si := &ssm.GetParametersByPathInput{
		Path: aws.String(path),
		WithDecryption: aws.Bool(d),
	}
	params, err := ec.Svc.GetParametersByPath(si)
	if err != nil {
		errors.CheckSSMByPath(err)
	}
	return params.Parameters
}

func displayParams(p []*ssm.Parameter) {
	for _, param := range p  {
		name := aurora.BrightYellow(*param.Name)
		fmt.Printf("%s: %s\n", name, *param.Value)
	}
}


// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "Recursively list a SSM path.",
	Long: `The list command gives you an easy way to recursively get all SSM parameters with a base path.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("list called")
		ec := client.New("us-east-1")
		var params []*ssm.Parameter
		decrypt, _ := cmd.Flags().GetBool("decrypt")

		params = getParamsByPath(ec, decrypt, args[0])
		displayParams(params)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	listCmd.Flags().BoolP("decrypt", "d", false, "Decrypt SSM secure strings.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
