//go:build tools
// +build tools

// tools.go serves as a dependedncy management file to bring in this module external go tools.
// this allows to run ocb and mdatagen without having them installed via //go:generate comments.
// see https://www.jvt.me/posts/2022/06/15/go-tools-dependency-management/ for more explainations.

package otelcolcustom

import (
	_ "github.com/mikefarah/yq/v4"
	_ "go.opentelemetry.io/collector/cmd/builder"

	_ "github.com/clement-casse/playground/otelcol-custom/exporter/cyphergraphexporter"
)
