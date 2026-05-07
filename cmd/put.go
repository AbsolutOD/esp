package cmd

import (
	"fmt"
	"strings"

	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
)

func formatParamName(n string) string {
	if strings.HasPrefix(n, esp.OrgPrefix) {
		return n
	}
	// SSM allows hyphens in parameter names, but esp's purpose is to
	// surface params as environment variables, where hyphens are not
	// valid. Normalize them to underscores alongside the upcasing.
	normalized := strings.ReplaceAll(strings.ToUpper(n), "-", "_")
	return esp.OrgPrefix + "_" + normalized
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
	Use:     "put",
	Aliases: []string{"add", "create"},
	Short:   "Creates an SSM parameter with the given value.",
	Long:    `Simple command to add values to SSM parameter store.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		param := buildEspParamInputFromCmd(cmd)
		if _, err := c.Save(param); err != nil {
			return err
		}
		savedParam, err := c.GetParam(common.GetOneInput{
			Name: param.Name,
		})
		if err != nil {
			return err
		}
		detailDisplay(savedParam)
		return nil
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
