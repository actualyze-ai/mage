//go:build ignore
// +build ignore

// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright (c) 2026 Actualyze AI
//
// NOTE: This file has been modified by Actualyze AI from the original upstream
// version (magefile/mage). See git history for details.

package main

import (
	"os"

	"github.com/actualyze-ai/mage/mage"
)

// This is a bootstrap builder, to build mage when you don't already *have* mage.
// Run it like
// go run bootstrap.go
// and it will install mage with all the right flags created for you.

func main() {
	os.Args = []string{os.Args[0], "-v", "install"}
	os.Exit(mage.Main())
}
