package settings

import (
	"embed"
)

// Assets holds the embedded asset fs.
//
//go:generate make generate
//go:embed assets/*
var Assets embed.FS
