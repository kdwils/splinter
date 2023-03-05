package parser

import "testing"

func TestResource_Kind(t *testing.T) {
	tests := []struct {
		name    string
		r       Resource
		want    string
		wantErr bool
	}{
		{
			name: "resource kind key exists",
			r: Resource{
				"apiVersion": "v1",
				"kind":       "Namespace",
				"metadata": Resource{
					"name": "my-namespace",
				},
			},
			want:    "Namespace",
			wantErr: false,
		},
		{
			name: "resource kind key does not exist",
			r: Resource{
				"apiVersion": "v1",
				"metadata": Resource{
					"name": "my-namespace",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Kind()
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.Kind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Resource.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}
