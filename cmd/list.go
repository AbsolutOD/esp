package cmd

import (
	"fmt"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func displayParams(ps []common.EspParam) {
	for _, p := range ps {
		name := aurora.BrightYellow(p.Name)
		fmt.Printf("%s: %s\n", name, p.Value)
	}
}

// getPath returns the SSM path to list. No args means "the project's
// base path"; one arg means "this exact path" (literal or short).
func getPath(cfg *app.Config, args []string) string {
	if len(args) == 0 {
		return cfg.GetAppPath()
	}
	return args[0]
}

func newListCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [path]",
		Aliases: []string{"ls"},
		Short:   "Recursively list a SSM path if given.",
		Long: `The list command gives you an easy way to recursively get all SSM parameters with a base path.
If you have a .espFile.yaml in the current directory this command will list all params under the project path.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runList(cmd, args, a.Client, a.Config)
		},
	}
	cmd.Flags().BoolP("decrypt", "d", false, "Decrypt SSM secure strings.")
	cmd.Flags().BoolP("path", "p", false, "Path to list parameters.")
	return cmd
}

func runList(cmd *cobra.Command, args []string, c *client.EspClient, cfg *app.Config) error {
	decrypt, _ := cmd.Flags().GetBool("decrypt")
	params, err := c.ListParams(common.ListParamInput{
		Path:      getPath(cfg, args),
		Decrypt:   decrypt,
		Recursive: true,
	})
	if err != nil {
		return err
	}
	displayParams(params)
	return nil
}
