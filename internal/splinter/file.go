package splinter

import (
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

func getFileNameFromPath(path string) string {
	if !strings.HasSuffix(path, ".yaml") {
		return ""
	}

	split := strings.Split(path, "/")
	if len(split) == 0 {
		return path
	}

	return split[len(split)-1]
}
