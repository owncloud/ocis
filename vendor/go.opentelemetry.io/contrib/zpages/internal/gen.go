// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/contrib/zpages/internal"

import "embed"

//go:embed templates/*
var Templates embed.FS
