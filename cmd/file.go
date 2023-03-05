package cmd

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func filesFromInput(input []string) []string {
	files := make([]string, 0)
	for _, p := range input {
		if strings.EqualFold(filepath.Ext(p), ".yaml") {
			files = append(files, p)
			continue
		}

		fileInfo, err := os.Stat(p)
		if err != nil {
			continue
		}

		if !fileInfo.IsDir() {
			continue
		}

		dir, err := os.ReadDir(p)
		if err != nil {
			continue
		}

		for _, file := range dir {
			files = append(files, path.Join(p, file.Name()))
		}
	}

	return files
}

func readFile(file string) (*bytes.Buffer, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

func createFile(path string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}

	return os.Create(path)
}
