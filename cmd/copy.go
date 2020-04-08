package cmd

//var copyDescription = "Copy SSM value from an existing path to a new path.\n"

/*type CopyCommand struct {
	source      string
	destination string
}

func (cc *CopyCommand) runCopy(region string) {
	ec := client.New(region)
	
}*/

// cpCmd represents the cp command
/*var copyCmd = &cobra.Command{
	Use:   "cp [OPTIONS] SRC_SSM_PATH DEST_SSM_PATH",
	Aliases: []string{"cp"},
	Short: "Copy a SSM Param from its current path to a new SSM Path",
	Long: "Copy SSM value from an existing path to a new path.\n",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		fmt.Println("cp called")
		if args[0] == "" {
			return fmt.Errorf("Source can not be empty")
		}
		if args[1] == "" {
			return fmt.Errorf("Destination can not be empty")
		}
		cc := &CopyCommand{
			source:     args[0],
			destination: args[1],
		}
		return cc.runCopy(region)
	},
	Example: "esp cp /ssm/path/key /ssm/new/path/key",
}*/

func init() {
	//rootCmd.AddCommand(copyCmd)

	// cpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
