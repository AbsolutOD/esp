/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/pinpt/esp/internal/common"

	"github.com/pinpt/esp/internal/client"
	"github.com/spf13/cobra"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	//Args:  cobra.MinimumNArgs(2),
	Short: "Creates an SSM parameter with the given value.",
	Long: `Simple command to add values to SSM parameter store.`,
	Run: func(cmd *cobra.Command, args []string) {
		ec := client.New(client.EspClient{ Backend: "ssm" })
		path, _ := cmd.Flags().GetString("path")
		name, _ := cmd.Flags().GetString("name")
		secure, _ := cmd.Flags().GetBool("secure")
		value, _ := cmd.Flags().GetString("value")
		param := common.EspParamInput{
			Path:   path,
			Name:   name,
			Secure: secure,
			Value:  value,
		}
		ec.Save(param)

	},
}

func init() {
	rootCmd.AddCommand(putCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// putCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// putCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	putCmd.Flags().StringP("name", "n", "", "The name for your parameter.")
	putCmd.Flags().StringP("path", "p", "", "Define the path for the SSM parameter.")
	putCmd.Flags().StringP("value", "v", "", "The value to be stored in the SSM.")
	putCmd.Flags().BoolP("secure", "s", false, "Sets the SSM parameter type to 'SecureString'.")
}
