// +build tools

package tools

// This file follows the recommendation at
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// on how to pin tooling dependencies to a go.mod file.
// This ensures that all systems use the same version of tools in addition to regular dependencies.

import (
	_ "golang.org/x/tools/cmd/goimports"
)
