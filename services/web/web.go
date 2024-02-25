package web

import (
	"embed"
)

//go:generate make generate

//go:embed all:assets/web/*
var WebAssets embed.FS

//go:embed all:assets/apps/*
var AppAssets embed.FS
