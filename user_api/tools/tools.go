//go:build tools
// +build tools

package tools

// This file ensures tool dependencies are kept in sync.  This is the
// recommended way of doing this according to
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// To install the following tools at the version used by this repo run:
// $ make bootstrap
// or
// $ go generate -tags tools tools/tools.go

//go:generate go install github.com/golang/mock/mockgen
import (
	_ "github.com/golang/mock/mockgen"
)
