// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright (c) 2026 Actualyze AI
//
// NOTE: This file has been modified by Actualyze AI from the original upstream
// version (magefile/mage). See git history for details.

package mage

import (
	"fmt"

	"github.com/actualyze-ai/mage/mg"
)

// BuildSubdir2 Builds stuff.
func BuildSubdir2() {
	fmt.Println("buildsubdir2")
}

// NS is a namespace.
type NS mg.Namespace

// Deploy2 deploys stuff.
func (NS) Deploy2() {
	fmt.Println("deploy2")
}
