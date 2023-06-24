package splinter

import (
	"os"
	"path"
	"path/filepath"

	"github.com/kdwils/splinter/pkg/parser"
)

func Merge(p *parser.Parser, input *Input) error {
	input.removeExclusions()
	files := filesFromInput(input.InputFiles)

	for _, f := range files {
		p.ReadFile(f)
	}

	p.Sanitize()

	if input.OutputPath == "" {
		out, err := getDefaultOutputPath()
		if err != nil {
			return err
		}

		input.OutputPath = filepath.Join(out, "manifest.yaml")
	}

	if input.Kustomize {
		kustomizeFilePath := path.Join(input.kustomizeFilePath(), "kustomization.yaml")
		resourceName := getFileNameFromPath(input.OutputPath)
		err := p.Create(kustomizeFilePath, kustomization(resourceName))
		if err != nil {
			return err
		}
	}

	err := p.Create(input.OutputPath, p.Resources...)
	if err != nil {
		return err
	}

	if input.Delete {
		for _, f := range files {
			os.Remove(f)
		}
	}

	return nil
}
