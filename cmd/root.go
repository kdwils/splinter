package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/kdwils/splinter/pkg/parser"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	input     []string
	output    string
	kustomize bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "splinter",
	Short:        "split manifests into multiple files by resource kind",
	Long:         `split manifests into multiple files by resource kind`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := parser.New()

		// https://stackoverflow.com/questions/22744443/check-if-there-is-something-to-read-on-stdin-in-golang
		if s, err := os.Stdin.Stat(); err == nil && (s.Mode()&os.ModeCharDevice) == 0 {
			p.Read(os.Stdin)
		}

		for _, f := range filesFromInput(input) {
			buf, err := readFile(f)
			if err != nil {
				return err
			}

			p.Read(buf)
		}

		for k, v := range p.Sort() {
			f, err := createFile(path.Join(output, p.YamlFile(k)))
			if err != nil {
				return err
			}

			err = p.Write(f, v...)
			if err != nil {
				return err
			}
		}

		if kustomize {
			f, err := createFile(path.Join(output, "kustomization.yaml"))
			if err != nil {
				return err
			}

			return p.Write(f, p.Kustomize())
		}

		return nil
	},
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.splinter.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringSliceVarP(&input, "input", "i", input, "/path/to/input.yaml or /path/to/dir, or both")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "/path/to/output/dir")
	rootCmd.Flags().BoolVarP(&kustomize, "kustomize", "k", false, "create a simple kustomization.yaml")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".splinter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".splinter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
