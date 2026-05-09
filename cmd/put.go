package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
)

// formatParamName ensures the param ends up as a valid env-var
// identifier under the org prefix. Inputs already prefixed pass
// through; everything else is uppercased and hyphens become
// underscores.
func formatParamName(cfg *app.Config, n string) string {
	if strings.HasPrefix(n, cfg.OrgPrefix) {
		return n
	}
	normalized := strings.ReplaceAll(strings.ToUpper(n), "-", "_")
	return cfg.OrgPrefix + "_" + normalized
}

// getFullPath resolves a put-target name. Leading "/" is a literal
// path; otherwise normalize the name and route through GetAppParamPath.
func getFullPath(cfg *app.Config, n string) string {
	if strings.HasPrefix(n, "/") {
		return n
	}
	return cfg.GetAppParamPath(formatParamName(cfg, n))
}

func buildEspParamInputFromCmd(cfg *app.Config, cmd *cobra.Command) common.EspParamInput {
	name, _ := cmd.Flags().GetString("name")
	secure, _ := cmd.Flags().GetBool("secure")
	value, _ := cmd.Flags().GetString("value")
	return common.EspParamInput{
		Name:   getFullPath(cfg, name),
		Secure: secure,
		Value:  value,
	}
}

func newPutCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "put",
		Aliases: []string{"add", "create"},
		Short:   "Creates an SSM parameter with the given value.",
		Long:    `Simple command to add values to SSM parameter store.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runPut(cmd, a.Client, a.Config)
		},
	}
	cmd.Flags().StringP("name", "n", "", "The name for your parameter.")
	cmd.Flags().StringP("value", "v", "", "The value to be stored in the SSM.")
	cmd.Flags().BoolP("secure", "s", false, "Sets the SSM parameter type to 'SecureString'.")
	if err := cobra.MarkFlagRequired(cmd.Flags(), "name"); err != nil {
		fmt.Fprintln(os.Stderr, "can't set flag --name as required")
	}
	if err := cobra.MarkFlagRequired(cmd.Flags(), "value"); err != nil {
		fmt.Fprintln(os.Stderr, "can't set flag --value as required")
	}
	return cmd
}

func runPut(cmd *cobra.Command, c *client.EspClient, cfg *app.Config) error {
	param := buildEspParamInputFromCmd(cfg, cmd)
	if _, err := c.Save(param); err != nil {
		return err
	}
	saved, err := c.GetParam(common.GetOneInput{Name: param.Name})
	if err != nil {
		return err
	}
	detailDisplay(saved)
	return nil
}
