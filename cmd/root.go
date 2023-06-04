package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/kdwils/splinter/pkg/parser"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	output     string
	input      []string
	kustomize  bool
	exclusions []string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "splinter",
	Short:        "split manifests into multiple files by resource kind",
	Long:         `split manifests into multiple files by resource kind`,
	SilenceUsage: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := parser.New()

		// https://stackoverflow.com/questions/22744443/check-if-there-is-something-to-read-on-stdin-in-golang
		if s, err := os.Stdin.Stat(); err == nil && (s.Mode()&os.ModeCharDevice) == 0 {
			p.Read(os.Stdin)
		}

		files := make([]string, 0)
		files = append(files, filesFromInput(args)...)
		files = append(files, filesFromInput(input)...)
		files = removeExclusions(files, exclusions)

		for _, f := range files {
			buf, err := readFile(f)
			if err != nil {
				return err
			}

			p.Read(buf)
		}

		for k, v := range p.Sort() {
			f, err := createFile(path.Join(output, parser.YamlFileName(k)))
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

	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("could not determine current working dir: %v", err)
	}

	rootCmd.Flags().StringVarP(&output, "output", "o", pwd, "provide /path/to/output/dir, defaults to current working dir")
	rootCmd.Flags().StringSliceVarP(&input, "input", "i", input, "provide /path/to/input/ or input.yaml")
	rootCmd.Flags().StringSliceVarP(&exclusions, "exclusions", "e", exclusions, "files or directories to exclusions")
	rootCmd.Flags().BoolVarP(&kustomize, "kustomize", "k", false, "output a simple kustomization.yaml as well")

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
