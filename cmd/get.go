package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/logrusorgru/aurora"
	"github.com/olekukonko/tablewriter"
	"github.com/pinpt/esp/internal/client"
	"github.com/pinpt/esp/internal/errors"
	"github.com/spf13/cobra"
)

// getParam Queries the ssm param
func getParam(ec client.EspConfig, d bool, key string) *ssm.Parameter {
	si := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(d),
	}
	resp, err := ec.Svc.GetParameter(si)
	if err != nil {
		errors.CheckSSMGetParameters(err)
	}

	return resp.Parameter
}

func display(p *ssm.Parameter, detail bool) {
	if detail {
		detailDisplay(p)
	} else {
		displayParam(p)
	}
}

func displayParam(p *ssm.Parameter) {
	name := aurora.BrightYellow(*p.Name)
	fmt.Printf("%s: %s\n", name, *p.Value)
}

func detailDisplay(p *ssm.Parameter) {
	data := [][]string{
		[]string{aurora.BrightYellow("ARN").String(), *p.ARN},
		[]string{aurora.BrightYellow("Last_Modified").String(), p.LastModifiedDate.String()},
		[]string{aurora.BrightYellow("Name").String(), *p.Name},
		[]string{aurora.BrightYellow("Type").String(), *p.Type},
		[]string{aurora.BrightYellow("Value").String(), *p.Value},
		[]string{aurora.BrightYellow("Version").String(), string(*p.Version)},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Keys", "Value"})
	table.AppendBulk(data)
	table.Render()
}

// getCmd represents the path command
var getCmd = &cobra.Command{
	Use:   "get [path]",
	Short: "Query path for SSM",
	Long:  `Allows you to get a specific ssm parameter with an exact path or recursively get params.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		ec := client.New(region)
		decrypt, _ := cmd.Flags().GetBool("decrypt")
		details, _ := cmd.Flags().GetBool("details")

		param := getParam(ec, decrypt, args[0])
		display(param, details)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolP("details", "t", false, "Show all of the attributes of a parameter.")
	getCmd.Flags().BoolP("decrypt", "d", false, "Decrypt SSM secure strings.")
}
