//go:build mage
// +build mage

// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright (c) 2026 Actualyze AI
//
// NOTE: This file has been modified by Actualyze AI from the original upstream
// version (magefile/mage). See git history for details.

package main

import "github.com/actualyze-ai/mage/mage/testdata/mixed_lib_files/subdir"

func Build() {
	subdir.Build()
}
