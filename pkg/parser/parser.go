package parser

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Parser struct {
	Resources  []Resource
	IndentSize int
}

const (
	defaultIndentSize = 2
)

type ParserOpt func(p *Parser)

func New(opts ...ParserOpt) *Parser {
	p := &Parser{
		Resources:  make([]Resource, 0),
		IndentSize: defaultIndentSize,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func WithIndentSize(size int) ParserOpt {
	return func(p *Parser) {
		p.IndentSize = size
	}
}

// ReadFile reads a manifest and appends the resource to the parser
func (p *Parser) ReadFile(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(b)
	p.Read(buf)
	return nil
}

// Read appends new resources to a parser
func (p *Parser) Read(reader io.Reader) {
	d := yaml.NewDecoder(reader)
	rs := make([]Resource, 0)

	var r Resource
	for d.Decode(&r) == nil {
		rs = append(rs, r)
		r = make(Resource)
	}

	rs = append(rs, p.Resources...)
	p.Resources = rs
}

// Sort returns a map[string][]Resource where the key is the resource kind
func (p *Parser) Sort() map[string][]Resource {
	m := make(map[string][]Resource)

	for _, r := range p.Resources {
		kind, err := r.Kind()
		if err != nil {
			continue
		}

		m[kind] = append(m[kind], r)
	}

	return m
}

// Kustomization produces a kustomize resource
func (p *Parser) Kustomization() Resource {
	resources := make([]string, 0)
	for k := range p.Sort() {
		if strings.EqualFold(k, "kustomization") {
			continue
		}

		resources = append(resources, fmt.Sprintf("%s.yaml", strings.ToLower(k)))
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i] < resources[j]
	})

	return Resource{
		"kind":       "Kustomization",
		"apiVersion": "kustomize.config.k8s.io/v1beta1",
		"resources":  resources,
	}
}

// Create generates the folder path for a given file and writes the manifest(s) to supplied path
func (p *Parser) Create(path string, resources ...Resource) error {
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	return p.Write(f, resources...)
}

// Write uses the supplied io.Writer to write resources to a yaml file
func (p *Parser) Write(writer io.Writer, resources ...Resource) error {

	e := yaml.NewEncoder(writer)
	e.SetIndent(p.IndentSize)
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
