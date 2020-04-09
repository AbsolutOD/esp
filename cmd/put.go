package cmd

import (
	"github.com/pinpt/esp/internal/common"

	"github.com/pinpt/esp/internal/client"
	"github.com/spf13/cobra"
)

// putCmd stores the parameter in the backend store
var putCmd = &cobra.Command{
	Use:   "put",
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
	putCmd.Flags().StringP("name", "n", "", "The name for your parameter.")
	putCmd.Flags().StringP("path", "p", "", "Define the path for the SSM parameter.")
	putCmd.Flags().StringP("value", "v", "", "The value to be stored in the SSM.")
	putCmd.Flags().BoolP("secure", "s", false, "Sets the SSM parameter type to 'SecureString'.")
}
