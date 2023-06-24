package splinter

import (
	"reflect"
	"testing"

	"github.com/kdwils/splinter/pkg/parser"
)

func TestInput_removeExclusions(t *testing.T) {
	type fields struct {
		InputFiles []string
		Exclusions []string
		Kustomize  bool
		OutputPath string
	}
	tests := []struct {
		name      string
		fields    fields
		wantFiles []string
	}{
		{
			name: "remove exlcusions",
			fields: fields{
				InputFiles: []string{"file1.yaml", "file2.yaml", "file3.yaml"},
				Exclusions: []string{"file3.yaml"},
			},
			wantFiles: []string{"file1.yaml", "file2.yaml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Input{
				InputFiles: tt.fields.InputFiles,
				Exclusions: tt.fields.Exclusions,
				Kustomize:  tt.fields.Kustomize,
				OutputPath: tt.fields.OutputPath,
			}
			i.removeExclusions()
		})
	}
}

func Test_removeIndex(t *testing.T) {
	type args struct {
		slice []string
		index int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "remove index",
			args: args{
				slice: []string{"1", "2", "3", "4"},
				index: 2,
			},
			want: []string{"1", "2", "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeIndex(tt.args.slice, tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInput_kustomizeFilePath(t *testing.T) {
	type fields struct {
		InputFiles []string
		Exclusions []string
		Kustomize  bool
		OutputPath string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "path is a dir",
			fields: fields{
				OutputPath: "/my/path/",
			},
			want: "/my/path/",
		},
		{
			name: "path is file in current wd",
			fields: fields{
				OutputPath: "my-file.yaml",
			},
			want: ".",
		},
		{
			name: "path is dir with file",
			fields: fields{
				OutputPath: "/path/to/my/file.yaml",
			},
			want: "/path/to/my/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Input{
				InputFiles: tt.fields.InputFiles,
				Exclusions: tt.fields.Exclusions,
				Kustomize:  tt.fields.Kustomize,
				OutputPath: tt.fields.OutputPath,
			}
			if got := i.kustomizeFilePath(); got != tt.want {
				t.Errorf("Input.kustomizeFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kustomization(t *testing.T) {
	type args struct {
		resources []string
	}
	tests := []struct {
		name string
		args args
		want parser.Resource
	}{
		{
			name: "kustomization",
			args: args{
				resources: []string{"resource-1.yaml", "resource-2.yaml"},
			},
			want: parser.Resource{
				"kind":       "Kustomization",
				"apiVersion": "kustomize.config.k8s.io/v1beta1",
				"resources":  []string{"resource-1.yaml", "resource-2.yaml"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kustomization(tt.args.resources...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("kustomization() = %v, want %v", got, tt.want)
			}
		})
	}
}
