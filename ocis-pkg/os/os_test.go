package os

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_mustUserConfigDir(t *testing.T) {
	configDir, _ := os.UserConfigDir()
	type args struct {
		prefix    string
		extension string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "fetch the default config location for the current user",
			args: args{
				prefix:    "ocis",
				extension: "testing",
			},
			want: filepath.Join(configDir, "ocis", "testing"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustUserConfigDir(tt.args.prefix, tt.args.extension); got != tt.want {
				t.Errorf("MustUserConfigDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
