package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	output     string
	input      []string
	kustomize  bool
	exclusions []string
	merge      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "splinter",
	Short:        "cli to manipulate kubernetes manifest files",
	Long:         `cli to manipulate kubernetes manifest files`,
	SilenceUsage: false,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.splinter.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
