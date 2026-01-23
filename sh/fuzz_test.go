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

package sh

import (
	"errors"
	"testing"
)

// FuzzEnvExpand tests environment variable expansion in command arguments.
// It passes fuzzed strings containing potential $VAR patterns to verify
// that expansion doesn't panic.
func FuzzEnvExpand(f *testing.F) {
	// Seed corpus with various expansion patterns
	seeds := []string{
		// Standard patterns
		"$FOO",
		"${FOO}",
		"$FOO_BAR",
		"$FOO$BAR",
		"${FOO}${BAR}",

		// With surrounding text
		"prefix$FOO",
		"$FOOsuffix",
		"prefix${FOO}suffix",
		"prefix$FOO$BARsuffix",

		// Edge cases
		"$",
		"$$",
		"${",
		"${}",
		"${FOO",
		"${FOO}extra}",
		"$123",
		"$_",
		"$_FOO",

		// No variables
		"no-vars-here",
		"",
		"plain text",

		// Special characters
		"$FOO$",
		"${FOO:-default}",
		"${FOO:=default}",
		"${FOO:+alternate}",
		"${FOO:?error}",

		// Nested (not valid but shouldn't panic)
		"${${FOO}}",
		"$${FOO}",

		// Unicode
		"$FOO_\u00e9",
		"\u00e9$FOO",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Use 'true' command which is safe, fast, and universally available
		// We're testing that the environment expansion doesn't panic
		// The actual command execution is secondary
		_, _ = Output("true", input)
	})
}

// FuzzExitStatus tests exit status extraction from various error types.
// It verifies that ExitStatus and CmdRan handle arbitrary inputs without panicking.
func FuzzExitStatus(f *testing.F) {
	// Seed with various integer values
	f.Add(0)
	f.Add(1)
	f.Add(-1)
	f.Add(127)
	f.Add(128)
	f.Add(255)
	f.Add(256)
	f.Add(-127)
	f.Add(2147483647)  // Max int32
	f.Add(-2147483648) // Min int32

	f.Fuzz(func(t *testing.T, _ int) {
		// Test with nil error - should always return 0
		if status := ExitStatus(nil); status != 0 {
			t.Errorf("ExitStatus(nil) = %d, want 0", status)
		}

		// Test CmdRan with nil - should always return true
		if !CmdRan(nil) {
			t.Error("CmdRan(nil) should be true")
		}

		// Test with generic error - should return 1
		genericErr := errors.New("generic error")
		if status := ExitStatus(genericErr); status != 1 {
			t.Errorf("ExitStatus(generic) = %d, want 1", status)
		}

		// Test CmdRan with generic error - should return false
		// (generic errors are unrecognized types, meaning command didn't run)
		if CmdRan(genericErr) {
			t.Error("CmdRan(generic) should be false for unrecognized error types")
		}
	})
}
