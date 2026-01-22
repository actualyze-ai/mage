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
	"os"
	"strings"
)

func TestWorkingDir() error {
	files, err := os.ReadDir(".")
	if err != nil {
		return err
	}
	var out []string
	for _, f := range files {
		out = append(out, f.Name())
	}

	fmt.Println(strings.Join(out, ", "))
	return nil
}
