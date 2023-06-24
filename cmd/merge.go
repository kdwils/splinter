package cmd

import (
	"log"
	"os"

	"github.com/kdwils/splinter/internal/splinter"
	"github.com/kdwils/splinter/pkg/parser"
	"github.com/spf13/cobra"
)

var (
	delete bool
	stdOut bool
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge split manifests back together",
	Long:  `merge split manifests back together`,
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
			Delete:     delete,
			StdOut:     stdOut,
		}

		err := splinter.Merge(p, in)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	mergeCmd.Flags().StringSliceVarP(&input, "input", "i", input, "provide /path/to/input/ or input.yaml")
	mergeCmd.Flags().StringSliceVarP(&exclusions, "exclusions", "e", exclusions, "files or directories to exclude")
	mergeCmd.Flags().BoolVarP(&kustomize, "kustomize", "k", false, "spit out a kustomization.yaml")
	mergeCmd.Flags().StringVarP(&output, "output", "o", output, "provide /path/to/output/dir")
	mergeCmd.Flags().BoolVarP(&delete, "delete", "d", false, "delete files that have been merged")
	mergeCmd.Flags().BoolVar(&stdOut, "std-out", false, "write to stdout instead of disk")
}
