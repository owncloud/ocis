package config

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9189",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: GRPC{
			Addr:      "127.0.0.1:9185",
			Namespace: "com.owncloud.api",
		},
		Service: Service{
			Name: "thumbnails",
		},
		Thumbnail: Thumbnail{
			Resolutions: []string{"16x16", "32x32", "64x64", "128x128", "1920x1080", "3840x2160", "7680x4320"},
			FileSystemStorage: FileSystemStorage{
				RootDirectory: path.Join(defaults.BaseDataPath(), "thumbnails"),
			},
			WebdavAllowInsecure: true,
			RevaGateway:         "127.0.0.1:9142",
			CS3AllowInsecure:    false,
		},
	}
}
