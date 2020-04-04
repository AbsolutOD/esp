package cmd

import (
	"fmt"
	"github.com/absolutod/esp/internal/common"

	"github.com/absolutod/esp/internal/client"
	"github.com/logrusorgru/aurora"

	"github.com/spf13/cobra"
)

func displayParams(p []common.EspParam) {
	for _, param := range p {
		name := aurora.BrightYellow(param.Name)
		fmt.Printf("%s: %s\n", name, param.Value)
	}
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "Recursively list a SSM path.",
	Long:  `The list command gives you an easy way to recursively get all SSM parameters with a base path.`,
	Run: func(cmd *cobra.Command, args []string) {
		ec := client.New(client.EspClient{Backend: "ssm"})
		decrypt, _ := cmd.Flags().GetBool("decrypt")
		params := ec.ListParams(common.ListParamInput{
			Path:    args[0],
			Decrypt: decrypt,
		})
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
