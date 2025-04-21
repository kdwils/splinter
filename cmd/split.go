package cmd

import (
	"log"
	"os"

	"github.com/kdwils/splinter/parser"
	"github.com/spf13/cobra"
)

var (
	splitInputFiles       []string
	splitOutputPath       string
	splitIncludeKustomize bool
	splitExclusions       []string
	splitCreateKustomize  bool
)

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "split a single kubernetes manifest into many",
	Long:  `split a single kubernetes manifest into many`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := parser.New()

		var stdin *os.File
		// shoutout https://stackoverflow.com/questions/22744443/check-if-there-is-something-to-read-on-stdin-in-golang
		if s, err := os.Stdin.Stat(); err == nil && (s.Mode()&os.ModeCharDevice) == 0 {
			stdin = os.Stdin
		}

		for _, a := range args {
			splitInputFiles = append(splitInputFiles, a)
		}

		err := p.Split(splitInputFiles, stdin, splitOutputPath, splitCreateKustomize)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	splitCmd.Flags().StringSliceVarP(&splitInputFiles, "input", "i", splitInputFiles, "provide /path/to/input/ or input.yaml")
	splitCmd.Flags().StringSliceVarP(&splitExclusions, "exclusions", "e", splitExclusions, "files or directories to exclude")
	splitCmd.Flags().BoolVarP(&splitCreateKustomize, "kustomize", "k", splitCreateKustomize, "spit out a kustomization.yaml")
	splitCmd.Flags().StringVarP(&splitOutputPath, "output", "o", splitOutputPath, "provide /path/to/output/dir")
	splitCmd.MarkFlagRequired("output")
}
