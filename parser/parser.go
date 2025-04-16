package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/kdwils/splinter/pkg/fio"
)

type Parser struct {
	indentSize int
	fio        fio.FileIO
}

const (
	defaultIndentSize = 2
)

type ParserOpt func(p *Parser)

func New(opts ...ParserOpt) *Parser {
	p := &Parser{
		indentSize: defaultIndentSize,
		fio:        fio.NewDefaultFileIO(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func WithIndentSize(size int) ParserOpt {
	return func(p *Parser) {
		p.indentSize = size
	}
}

func WithFileIO(fio fio.FileIO) ParserOpt {
	return func(p *Parser) {
		p.fio = fio
	}
}

func (p *Parser) Merge(files []string, stdin io.Reader, outputPath string) error {
	resources := make([]Resource, 0)

	if stdin != nil {
		resources = append(resources, readResource(stdin)...)
	}

	for _, f := range p.filesFromInput(files) {
		buf, err := p.readFileToBuffer(f)
		if err != nil {
			return err
		}
		for _, r := range readResource(buf) {
			_, err := r.Kind()
			if err != nil {
				continue
			}

			resources = append(resources, r)
		}
	}

	if outputPath != "" {
		return p.write(outputPath, p.indentSize, resources...)
	}

	return write(os.Stdout, p.indentSize, resources...)
}

func (p *Parser) Split(inputFiles []string, stdin io.Reader, outputPath string, kustomize bool) error {
	resources := make([]Resource, 0)

	if stdin != nil {
		resources = append(resources, readResource(stdin)...)
	}

	for _, f := range p.filesFromInput(inputFiles) {
		buf, err := p.readFileToBuffer(f)
		if err != nil {
			return err
		}
		for _, r := range readResource(buf) {
			kind, err := r.Kind()
			if err != nil {
				continue
			}
			if strings.EqualFold(kind, "kustomization") {
				continue
			}

			resources = append(resources, r)
		}
	}

	resourceMap := resourcesToMap(resources)

	if kustomize {
		resources := make([]string, 0)
		for k := range resourceMap {
			resources = append(resources, fmt.Sprintf("%s.yaml", strings.ToLower(k)))
		}
		resourceMap["kustomization"] = append(resourceMap["kustomization"], newKustomizeResource(resources...))
	}

	for k, v := range resourceMap {
		filepath := path.Join(outputPath, fmt.Sprintf("%s.yaml", strings.ToLower(k)))
		err := p.write(filepath, p.indentSize, v...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) readFileToBuffer(file string) (*bytes.Buffer, error) {
	b, err := p.fio.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

func (p *Parser) filesFromInput(input []string) []string {
	files := make([]string, 0)
	for _, f := range input {
		if strings.EqualFold(filepath.Ext(f), ".yaml") {
			files = append(files, f)
			continue
		}

		fileInfo, err := p.fio.Stat(f)
		if err != nil {
			continue
		}

		if !fileInfo.IsDir() {
			continue
		}

		dir, err := p.fio.ReadDir(f)
		if err != nil {
			continue
		}

		for _, file := range dir {
			files = append(files, path.Join(f, file.Name()))
		}
	}

	return files
}

func (p *Parser) write(path string, indentSize int, resources ...Resource) error {
	if _, err := p.fio.Stat(filepath.Dir(path)); errors.Is(err, os.ErrNotExist) {
		err := p.fio.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := p.fio.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return write(f, indentSize, resources...)
}

func write(writer io.Writer, indentSize int, resources ...Resource) error {
	e := yaml.NewEncoder(writer)
	e.SetIndent(indentSize)
	defer e.Close()

	var err error
	for _, r := range resources {
		err = e.Encode(r)
		if err != nil {
			return err
		}
	}

	return nil
}
