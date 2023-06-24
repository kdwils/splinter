package splinter

import (
	"fmt"
	"path"
	"strings"

	"github.com/kdwils/splinter/pkg/parser"
)

func Split(p *parser.Parser, input *Input) error {
	input.removeExclusions()
	files := filesFromInput(input.InputFiles)

	for _, f := range files {
		p.ReadFile(f)
	}

	p.Sanitize()

	var err error
	if input.OutputPath == "" {
		input.OutputPath, err = getDefaultOutputPath()
		if err != nil {
			return err
		}
	}

	sorted := p.Sort()
	if input.Kustomize {
		kustomizeFilePath := path.Join(input.kustomizeFilePath(), "kustomization.yaml")

		resources := make([]string, 0)
		for k := range sorted {
			if strings.EqualFold(k, "kustomization") {
				continue
			}

			resources = append(resources, fmt.Sprintf("%s.yaml", strings.ToLower(k)))
		}

		err = p.Create(kustomizeFilePath, kustomization(resources...))
		if err != nil {
			return err
		}
	}

	for k, v := range p.Sort() {
		filepath := path.Join(input.OutputPath, fmt.Sprintf("%s.yaml", strings.ToLower(k)))
		err = p.Create(filepath, v...)
		if err != nil {
			return err
		}
	}

	return nil
}
