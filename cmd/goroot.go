package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yi-jiayu/govm/lib"
)

// gorootCmd represents the goroot command
var gorootCmd = &cobra.Command{
	Use:   "root [goroot]",
	Short: "Print the current GOROOT",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			// fixme: this does not show the actual value of goroot but the normalised value
			fmt.Printf("GOROOT is currently set to: %s\n", lib.Goroot())
		case 1:
			fmt.Printf("You are trying to set GOROOT to: %s\n", args[0])
			fmt.Println("not implemented")
		default:
			cmd.Usage()
		}
	},
}

func init() {
	RootCmd.AddCommand(gorootCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gorootCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gorootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
