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
	"github.com/pinpt/esp/pkg/client"
	"github.com/pinpt/esp/pkg/errors"

	"github.com/logrusorgru/aurora"

	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/spf13/cobra"
)


// getParam Queries the ssm param
func getParam(ec client.EspConfig, key string) []*ssm.Parameter {
	si := &ssm.GetParameterInput{
		Name: &key,
	}
	param, err := ec.Client.GetParameter(si)
	if err != nil {
		errors.CheckSSMError(err)
	}

	return [param.Parameter]
}

func getParamsByPath(ec client.EspConfig, path string) *ssm.GetParametersByPathOutput {
	si := &ssm.GetParametersByPathInput{
		Path: aws.String(path),
	}
	params, err := ec.Client.GetParametersByPath(si)
	if err != nil {
		errors.CheckSSMByPath(err)
	}
	return params
}


// getCmd represents the path command
var getCmd = &cobra.Command{
	Use:   "get [path]",
	Short: "Query path for SSM",
	Long:  `Allows you to get a specific ssm parameter with an exact path.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ec := client.New("us-east-1")

		//fmt.Printf("Getting: %v\n", args[0])
		if flag, _ := cmd.Flags().GetBool("recursive"); flag {
			getParamsByPath(ec, args[0])
		} else {
			getParam(ec, args[0])
		}

		keystr := aurora.BrightYellow(args[0])
		fmt.Printf("%s: %s\n", keystr, *param.Parameter.Value)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	getCmd.Flags().BoolP("recursive", "r", false, "Does a recursive query given a base path.")
	//getCmd.Flags().BoolP("recursive", "r", false, "Does a recursive query given a base path.")
}
