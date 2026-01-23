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

package mage

import (
	"bytes"
	"testing"
)

// FuzzParse tests CLI argument parsing for panics and crashes.
// It splits fuzz data into argument slices and passes them to Parse().
func FuzzParse(f *testing.F) {
	// Seed corpus with valid CLI patterns
	seeds := []string{
		"-v",
		"-debug",
		"-h",
		"--help",
		"-version",
		"-init",
		"-clean",
		"-t 30s",
		"-t 1ms",
		"-t 5m",
		"-d /tmp/test",
		"-w /tmp/work",
		"-gocmd go",
		"-gocmd /usr/local/go/bin/go",
		"-f magefile.go",
		"-goarch amd64",
		"-goos linux",
		"-ldflags -s -w",
		"-compile ./output",
		"build",
		"build deploy",
		"test:unit",
		"-v build",
		"-v -debug build",
		"-t 1ms build",
		"-d . -v build",
		"",
		"-l",
		"target arg1 arg2",
		"target:subtarget",
		"ns:target arg",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data string) {
		// Split data into args
		args := splitArgs(data)

		stderr := &bytes.Buffer{}
		stdout := &bytes.Buffer{}

		// Parse should not panic regardless of input
		_, _, _ = Parse(stderr, stdout, args)
	})
}

// splitArgs converts a string into an argument slice by splitting on whitespace.
// This simulates how shell arguments would be passed to the program.
func splitArgs(data string) []string {
	var args []string
	var current []byte

	for i := 0; i < len(data); i++ {
		b := data[i]
		if b == ' ' || b == '\n' || b == '\t' || b == '\r' {
			if len(current) > 0 {
				args = append(args, string(current))
				current = nil
			}
		} else {
			current = append(current, b)
		}
	}
	if len(current) > 0 {
		args = append(args, string(current))
	}
	return args
}
