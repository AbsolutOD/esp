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
	"fmt"

	"github.com/pinpt/esp/internal/client"
	"github.com/spf13/cobra"
)

//var copyDescription = "Copy SSM value from an existing path to a new path.\n"

type CopyCommand struct {
	source      string
	destination string
}

func (cc *CopyCommand) runCopy(region string) {
	ec := client.New(region)
	
}

// cpCmd represents the cp command
var copyCmd = &cobra.Command{
	Use:   "cp [OPTIONS] SRC_SSM_PATH DEST_SSM_PATH",
	Aliases: []string{"cp"},
	Short: "Copy a SSM Param from its current path to a new SSM Path",
	Long: "Copy SSM value from an existing path to a new path.\n",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		fmt.Println("cp called")
		if args[0] == "" {
			return fmt.Errorf("Source can not be empty")
		}
		if args[1] == "" {
			return fmt.Errorf("Destination can not be empty")
		}
		cc := &CopyCommand{
			source:     args[0],
			destination: args[1],
		}
		return cc.runCopy(region)
	},
	Example: "esp cp /ssm/path/key /ssm/new/path/key",
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
