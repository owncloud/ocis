package web

import (
	"embed"
)

//go:embed assets/*
//go:embed assets/js/*
var Assets embed.FS
