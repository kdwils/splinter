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
			p := &Parser{
				IndentSize: 2,
			}
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

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Parser
	}{
		{
			name: "new",
			want: &Parser{
				Resources:  make([]Resource, 0),
				IndentSize: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(WithIndentSize(2)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Sanitize(t *testing.T) {
	type fields struct {
		Resources  []Resource
		IndentSize int
	}
	tests := []struct {
		name          string
		fields        fields
		wantResources []Resource
	}{
		{
			name: "sanitize",
			fields: fields{
				Resources: []Resource{
					{
						"kind": "my-resource",
					},
					{
						"not-kind": "remove-me",
					},
				},
			},
			wantResources: []Resource{
				{
					"kind": "my-resource",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Resources:  tt.fields.Resources,
				IndentSize: tt.fields.IndentSize,
			}
			p.Sanitize()
			if !reflect.DeepEqual(tt.wantResources, p.Resources) {
				t.Errorf("TestParser_Sanitize() got = %v, want = %v", p.Resources, tt.wantResources)
			}
		})
	}
}

func TestParser_Remove(t *testing.T) {
	type fields struct {
		Resources  []Resource
		IndentSize int
	}
	type args struct {
		target Resource
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantResources []Resource
	}{
		{
			name: "remove target from resource",
			fields: fields{
				Resources: []Resource{
					{
						"kind": "my-resource",
					},
					{
						"not-kind": "dont-remove-me",
					},
				},
			},
			args: args{
				target: Resource{
					"kind": "my-resource",
				},
			},
			wantResources: []Resource{
				{
					"not-kind": "dont-remove-me",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Resources:  tt.fields.Resources,
				IndentSize: tt.fields.IndentSize,
			}
			p.Remove(tt.args.target)

			if !reflect.DeepEqual(tt.wantResources, p.Resources) {
				t.Errorf("TestParser_Remove() got = %v, want = %v", p.Resources, tt.wantResources)
			}
		})
	}
}
