package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yi-jiayu/govm/lib"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a new Go version",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Aliases: []string{"i"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		version := args[0]

		// validate version
		if !lib.ValidateSemver(version) {
			fmt.Printf("%s is not a valid Go version string.\n", version)
			return nil
		}

		// check if version is already installed
		ivs, err := lib.InstalledGoVersions()
		if err != nil {
			return err
		}

		for _, iv := range ivs {
			if version == iv {
				fmt.Printf("Go version %s is already installed.\n", version)
				return nil
			}
		}

		// todo: clean temp files if an error occurs or deterministically cache them for reuse
		// download new go version
		fmt.Printf("Downloading Go version %s...\n", version)
		dl, err := lib.DownloadRemoteVersion(version)
		if err != nil {
			return err
		}
		fmt.Println("Done!")

		// extract new go version
		fmt.Println("Extracting...")
		temp, err := lib.ExtractDownloadedGoVersion(dl)
		if err != nil {
			return err
		}
		fmt.Println("Done!")

		fmt.Println("Installing...")
		err = lib.InstallGoVersion(version, temp)
		if err != nil {
			return err
		}
		fmt.Println("Done!")

		// clean download from temp directory
		fmt.Println("Cleaning up...")
		err = os.Remove(dl)
		if err != nil {
			return err
		}
		fmt.Println("Done!")

		fmt.Printf("Go version %s successfully installed!\n", version)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
