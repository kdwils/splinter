package cmd

import (
	"os"
	"path/filepath"

	fh "github.com/kdwils/splinter/internal/file"
	"github.com/kdwils/splinter/pkg/parser"
	"github.com/spf13/cobra"
)

// flattenCmd represents the flatten command
var flattenCmd = &cobra.Command{
	Use:          "flatten",
	Short:        "flatten multiple kubernetes yaml resources into one file",
	Long:         `flatten multiple kubernetes yaml resources into one file`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fs, err := fh.ListFilesFromInput(input)
		if err != nil {
			return err
		}

		p := parser.New()
		rs := make([]parser.Resource, 0)
		for _, f := range fs {
			buf, err := fh.FileToBuffer(f)
			if err != nil {
				return err
			}

			rs = append(rs, p.Flatten(buf)...)
		}

		err = os.MkdirAll(filepath.Dir(output), os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.Create(output)
		if err != nil {
			return err
		}

		return p.Write(f, rs...)
	},
}

func init() {
	rootCmd.AddCommand(flattenCmd)

	flattenCmd.Flags().StringSliceVarP(&input, "input", "i", input, "/path/to/input.yaml or /path/to/dir, or both")
	flattenCmd.Flags().StringVarP(&output, "output", "o", "manifest.yaml", "/path/to/output.yaml")
}
