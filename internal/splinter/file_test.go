package splinter

import "testing"

func Test_getFileNameFromPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no file provided",
			args: args{
				"/my/path/",
			},
			want: "",
		},
		{
			name: "single file, no path",
			args: args{
				path: "my-file.yaml",
			},
			want: "my-file.yaml",
		},
		{
			name: "file at end of path",
			args: args{
				"/path/to/my/file.yaml",
			},
			want: "file.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFileNameFromPath(tt.args.path); got != tt.want {
				t.Errorf("getFileNameFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
