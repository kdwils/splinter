package cmd

import (
	"log"
	"os"

	"github.com/kdwils/splinter/internal/splinter"
	"github.com/kdwils/splinter/pkg/parser"
	"github.com/spf13/cobra"
)

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "split a single kubernetes manifest into many",
	Long:  `split a single kubernetes manifest into many`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := parser.New()

		// shoutout https://stackoverflow.com/questions/22744443/check-if-there-is-something-to-read-on-stdin-in-golang
		if s, err := os.Stdin.Stat(); err == nil && (s.Mode()&os.ModeCharDevice) == 0 {
			p.Read(os.Stdin)
		}

		in := &splinter.Input{
			InputFiles: append(input, args...),
			OutputPath: output,
			Kustomize:  kustomize,
			Exclusions: exclusions,
		}

		err := splinter.Split(p, in)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	splitCmd.Flags().StringSliceVarP(&input, "input", "i", input, "provide /path/to/input/ or input.yaml")
	splitCmd.Flags().StringSliceVarP(&exclusions, "exclusions", "e", exclusions, "files or directories to exclude")
	splitCmd.Flags().BoolVarP(&kustomize, "kustomize", "k", false, "spit out a kustomization.yaml")
	splitCmd.Flags().StringVarP(&output, "output", "o", output, "provide /path/to/output/dir")
}
