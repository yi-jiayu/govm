package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yi-jiayu/govm/lib"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [version]",
	Short: "Uninstall an installed Go version",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		// get destination directory for govm install
		installDir := viper.GetString("install-dir")

		version := args[0]

		// validate version
		if !lib.ValidateSemver(version) {
			fmt.Printf("%s is not a valid Go version string.\n", version)
			return nil
		}

		// check if version is actually installed
		ivs, err := lib.InstalledGoVersions(installDir)
		if err != nil {
			return err
		}

		installed := false
		for _, iv := range ivs {
			if version == iv {
				installed = true
				break
			}
		}

		if !installed {
			fmt.Printf("Go version %s is not currently installed.\n", version)
			return nil
		}

		// check if version is the currently active Go installation
		cv, err := lib.CurrentGoVersion(installDir)
		if err != nil {
			return err
		}

		if version == cv {
			fmt.Printf("You are currently using Go %s. Change to a different Go version before you can uninstall it.\n", cv)
			return nil
		}

		// uninstall
		fmt.Printf("Uninstalling Go %s...\n", version)
		err = lib.UninstallGoVersion(version, installDir)
		if err != nil {
			return err
		}
		fmt.Println("Done!")

		fmt.Printf("Go version %s successfully uninstalled.\n", version)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(uninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
