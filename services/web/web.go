package web

import (
	"embed"
)

//go:generate make generate

//go:embed all:assets/*
var Assets embed.FS

//go:embed all:apps/*
var Apps embed.FS
