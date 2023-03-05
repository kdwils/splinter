package parser

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Parser struct {
	Resources []Resource
}

const (
	indentSize = 2
)

func New() *Parser {
	return &Parser{
		Resources: make([]Resource, 0),
	}
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

func (p *Parser) kustomizeResourcePaths() []string {
	keys := make([]string, 0)
	for k := range p.Sort() {
		if strings.EqualFold(k, "kustomization") {
			continue
		}

		keys = append(keys, p.YamlFile(k))
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}

// YamlFileName returns a lowercase name.yaml
func (p *Parser) YamlFile(name string) string {
	return fmt.Sprintf("%s.yaml", strings.ToLower(name))
}

// Kustomize produces a kustomize resource
func (p *Parser) Kustomize() Resource {
	return Resource{
		"kind":       "Kustomization",
		"apiVersion": "kustomize.config.k8s.io/v1beta1",
		"resources":  p.kustomizeResourcePaths(),
	}
}

func (p *Parser) Write(writer io.Writer, resources ...Resource) error {
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
