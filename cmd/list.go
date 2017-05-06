package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yi-jiayu/govm/lib"
	"log"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List installed Go versions",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get destination directory for govm install
		installDir := viper.GetString("install-dir")

		vs, err := lib.InstalledGoVersions(installDir)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		cv, err := lib.CurrentGoVersion(installDir)
		if err != nil {
			cv = ""
		}

		if len(vs) > 0 {
			if cv == "" {
				fmt.Printf("govm did found the following Go installations in the current govm install directory (%s).\n", installDir)
				fmt.Println("You can run \"govm use [version]\" to use one of them.")
			}

			for _, v := range vs {
				if v == cv {
					fmt.Printf("  * %s\n", v)
				} else {
					fmt.Printf("    %s\n", v)
				}
			}
		} else {
			fmt.Printf("govm did not find any Go installations in the current govm install directory (%s).\n", installDir)
			fmt.Println("You can run \"govm install [version]\" and \"govm use [version]\" to install and use a new Go version.")
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
