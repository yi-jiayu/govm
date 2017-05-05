package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yi-jiayu/govm/lib"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Switch to a different Go version",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// get destination directory for govm install
		searchDir, err := cmd.Flags().GetString("govm_home")
		if err != nil {
			return err
		}

		cv, err := lib.CurrentGoVersion(searchDir)
		if err != nil {
			return err
		}

		if len(args) != 1 {
			return cmd.Usage()
		}

		version := args[0]

		// validate version
		if !lib.ValidateSemver(version) {
			fmt.Printf("%s is not a valid Go version string.\n", version)
			return nil
		}

		gv, err := lib.GoVersionOutput()
		if err != nil {
			return err
		}

		fmt.Printf("Now using: %s", string(gv))
		fmt.Printf("You are trying to switch to Go version: %s\n", version)

		vs, err := lib.InstalledGoVersions(searchDir)
		if err != nil {
			return err
		}

		found := false
		for _, v := range vs {
			if v == version {
				found = true
				break
			}
		}

		if found {
			if cv == version {
				fmt.Printf("Go version %s is already the currently active version.\n", cv)

				return nil
			} else {
				fmt.Printf("Changing to Go version %s...\n", version)

				err := lib.SwitchGoVersion(version, searchDir)
				if err != nil {
					return err
				}

				gv, err := lib.GoVersionOutput()
				if err != nil {
					return err
				}

				fmt.Printf("Now using: %s", string(gv))

				return nil
			}
		} else {
			fmt.Printf("Go version %s is not currently installed.\n", version)

			return nil
		}
	},
}

func init() {
	RootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
