package parser

import (
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ListFilesFromInput(inputs []string) ([]string, error) {
	files := make([]string, 0)
	for _, p := range inputs {
		if strings.EqualFold(filepath.Ext(p), ".yaml") {
			files = append(files, p)
			continue
		}

		fileInfo, err := os.Stat(p)
		if err != nil {
			return nil, err
		}

		if !fileInfo.IsDir() {
			continue
		}

		dir, err := os.ReadDir(p)
		if err != nil {
			return nil, err
		}

		for _, file := range dir {
			files = append(files, path.Join(p, file.Name()))
		}
	}

	return files, nil
}

func FileToBuffer(files ...string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)

	for _, f := range files {
		copyBytes(f, buf)
	}

	return buf, nil
}

func copyBytes(f string, buf *bytes.Buffer) error {
	o, err := os.Open(f)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = io.Copy(buf, o)
	if err != nil {
		return err
	}

	return nil
}
