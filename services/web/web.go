package web

import (
	"embed"
)

// Assets holds the embedded asset fs.
//
//go:generate make generate
//go:embed all:assets/*
var Assets embed.FS
