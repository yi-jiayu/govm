package cmd

import (
	"fmt"

	"bufio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yi-jiayu/govm/lib"
	"log"
	"os"
	"os/exec"
	"strings"
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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			return
		}

		// get destination directory for govm install
		installDir := viper.GetString("install-dir")

		ivs, err := lib.InstalledGoVersions(installDir)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		if len(ivs) == 0 {
			fmt.Printf("govm did not find any Go installations in the current govm install directory (%s).\n", installDir)
			fmt.Println("You can run \"govm install [version]\" and \"govm use [version]\" to install and use a new Go version.")
			return
		}

		version := args[0]

		// validate version
		if !lib.ValidateSemver(version) {
			fmt.Printf("%s is not a valid Go version string.\n", version)
			return
		}

		var alreadyExists bool
		cv, err := lib.CurrentGoVersion(installDir)
		if err != nil && err != lib.ErrNotManaged {
			if err == lib.ErrNotASymlink {
				alreadyExists = true
			} else {
				log.Fatalf("error: %v", err)
			}
		}

		gv, err := lib.GoVersionOutput()
		if err != nil {
			gv = "not applicable"
		}
		gv = strings.TrimPrefix(gv, "go version ")
		gv = strings.TrimSpace(gv)
		fmt.Printf("Current Go version: %s\n", gv)

		found := false
		for _, v := range ivs {
			if v == version {
				found = true
				break
			}
		}

		if found {
			if cv == version {
				fmt.Printf("Go version %s is already the currently active version.\n", cv)
				return
			} else {
				if alreadyExists {
					fmt.Printf("GOROOT (%s) already exists but is not a symlink to a Go installation in the current govm install dir (%s).\n", lib.Goroot(), installDir)
					fmt.Print("Do you want govm to remove and relink GOROOT? (y/N) ")

					r := bufio.NewReader(os.Stdin)
					input, err := r.ReadString('\n')
					if err != nil {
						log.Fatalf("error: %v", err)
					}

					input = strings.TrimSpace(input)

					if input == "y" || input == "Y" || input == "yes" {
						// todo: prompt for confirmation if goroot is a non-empty directory
						fmt.Printf("Removing %s...\n", lib.Goroot())
						err := os.RemoveAll(lib.Goroot())
						if err != nil {
							log.Fatalf("error: %v", err)
						}
					} else {
						fmt.Println("Aborting...")
						os.Exit(0)
					}
				}

				fmt.Printf("Changing to Go version %s...\n", version)

				err := lib.SwitchGoVersion(version, installDir)
				if err != nil {
					log.Fatalf("error: %v", err)
				}

				gv, err := lib.GoVersionOutput()
				if err != nil {
					if e, ok := err.(*exec.ExitError); ok {
						log.Fatalf("error: %s", string(e.Stderr))
					} else {
						log.Fatalf("%v", err)
					}
				}
				gv = strings.TrimPrefix(gv, "go version ")
				gv = strings.TrimSpace(gv)

				fmt.Printf("Current Go version: %s", string(gv))
				return
			}
		} else {
			fmt.Printf("Go version %s is not currently installed.\n", version)
			fmt.Printf("Use \"govm install %s\" to install it first.\n", version)
			return
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
