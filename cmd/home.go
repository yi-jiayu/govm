package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// homeCmd represents the home command
var homeCmd = &cobra.Command{
	Use:   "home [path]",
	Short: "Get or set the current govm home",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			home, err := cmd.Flags().GetString("govm_home")
			if err != nil {
				fmt.Printf("An error occured: err %v\n", err)
			}

			fmt.Printf("govm home: %s\n", home)
		case 1:
			fmt.Printf("Setting govm home to %s...\n", args[0])
			fmt.Println("not implemented")
		default:
			cmd.Usage()
		}
	},
}

func init() {
	RootCmd.AddCommand(homeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// homeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// homeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
