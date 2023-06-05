package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

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
	merge      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "splinter",
	Short:        "split manifests into multiple files by resource kind",
	Long:         `split manifests into multiple files by resource kind`,
	SilenceUsage: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := parser.New()

		// shoutout https://stackoverflow.com/questions/22744443/check-if-there-is-something-to-read-on-stdin-in-golang
		if s, err := os.Stdin.Stat(); err == nil && (s.Mode()&os.ModeCharDevice) == 0 {
			p.Read(os.Stdin)
		}

		files := make([]string, 0)
		files = append(files, filesFromInput(args)...)
		files = append(files, filesFromInput(input)...)
		files = removeExclusions(files, exclusions)

		for _, f := range files {
			p.ReadFile(f)
		}

		if kustomize {
			kustomizeFilePath := path.Join(getKustomizePath(output), "kustomization.yaml")
			err := p.Create(kustomizeFilePath, p.Kustomization())
			if err != nil {
				return err
			}
		}

		if merge {
			return p.Create(output, p.Resources...)
		}

		for k, v := range p.Sort() {
			filepath := path.Join(output, fmt.Sprintf("%s.yaml", strings.ToLower(k)))
			if err := p.Create(filepath, v...); err != nil {
				return err
			}
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

	rootCmd.Flags().StringVarP(&output, "output", "o", pwd, "provide /path/to/output/dir")
	rootCmd.Flags().StringSliceVarP(&input, "input", "i", input, "provide /path/to/input/ or input.yaml")
	rootCmd.Flags().StringSliceVarP(&exclusions, "exclusions", "e", exclusions, "files or directories to exclude")
	rootCmd.Flags().BoolVar(&merge, "merge", false, "merge splintered manifests back together")
	rootCmd.Flags().BoolVarP(&kustomize, "kustomize", "k", false, "spit out a kustomization.yaml")
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

func filesFromInput(input []string) []string {
	files := make([]string, 0)
	for _, p := range input {
		if strings.EqualFold(filepath.Ext(p), ".yaml") {
			files = append(files, p)
			continue
		}

		fileInfo, err := os.Stat(p)
		if err != nil {
			continue
		}

		if !fileInfo.IsDir() {
			continue
		}

		dir, err := os.ReadDir(p)
		if err != nil {
			continue
		}

		for _, file := range dir {
			files = append(files, path.Join(p, file.Name()))
		}
	}

	return files
}

func removeExclusions(input, exclusions []string) []string {
	for i, in := range input {
		for _, e := range exclusions {
			if strings.EqualFold(in, e) {
				input = removeIndex(input, i)
			}
		}
	}

	return input
}

func removeIndex(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func getKustomizePath(output string) string {
	if strings.HasSuffix(output, ".yaml") {
		dirs := strings.Split(output, "/")
		if len(dirs) == 0 {
			return "."
		}

		return strings.Join(dirs[:len(dirs)-1], "/")
	}

	return output
}
