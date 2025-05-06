// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/contrib/zpages/internal"

import "embed"

// Templates embeds all the HTML templates used used to serve the tracez
// endpoint
//
//go:embed templates/*
var Templates embed.FS
