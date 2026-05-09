package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/logrusorgru/aurora/v4"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func display(p common.EspParam, detail bool) {
	if detail {
		detailDisplay(p)
	} else {
		displayParam(p)
	}
}

func displayParam(p common.EspParam) {
	name := aurora.BrightYellow(p.Name)
	fmt.Printf("%s: %s\n", name, p.Value)
}

func detailDisplay(p common.EspParam) {
	data := [][]string{
		{aurora.BrightYellow("ID").String(), p.Id},
		{aurora.BrightYellow("Last_Modified").String(), p.LastModifiedDate.String()},
		{aurora.BrightYellow("Name").String(), p.Name},
		{aurora.BrightYellow("Type").String(), p.Type},
		{aurora.BrightYellow("Value").String(), p.Value},
		{aurora.BrightYellow("Version").String(), strconv.FormatInt(p.Version, 10)},
	}
	table := tablewriter.NewTable(os.Stdout)
	table.Header("Keys", "Value")
	if err := table.Bulk(data); err != nil {
		fmt.Fprintf(os.Stderr, "table render error: %v\n", err)
		return
	}
	if err := table.Render(); err != nil {
		fmt.Fprintf(os.Stderr, "table render error: %v\n", err)
	}
}

// getParamPath resolves an argument to a full SSM path. Leading "/"
// means the caller passed a literal path; everything else is routed
// through the project-aware GetAppParamPath.
func getParamPath(cfg *app.Config, p string) string {
	if strings.HasPrefix(p, "/") {
		return p
	}
	return cfg.GetAppParamPath(p)
}

func newGetCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [path]",
		Short: "Query path for SSM",
		Long:  `Allows you to get a specific ssm parameter with an exact path or recursively get params.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runGet(cmd, args, a.Client, a.Config)
		},
	}
	cmd.Flags().BoolP("details", "t", false, "Show all of the attributes of a parameter.")
	cmd.Flags().BoolP("decrypt", "d", false, "Decrypt SSM secure strings.")
	return cmd
}

func runGet(cmd *cobra.Command, args []string, c *client.EspClient, cfg *app.Config) error {
	decrypt, _ := cmd.Flags().GetBool("decrypt")
	details, _ := cmd.Flags().GetBool("details")
	param, err := c.GetParam(common.GetOneInput{
		Name:    getParamPath(cfg, args[0]),
		Decrypt: decrypt,
	})
	if err != nil {
		return err
	}
	display(param, details)
	return nil
}
