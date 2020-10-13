package cmd

import (
	"fmt"
	"github.com/pinpt/esp/internal/common"
	"strings"

	"github.com/spf13/cobra"
)

func formatParamName(n string) string {
	if strings.HasPrefix(n, esp.OrgPrefix) {
		return n
	}
	return esp.OrgPrefix + "_" + strings.ToUpper(n)
}

func getFullPath(n string) string {
	if strings.HasPrefix(n, "/") {
		return n
	}
	name := formatParamName(n)
	return esp.GetAppParamPath(name)
}

func buildEspParamInputFromCmd(cmd *cobra.Command) common.EspParamInput {
	name, _ := cmd.Flags().GetString("name")
	secure, _ := cmd.Flags().GetBool("secure")
	value, _ := cmd.Flags().GetString("value")
	fullName := getFullPath(name)
	param := common.EspParamInput{
		Name:   fullName,
		Secure: secure,
		Value:  value,
	}
	return param
}

// putCmd stores the parameter in the backend store
var putCmd = &cobra.Command{
	Use:   "put",
	Aliases: []string{"add", "create"},
	Short: "Creates an SSM parameter with the given value.",
	Long: `Simple command to add values to SSM parameter store.`,
	Run: func(cmd *cobra.Command, args []string) {
		param := buildEspParamInputFromCmd(cmd)
		c.Save(param)
		savedParam := c.GetParam(common.GetOneInput{
			Name: param.Name,
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
