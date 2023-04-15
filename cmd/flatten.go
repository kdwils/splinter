package cmd

import (
	"github.com/kdwils/splinter/pkg/parser"
	"github.com/spf13/cobra"
)

var (
	flattenOutput string
)

// flattenCmd represents the flatten command
var flattenCmd = &cobra.Command{
	Use:          "flatten",
	Short:        "flatten multiple kubernetes yaml resources into one file",
	Long:         `flatten multiple kubernetes yaml resources into one file`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := parser.New()
		for _, f := range filesFromInput(input) {
			buf, err := readFile(f)
			if err != nil {
				return err
			}

			p.Read(buf)
		}

		f, err := createFile(flattenOutput)
		if err != nil {
			return err
		}

		return p.Write(f, p.Resources...)
	},
}

func init() {
	rootCmd.AddCommand(flattenCmd)

	flattenCmd.Flags().StringSliceVarP(&input, "input", "i", input, "/path/to/input.yaml or /path/to/dir, or both")
	flattenCmd.Flags().StringVarP(&flattenOutput, "output", "o", "manifest.yaml", "/path/to/output.yaml")
}
