package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yi-jiayu/govm/lib"
	"path/filepath"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "govm",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("version")

		if v {
			PrintVersion()
		} else {
			cmd.Help()
		}
	},

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		v, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}

		if v {
			lib.SetVerbose(true)
		}

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.govm/.govm.yaml)")
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Print additional output")
	RootCmd.PersistentFlags().String("govm_home", "C:\\", "Set storage location for Go installations")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().Bool("version", false, "Show version information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".govm")                                          // name of config file (without extension)
	viper.AddConfigPath(".")                                              // search in current working directory for config file
	viper.AddConfigPath(filepath.Join(os.Getenv("USERPROFILE"), ".govm")) // adding USERPROFILE/.govm as second search path

	viper.SetEnvPrefix("govm")
	viper.AutomaticEnv() // read in environment variables that match

	viper.BindPFlag("govm_home", RootCmd.Flags().Lookup("govm_home"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
