//go:build mage
// +build mage

// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright (c) 2026 Actualyze AI
//
// NOTE: This file has been modified by Actualyze AI from the original upstream
// version (magefile/mage). See git history for details.

package main

import (
	"fmt"
	"log"

	"github.com/actualyze-ai/mage/mg"
)

var Default = SomePig

// this should not be a target because it returns a string
func ReturnsString() string {
	fmt.Println("more stuff")
	return ""
}

func TestVerbose() {
	log.Println("hi!")
}

// This is the synopsis for SomePig.  There's more data that won't show up.
func SomePig() {
	mg.Deps(f)
}

func f() {}
