package cmd

import (
	"fmt"
	"github.com/pinpt/esp/internal/common"

	"github.com/spf13/cobra"
)

// putCmd stores the parameter in the backend store
var putCmd = &cobra.Command{
	Use:   "put",
	Aliases: []string{"add", "create"},
	Short: "Creates an SSM parameter with the given value.",
	Long: `Simple command to add values to SSM parameter store.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		secure, _ := cmd.Flags().GetBool("secure")
		value, _ := cmd.Flags().GetString("value")
		param := common.EspParamInput{
			Name:   name,
			Secure: secure,
			Value:  value,
		}
		c.Save(param)
		savedParam := c.GetParam(common.GetOneInput{
			Name: name,
		})
		detailDisplay(savedParam)
	},
}

func init() {
	rootCmd.AddCommand(putCmd)
	putCmd.Flags().StringP("name", "n", "", "The name for your parameter.")
	putCmd.Flags().StringP("value", "v", "", "The value to be stored in the SSM.")
	putCmd.Flags().BoolP("secure", "s", false, "Sets the SSM parameter type to 'SecureString'.")
	err := cobra.MarkFlagRequired(putCmd.Flags(), "name")
	if err != nil {
		fmt.Print("can't set flag --name as required")
	}
	err = cobra.MarkFlagRequired(putCmd.Flags(), "value")
	if err != nil {
		fmt.Print("can't set flag --name as required")
	}
}
