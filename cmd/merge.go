package cmd

import (
	"log"
	"os"

	"github.com/kdwils/splinter/parser"
	"github.com/spf13/cobra"
)

var (
	mergeInputFiles       []string
	mergeOutputPath       string
	mergeIncludeKustomize bool
	mergeExclusions       []string
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge split manifests back together",
	Long:  `merge split manifests back together`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := parser.New()

		var stdin *os.File
		// shoutout https://stackoverflow.com/questions/22744443/check-if-there-is-something-to-read-on-stdin-in-golang
		if s, err := os.Stdin.Stat(); err == nil && (s.Mode()&os.ModeCharDevice) == 0 {
			stdin = os.Stdin
		}

		err := p.Merge(mergeInputFiles, stdin, mergeOutputPath)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().StringSliceVarP(&mergeInputFiles, "input", "i", mergeInputFiles, "provide /path/to/input/ or input.yaml")
	mergeCmd.Flags().StringSliceVarP(&mergeExclusions, "exclusions", "e", mergeExclusions, "files or directories to exclude")
	mergeCmd.Flags().BoolVarP(&mergeIncludeKustomize, "kustomize", "k", false, "spit out a kustomization.yaml")
	mergeCmd.Flags().StringVarP(&mergeOutputPath, "output", "o", mergeOutputPath, "provide /path/to/output/file.yaml")
}
