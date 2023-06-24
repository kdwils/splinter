package splinter

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/kdwils/splinter/pkg/parser"
)

type Input struct {
	InputFiles []string
	Exclusions []string
	Kustomize  bool
	Delete     bool
	StdOut     bool
	OutputPath string
}

func (i *Input) removeExclusions() {
	for index, in := range i.InputFiles {
		for _, e := range i.Exclusions {
			if strings.EqualFold(in, e) {
				i.InputFiles = removeIndex(i.InputFiles, index)
			}
		}
	}
}

func removeIndex(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func (i *Input) kustomizeFilePath() string {
	if !strings.HasSuffix(i.OutputPath, ".yaml") && i.OutputPath != "" {
		return i.OutputPath
	}

	dirs := strings.Split(i.OutputPath, "/")
	if len(dirs) == 1 {
		return "."
	}

	return fmt.Sprintf("%s/", strings.Join(dirs[:len(dirs)-1], "/"))
}

func getDefaultOutputPath() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return pwd, nil
}

// Kustomization produces a kustomize resource
func kustomization(resources ...string) parser.Resource {
	sort.Slice(resources, func(i, j int) bool {
		return resources[i] < resources[j]
	})

	return parser.Resource{
		"kind":       "Kustomization",
		"apiVersion": "kustomize.config.k8s.io/v1beta1",
		"resources":  resources,
	}
}
