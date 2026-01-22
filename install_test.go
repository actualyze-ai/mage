//go:build CI
// +build CI

// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright (c) 2026 Actualyze AI
//
// NOTE: This file has been modified by Actualyze AI from the original upstream
// version (magefile/mage). See git history for details.

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBootstrap(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	s, err := run("go", "run", "bootstrap.go")
	if err != nil {
		t.Fatal(s)
	}
	name := "mage"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	// Use `go env GOBIN` to determine install location, matching bootstrap.go behavior.
	binDir, err := run("go", "env", "GOBIN")
	if err != nil {
		t.Fatalf("failed to get GOBIN: %v", err)
	}
	binDir = filepath.Clean(strings.TrimSpace(binDir))
	if _, err := os.Stat(filepath.Join(binDir, name)); err != nil {
		t.Fatal(err)
	}
}

func run(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	b, err := c.CombinedOutput()
	return string(b), err
}
