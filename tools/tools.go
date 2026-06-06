//go:build tools
// +build tools

package tools

import (
	// github.com/AlekSi/gocov-xml is package main (no importable API).
	// This blank import is syntactically valid for go mod tidy pinning only —
	// do not build with -tags tools.
	_ "github.com/AlekSi/gocov-xml"
	_ "github.com/axw/gocov"
)
