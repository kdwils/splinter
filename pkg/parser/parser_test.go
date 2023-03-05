package parser

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
		IndentSize int
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string][]Resource
	}{
		{
			name: "split",
			args: args{
				reader: bytes.NewReader(testFlatBytes()),
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
			p := &Parser{}
			if got := p.Sort(tt.args.reader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Flatten(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		p    *Parser
		args args
		want []Resource
	}{
		{
			name: "flatten",
			args: args{
				bytes.NewReader(testFlatBytes()),
			},
			want: []Resource{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			if got := p.Flatten(tt.args.reader); !assert.ElementsMatch(t, got, tt.want) {
				t.Errorf("Parser.Flatten() = %v, want %v", got, tt.want)
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
