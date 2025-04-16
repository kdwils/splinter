package parser

import (
	"bytes"
	"os"
	"testing"

	"github.com/kdwils/splinter/pkg/fio/mocks"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		p := New()
		if p.indentSize != defaultIndentSize {
			t.Errorf("expected indent size %d, got %d", defaultIndentSize, p.indentSize)
		}
	})

	t.Run("with custom indent size", func(t *testing.T) {
		p := New(WithIndentSize(4))
		if p.indentSize != 4 {
			t.Errorf("expected indent size 4, got %d", p.indentSize)
		}
	})

	t.Run("with custom file io", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFio := mocks.NewMockFileIO(ctrl)
		p := New(WithFileIO(mockFio))
		if p.fio != mockFio {
			t.Error("expected mock file io to be set")
		}
	})
}

func TestParser_Merge(t *testing.T) {
	t.Run("merge from file to output file", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		deploymentYaml, err := os.ReadFile("./testing/deployment.yaml")
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
			t.FailNow()
		}

		serviceYaml, err := os.ReadFile("./testing/service.yaml")
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
			t.FailNow()
		}

		serviceSeparatorYaml, err := os.ReadFile("./testing/serviceseparator.yaml")
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
			t.FailNow()
		}

		mockFio := mocks.NewMockFileIO(ctrl)
		mockFile := mocks.NewMockWriteCloser(ctrl)

		mockFio.EXPECT().ReadFile("./testing/deployment.yaml").Return(deploymentYaml, nil)
		mockFio.EXPECT().ReadFile("./testing/service.yaml").Return(serviceYaml, nil)
		mockFio.EXPECT().Stat(".").Return(nil, os.ErrNotExist)
		mockFio.EXPECT().MkdirAll(".", os.ModePerm).Return(nil)
		mockFio.EXPECT().Create("output.yaml").Return(mockFile, nil)
		mockFile.EXPECT().Write(deploymentYaml).Return(0, nil)
		mockFile.EXPECT().Write(serviceSeparatorYaml).Return(0, nil)
		mockFile.EXPECT().Close().Return(nil)

		p := New(WithFileIO(mockFio))
		err = p.Merge([]string{"./testing/deployment.yaml", "./testing/service.yaml"}, nil, "output.yaml")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func TestParser_Split(t *testing.T) {
	t.Run("split single file with kustomize", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFio := mocks.NewMockFileIO(ctrl)
		mockKustomizeFile := mocks.NewMockWriteCloser(ctrl)
		mockServiceFile := mocks.NewMockWriteCloser(ctrl)
		mockDeployFile := mocks.NewMockWriteCloser(ctrl)

		input, err := os.ReadFile("./testing/input.yaml")
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
			t.FailNow()
		}

		serviceYaml, err := os.ReadFile("./testing/service.yaml")
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
			t.FailNow()
		}

		deploymentYaml, err := os.ReadFile("./testing/deployment.yaml")
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
			t.FailNow()
		}

		kustomizationYaml, err := os.ReadFile("./testing/kustomization.yaml")
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
			t.FailNow()
		}

		mockFio.EXPECT().ReadFile("input.yaml").Return([]byte(input), nil)

		mockFio.EXPECT().Stat("output").Return(nil, os.ErrNotExist)
		mockFio.EXPECT().MkdirAll("output", os.ModePerm).Return(nil)
		mockFio.EXPECT().Create("output/service.yaml").Return(mockServiceFile, nil)
		mockServiceFile.EXPECT().Write(serviceYaml).Return(1, nil)
		mockServiceFile.EXPECT().Close().Return(nil)

		mockFio.EXPECT().Stat("output").Return(nil, os.ErrNotExist)
		mockFio.EXPECT().MkdirAll("output", os.ModePerm).Return(nil)
		mockFio.EXPECT().Create("output/kustomization.yaml").Return(mockKustomizeFile, nil)
		mockKustomizeFile.EXPECT().Write(kustomizationYaml).Return(1, nil)
		mockKustomizeFile.EXPECT().Close().Return(nil)

		mockFio.EXPECT().Stat("output").Return(nil, os.ErrNotExist)
		mockFio.EXPECT().MkdirAll("output", os.ModePerm).Return(nil)
		mockFio.EXPECT().Create("output/deployment.yaml").Return(mockDeployFile, nil)
		mockDeployFile.EXPECT().Write(deploymentYaml).Return(1, nil)
		mockDeployFile.EXPECT().Close().Return(nil)

		p := New(WithFileIO(mockFio))
		err = p.Split([]string{"input.yaml"}, nil, "output", true)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func TestWrite(t *testing.T) {
	t.Run("write resources to writer", func(t *testing.T) {
		buf := new(bytes.Buffer)
		resources := []Resource{
			{
				"kind":       "Service",
				"apiVersion": "v1",
				"metadata": map[string]any{
					"name": "test",
				},
			},
		}

		err := write(buf, 2, resources...)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if buf.Len() == 0 {
			t.Error("expected buffer to contain yaml content")
		}
	})
}
