package parser

import (
	"errors"
	"io"

	"slices"

	"gopkg.in/yaml.v3"
)

var (
	ErrKindKeyNotFound = errors.New("no resource key 'kind' found")
)

// Resource is an alias for map[string]any
type Resource map[string]any

// Kind returns the kind of the resource if the key exists, otherwise returns an error
func (r Resource) Kind() (string, error) {
	k, ok := r["kind"]
	if !ok {
		return "", ErrKindKeyNotFound
	}

	return k.(string), nil
}

func newKustomizeResource(resources ...string) Resource {
	slices.Sort(resources)
	return Resource{
		"kind":       "Kustomization",
		"apiVersion": "kustomize.config.k8s.io/v1beta1",
		"resources":  resources,
	}
}

func readResource(reader io.Reader) []Resource {
	d := yaml.NewDecoder(reader)
	rs := make([]Resource, 0)

	var r Resource
	for d.Decode(&r) == nil {
		rs = append(rs, r)
		r = make(Resource)
	}

	return rs
}

func resourcesToMap(resources []Resource) map[string][]Resource {
	m := make(map[string][]Resource)
	for _, r := range resources {
		kind, err := r.Kind()
		if err != nil {
			continue
		}

		m[kind] = append(m[kind], r)
	}

	return m
}
