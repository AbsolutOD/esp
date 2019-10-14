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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/olekukonko/tablewriter"
	"github.com/pinpt/esp/pkg/client"
	"github.com/pinpt/esp/pkg/errors"
	"os"

	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/spf13/cobra"
)

// getParam Queries the ssm param
func getParam(ec client.EspConfig, d bool, key string) *ssm.Parameter {
	si := &ssm.GetParameterInput{
		Name: aws.String(key),
		WithDecryption: aws.Bool(d),
	}
	resp, err := ec.Svc.GetParameter(si)
	if err != nil {
		errors.CheckSSMGetParameters(err)
	}

	return resp.Parameter
}

func displayParam(p *ssm.Parameter) {
	data := [][]string{
		[]string{"ARN", *p.ARN},
		[]string{"Last Modified Date", p.LastModifiedDate.String()},
		[]string{"Name", *p.Name},
		[]string{"Type", *p.Type},
		[]string{"Value", *p.Value},
		[]string{"Version", string(*p.Version)},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Keys", "Value"})
	table.AppendBulk(data)
	table.Render()
}

// getCmd represents the path command
var getCmd = &cobra.Command{
	Use:   "get [path]",
	Short: "Query path for SSM",
	Long:  `Allows you to get a specific ssm parameter with an exact path or recursively get params.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Printf("Got: %s\n", args[0])
		ec := client.New("us-east-1")
		decrypt, _ := cmd.Flags().GetBool("decrypt")

		param := getParam(ec, decrypt, args[0])
		displayParam(param)
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

	getCmd.Flags().BoolP("decrypt", "d", false, "Decrypt SSM secure strings.")
}
