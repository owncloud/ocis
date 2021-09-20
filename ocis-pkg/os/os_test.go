package os_test

import (
	"os"
	"path/filepath"
	"testing"

	pkgos "github.com/owncloud/ocis/ocis-pkg/os"
)

func TestMustUserConfigDir(t *testing.T) {
	configDir, _ := os.UserConfigDir()
	type args struct {
		prefix    string
		extension string
	}
	tests := []struct {
		name      string
		args      args
		want      string
		resetHome bool
		panic     bool
	}{
		{
			name: "fetch the default config location for the current user",
			args: args{
				prefix:    "ocis",
				extension: "testing",
			},
			want: filepath.Join(configDir, "ocis", "testing"),
		},
		{
			name: "location cannot be determined because $HOME is not set",
			args: args{
				prefix:    "ocis",
				extension: "testing",
			},
			want:      filepath.Join(configDir, "ocis", "testing"),
			resetHome: true,
			panic:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resetHome {
				unsetHome(t)
			}

			defer func() {
				if r := recover(); r != nil && !tt.panic {
					t.Errorf("should have panicked but didn't")
				}
			}()

			if got := pkgos.MustUserConfigDir(tt.args.prefix, tt.args.extension); got != tt.want {
				t.Errorf("MustUserConfigDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func unsetHome(t *testing.T) {
	if err := os.Setenv("HOME", ""); err != nil {
		t.Error(err)
	}
}
