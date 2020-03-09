// +build tools

package main

import (
	_ "github.com/haya14busa/goverage"
	_ "github.com/restic/calens"
	_ "golang.org/x/lint/golint"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
