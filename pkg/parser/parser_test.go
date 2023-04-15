package parser

import (
	"bytes"
	"reflect"
	"testing"
)

func testFlatBytes() []byte {
	s := `apiVersion: v1
kind: Namespace
metadata:
  name: my-namespace
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-service-account
  namespace: my-namespace
`
	return []byte(s)
}

func TestParser_Split(t *testing.T) {
	type fields struct {
		Resources []Resource
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string][]Resource
	}{
		{
			name: "split",
			fields: fields{
				Resources: []Resource{
					{
						"apiVersion": "v1",
						"kind":       "Namespace",
						"metadata": Resource{
							"name": "my-namespace",
						},
					},
					{
						"apiVersion": "v1",
						"kind":       "ServiceAccount",
						"metadata": Resource{
							"name":      "my-service-account",
							"namespace": "my-namespace",
						},
					},
				},
			},
			want: map[string][]Resource{
				"Namespace": {
					{
						"apiVersion": "v1",
						"kind":       "Namespace",
						"metadata": Resource{
							"name": "my-namespace",
						},
					},
				},
				"ServiceAccount": {
					{
						"apiVersion": "v1",
						"kind":       "ServiceAccount",
						"metadata": Resource{
							"name":      "my-service-account",
							"namespace": "my-namespace",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Resources: tt.fields.Resources,
			}
			if got := p.Sort(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Write(t *testing.T) {
	type args struct {
		resources []Resource
	}
	tests := []struct {
		name       string
		p          *Parser
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name: "writer",
			args: args{
				resources: []Resource{
					{
						"apiVersion": "v1",
						"kind":       "Namespace",
						"metadata": Resource{
							"name": "my-namespace",
						},
					},
				},
			},
			wantErr: false,
			wantWriter: `apiVersion: v1
kind: Namespace
metadata:
  name: my-namespace
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			writer := &bytes.Buffer{}
			if err := p.Write(writer, tt.args.resources...); (err != nil) != tt.wantErr {
				t.Errorf("Parser.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Parser.Write() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestParser_YamlFile(t *testing.T) {
	type fields struct {
		Resources []Resource
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "yaml file name",
			args: args{
				name: "my-FILE",
			},
			want: "my-file.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := YamlFileName(tt.args.name); got != tt.want {
				t.Errorf("Parser.YamlFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_kustomizeResourcePaths(t *testing.T) {
	type fields struct {
		Resources []Resource
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "kustomize resource paths",
			fields: fields{
				Resources: []Resource{
					{
						"kind": "Kustomization",
					},
					{
						"kind": "Service",
					},
					{
						"kind": "Service",
					},
					{
						"kind": "Deployment",
					},
				},
			},
			want: []string{"deployment.yaml", "service.yaml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Resources: tt.fields.Resources,
			}
			if got := p.kustomizeResourcePaths(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.kustomizeResourcePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Kustomize(t *testing.T) {
	type fields struct {
		Resources []Resource
	}
	tests := []struct {
		name   string
		fields fields
		want   Resource
	}{
		{
			name: "kustomize",
			fields: fields{
				Resources: []Resource{
					{
						"kind": "Kustomization",
					},
					{
						"kind": "Service",
					},
					{
						"kind": "Service",
					},
					{
						"kind": "Deployment",
					},
				},
			},
			want: Resource{
				"kind":       "Kustomization",
				"apiVersion": "kustomize.config.k8s.io/v1beta1",
				"resources":  []string{"deployment.yaml", "service.yaml"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Resources: tt.fields.Resources,
			}
			if got := p.Kustomize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Kustomize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Parser
	}{
		{
			name: "new",
			want: &Parser{
				Resources: make([]Resource, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
