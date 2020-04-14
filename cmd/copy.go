package cmd

import (
	"fmt"
	"github.com/pinpt/esp/internal/client"
	"github.com/pinpt/esp/internal/common"
	"github.com/spf13/cobra"
	"os"
)

// cpCmd represents the cp command
var copyCmd = &cobra.Command{
	Use:   "cp [OPTIONS] SRC_SSM_PATH DEST_SSM_PATH",
	Aliases: []string{"copy"},
	Short: "Copy a SSM Param from its current path to a new SSM Path",
	Long: "Copy SSM value from an existing path to a new path.\n",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "" {
			fmt.Errorf("source can not be empty")
			os.Exit(1)
		}
		if args[1] == "" {
			fmt.Errorf("destination can not be empty")
			os.Exit(1)
		}
		cc := common.CopyCommand{
			Source:     args[0],
			Destination: args[1],
		}
		ec := client.New(client.EspClient{ Backend: "ssm" })
		ec.Copy(cc)
	},
	Example: "esp cp /ssm/path/key /ssm/new/path/key",
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// cpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
