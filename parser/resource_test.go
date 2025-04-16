package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestResource_Kind(t *testing.T) {
	tests := []struct {
		name    string
		r       Resource
		want    string
		wantErr error
	}{
		{
			name: "valid kind",
			r: Resource{
				"kind": "Deployment",
			},
			want:    "Deployment",
			wantErr: nil,
		},
		{
			name:    "missing kind",
			r:       Resource{},
			want:    "",
			wantErr: ErrKindKeyNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Kind()
			if err != tt.wantErr {
				t.Errorf("Resource.Kind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Resource.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newKustomizeResource(t *testing.T) {
	type args struct {
		resources []string
	}
	tests := []struct {
		name string
		args args
		want Resource
	}{
		{
			name: "creates kustomization with sorted resources",
			args: args{
				resources: []string{"b.yaml", "a.yaml", "c.yaml"},
			},
			want: Resource{
				"kind":       "Kustomization",
				"apiVersion": "kustomize.config.k8s.io/v1beta1",
				"resources":  []string{"a.yaml", "b.yaml", "c.yaml"},
			},
		},
		{
			name: "empty resources",
			args: args{
				resources: []string{},
			},
			want: Resource{
				"kind":       "Kustomization",
				"apiVersion": "kustomize.config.k8s.io/v1beta1",
				"resources":  []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newKustomizeResource(tt.args.resources...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newKustomizeResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readResource(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Resource
	}{
		{
			name: "single resource",
			input: `kind: Deployment
apiVersion: apps/v1`,
			want: []Resource{
				{
					"kind":       "Deployment",
					"apiVersion": "apps/v1",
				},
			},
		},
		{
			name: "multiple resources",
			input: `kind: Deployment
apiVersion: apps/v1
---
kind: Service
apiVersion: v1`,
			want: []Resource{
				{
					"kind":       "Deployment",
					"apiVersion": "apps/v1",
				},
				{
					"kind":       "Service",
					"apiVersion": "v1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readResource(strings.NewReader(tt.input)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourcesToMap(t *testing.T) {
	tests := []struct {
		name      string
		resources []Resource
		want      map[string][]Resource
	}{
		{
			name: "group by kind",
			resources: []Resource{
				{"kind": "Deployment", "name": "app1"},
				{"kind": "Service", "name": "svc1"},
				{"kind": "Deployment", "name": "app2"},
			},
			want: map[string][]Resource{
				"Deployment": {
					{"kind": "Deployment", "name": "app1"},
					{"kind": "Deployment", "name": "app2"},
				},
				"Service": {
					{"kind": "Service", "name": "svc1"},
				},
			},
		},
		{
			name: "skip invalid resources",
			resources: []Resource{
				{"kind": "Deployment", "name": "app1"},
				{"name": "invalid"},
				{"kind": "Service", "name": "svc1"},
			},
			want: map[string][]Resource{
				"Deployment": {
					{"kind": "Deployment", "name": "app1"},
				},
				"Service": {
					{"kind": "Service", "name": "svc1"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourcesToMap(tt.resources); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourcesToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
