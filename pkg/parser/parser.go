package parser

import (
	"io"

	"gopkg.in/yaml.v3"
)

type Parser struct {}

const (
	indentSize = 2
)

func New() *Parser {
	return new(Parser)
}

// Sort organizes resources by kind from reader
func (p *Parser) Sort(reader io.Reader) map[string][]Resource {
	d := yaml.NewDecoder(reader)
	m := make(map[string][]Resource)

	var r Resource
	for d.Decode(&r) == nil {
		kind, err := r.Kind()
		if err != nil {
			continue
		}

		m[kind] = append(m[kind], r)
		r = make(Resource)
	}

	return m
}

// Flatten slice of resources from reader
func (p *Parser) Flatten(reader io.Reader) []Resource {
	d := yaml.NewDecoder(reader)
	rs := make([]Resource, 0)

	var r Resource
	for d.Decode(&r) == nil {
		rs = append(rs, r)
		r = make(Resource)
	}

	return rs
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
