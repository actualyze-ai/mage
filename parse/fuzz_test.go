// SPDX-License-Identifier: Apache-2.0
//
// Copyright (c) 2026 Actualyze AI
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parse

import (
	"os"
	"path/filepath"
	"testing"
)

// FuzzParseGoSource tests parsing of Go source files with various function signatures.
// It creates temporary Go files with fuzzed content and attempts to parse them,
// looking for panics rather than errors.
func FuzzParseGoSource(f *testing.F) {
	// Seed corpus with valid Go source patterns
	seeds := []string{
		// Simple function
		`package main

func Simple() {}`,

		// Function with error return
		`package main

func WithError() error { return nil }`,

		// Function with context parameter
		`package main

import "context"

func WithContext(ctx context.Context) error { return nil }`,

		// Namespace method
		`package main

import "github.com/actualyze-ai/mage/mg"

type Build mg.Namespace

func (Build) Target() error { return nil }`,

		// Function with typed arguments
		`package main

func TakesArgs(name string, count int) error { return nil }`,

		// Function with duration argument
		`package main

import "time"

func TakesDuration(d time.Duration) {}`,

		// Function with float64 argument
		`package main

func TakesFloat(val float64) error { return nil }`,

		// Function with bool argument
		`package main

func TakesBool(flag bool) {}`,

		// Multiple functions
		`package main

func First() error { return nil }
func Second() {}
func Third(s string) error { return nil }`,

		// With default target
		`package main

var Default = Build

func Build() error { return nil }`,

		// With aliases
		`package main

var Aliases = map[string]interface{}{
	"b": Build,
}

func Build() error { return nil }`,

		// Empty package
		`package main`,

		// Only imports
		`package main

import (
	"fmt"
	"os"
)`,

		// Private functions only (should be skipped)
		`package main

func private() error { return nil }`,

		// Mixed public and private
		`package main

func Public() error { return nil }
func private() {}`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, source string) {
		// Create temp directory with Go file
		dir := t.TempDir()
		goFile := filepath.Join(dir, "magefile.go")

		// Prepend build tag if not present
		content := "//go:build mage\n\n" + source
		if err := os.WriteFile(goFile, []byte(content), 0o644); err != nil {
			return // Invalid write, skip
		}

		// Try to parse - we're looking for panics, not errors
		// Errors are expected for invalid Go source
		_, _ = Package(dir, []string{"magefile.go"})
	})
}
