# AGENTS.md - AI Coding Assistant Context

This file provides context for AI coding assistants (Claude Code, OpenCode, Kiro.dev, etc.) working with this repository.

## Repository Overview

**Project:** Mage - A make/rake-like build tool using Go
**Fork:** This is a fork of [magefile/mage](https://github.com/magefile/mage), maintained by Actualyze AI
**Module Path:** `github.com/actualyze-ai/mage`
**Go Version:** 1.25
**License:** Apache 2.0

### Why This Fork Exists

The upstream mage project (github.com/magefile/mage) has become dormant:
- Last release: v1.15.0 (May 2023)
- Still targets Go 1.12 (released February 2019)
- No security tooling (no linting, no vulnerability scanning, no SBOM)
- Uses deprecated `ioutil` package
- Tests against Go versions as old as 1.11.x

This fork modernizes mage for current Go versions and adds security best practices:
- Go 1.25 support with modern toolchain
- Comprehensive linting with golangci-lint v2
- Security scanning (govulncheck, Grype)
- SBOM generation for releases
- Dependabot for dependency monitoring
- 65% code coverage threshold

---

## Critical Workflow Rules

### 1. Never Commit Directly to Master
**ALWAYS** create a branch and propose changes via Pull Request. Never push directly to `master`.

### 2. Branch Naming Convention
```
feature/<issue-number>-<issue-subject-with-dashes>
fix/<issue-number>-<issue-subject-with-dashes>
chore/<description-with-dashes>
```

Examples:
- `feature/168-milestone-2-linting-security`
- `fix/170-broken-test-on-windows`
- `chore/add-agents-md`
- `chore/update-dependencies`

### 3. Issue Tracking
Issues are tracked in **this repository**:
- **Issues:** https://github.com/actualyze-ai/mage/issues
- **Project Board:** https://github.com/orgs/actualyze-ai/projects/8
- Reference issues as: `#<number>`

When creating a new issue, add it to the [Actualyze AI project board](https://github.com/orgs/actualyze-ai/projects/8).

### 4. Keep Documentation Updated
After making changes, update:
- This file (`AGENTS.md`) if architecture or workflows change
- `NOTICE` file with high-level summary of modifications

### 5. Required File Headers (Fork Compliance)

Since this is a fork of an Apache 2.0 licensed project, every file we modify or create MUST have appropriate headers.

**For MODIFIED Go files:**
```go
// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright (c) 2026 Actualyze AI
//
// NOTE: This file has been modified by Actualyze AI from the original upstream
// version (magefile/mage). See git history for details.

package example
```

**For NEW Go files:**
```go
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

package example
```

**For MODIFIED YAML files (workflows, configs):**
```yaml
# SPDX-License-Identifier: Apache-2.0
# Modifications Copyright (c) 2026 Actualyze AI
#
# NOTE: This file has been modified by Actualyze AI from the original upstream
# version (magefile/mage). See git history for details.
```

**For NEW YAML files:**
```yaml
# SPDX-License-Identifier: Apache-2.0
#
# Copyright (c) 2026 Actualyze AI
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
```

**For MODIFIED Markdown files:**
```markdown
<!--
SPDX-License-Identifier: Apache-2.0
Modifications Copyright (c) 2026 Actualyze AI

NOTE: This file has been modified by Actualyze AI from the original upstream
version (magefile/mage). See git history for details.
-->
```

**For NEW Markdown files:**
```markdown
<!--
SPDX-License-Identifier: Apache-2.0

Copyright (c) 2026 Actualyze AI

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->
```

### 6. Run CI Checks Before Committing

**ALWAYS** run these checks locally before creating a PR:

```bash
# Linting (required to pass)
golangci-lint run ./...

# Tests with race detection (required to pass)
go test -tags CI -race ./...

# Verify code compiles
go build ./...
```

The CI pipeline will run these same checks and the PR cannot be merged if they fail.

---

## Project Architecture

### What is Mage?

Mage is a build tool that lets you write build automation in Go instead of Makefiles or shell scripts. You write regular Go functions with a `//go:build mage` tag, and mage compiles them into a runnable binary.

Example magefile:
```go
//go:build mage

package main

import "github.com/actualyze-ai/mage/mg"
import "github.com/actualyze-ai/mage/sh"

// Build compiles the project
func Build() error {
    mg.Deps(Generate) // Run Generate first
    return sh.Run("go", "build", "./...")
}

// Generate runs code generation
func Generate() error {
    return sh.Run("go", "generate", "./...")
}
```

### Package Structure

```
github.com/actualyze-ai/mage/
├── main.go           # CLI entry point - calls mage.Main()
├── magefile.go       # Mage's own build script (dogfooding)
├── bootstrap.go      # Bootstrap installer
│
├── mage/             # Core orchestration package
│   ├── main.go       # Main(), Invoke(), Parse(), Compile()
│   ├── template.go   # Generated mainfile template
│   └── testdata/     # Test fixtures
│
├── mg/               # Runtime utilities for magefile authors
│   ├── deps.go       # Deps(), CtxDeps(), SerialDeps()
│   ├── fn.go         # Fn interface, F() wrapper
│   ├── runtime.go    # Verbose(), Debug(), GoCmd(), CacheDir()
│   ├── errors.go     # Fatal(), Fatalf(), ExitStatus()
│   └── color.go      # Terminal color support
│
├── sh/               # Shell command execution
│   ├── cmd.go        # Run(), RunV(), Output(), Exec()
│   └── helpers.go    # Rm(), Copy()
│
├── parse/            # Go AST parsing for target extraction
│   ├── parse.go      # PrimaryPackage(), Package()
│   └── testdata/     # Test fixtures
│
├── target/           # File freshness checking (like make)
│   ├── target.go     # Path(), Glob(), Dir()
│   └── newer.go      # PathNewer(), OldestModTime(), NewestModTime()
│
├── internal/         # Internal utilities (not for external use)
│   └── run.go        # RunDebug(), OutputDebug()
│
└── site/             # Documentation website source
```

### Package Dependencies

```
main.go
└── mage/
    ├── internal/
    ├── mg/
    ├── parse/
    │   └── internal/
    └── sh/
        └── mg/

User magefiles import:
├── mg/      # Dependency management, configuration
├── sh/      # Command execution
└── target/  # Freshness checking
```

### Key Types and Interfaces

**`mg.Fn` - Dependency Function Interface:**
```go
type Fn interface {
    Name() string                    // Fully qualified function name
    ID() string                      // Uniqueness identifier (for args)
    Run(ctx context.Context) error   // Execute the function
}
```

**`mg.Namespace` - Target Grouping:**
```go
type Namespace struct{}

// Types embedding Namespace have methods as namespaced targets
type Build mg.Namespace

func (Build) Docker() error { ... }  // Called as: mage build:docker
```

**Valid Target Function Signatures:**
```go
func()
func() error
func(context.Context)
func(context.Context) error
// With optional typed arguments:
func(ctx context.Context, name string, count int) error
// Supported arg types: string, int, float64, bool, time.Duration
```

### Execution Flow

1. **Entry:** `main.go` → `mage.Main()`
2. **Parse CLI:** Flags parsed into `Invocation` struct
3. **Find Magefiles:** Look for `//go:build mage` tagged files
4. **Parse Targets:** Use Go AST to extract exported functions
5. **Check Cache:** Look for compiled binary in `~/.magefile/`
6. **Generate Main:** Create temporary wrapper with target dispatch
7. **Compile:** Run `go build` on magefiles + generated main
8. **Execute:** Run compiled binary with user's target arguments
9. **Dependency Resolution:** `mg.Deps()` runs deps in parallel, exactly once

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MAGEFILE_VERBOSE` | `false` | Print target names as they run |
| `MAGEFILE_DEBUG` | `false` | Print detailed debug information |
| `MAGEFILE_GOCMD` | `go` | Go command to use |
| `MAGEFILE_CACHE` | `~/.magefile` | Binary cache directory |
| `MAGEFILE_HASHFAST` | `false` | Skip content hashing for faster rebuilds |
| `MAGEFILE_IGNOREDEFAULT` | `false` | Ignore default target |
| `MAGEFILE_ENABLE_COLOR` | auto | Enable/disable color output |
| `MAGEFILE_TARGET_COLOR` | `cyan` | Color for target names |

---

## CI/CD Configuration

### CI Workflow (`.github/workflows/ci.yml`)

Runs on every push and pull request:

**Lint Job:**
- golangci-lint v2 with comprehensive rule set
- Uses `.golangci.yml` configuration

**Test Job (matrix: stable + go.mod version):**
- `go test -tags CI -race ./...`
- Coverage threshold: 65% minimum
- Coverage artifacts uploaded

### Release Workflow (`.github/workflows/release.yml`)

Triggered by `v*` tags (e.g., `v1.16.0`):
- Uses goreleaser v2
- Builds: linux/amd64, darwin/arm64
- Generates SBOM in SPDX JSON format
- Creates GitHub Release with binaries and checksums

### Security Workflows (`.github/workflows/security.yml`)

- govulncheck for Go vulnerability scanning
- Grype for container/dependency scanning
- SARIF upload to GitHub Security tab
- OSSF Scorecard for supply chain security

### Dependabot (`.github/dependabot.yml`)

- Weekly updates for Go modules (Mondays)
- Weekly updates for GitHub Actions (Mondays)

---

## Common Development Tasks

### Running Tests
```bash
# Standard test run
go test -tags CI -race ./...

# With verbose output
go test -tags CI -race -v ./...

# Specific package
go test -tags CI -race ./mage/...

# With coverage
go test -tags CI -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Running Linter
```bash
# Run all linters
golangci-lint run ./...

# With auto-fix where possible
golangci-lint run --fix ./...

# Specific linters only
golangci-lint run --enable=gosec,govet ./...
```

### Building Mage
```bash
# Install to GOBIN with version info
go run bootstrap.go

# Or using mage itself (if already installed)
mage install

# Build without version info
go build -o mage .

# Build with custom version info
go build -ldflags="-X github.com/actualyze-ai/mage/mage.gitTag=dev" -o mage .
```

### Creating a Release
Releases are automated via GitHub Actions:
```bash
# Create and push an annotated tag
git tag -a v1.17.0 -m "v1.17.0

Release notes here..."
git push origin v1.17.0
```

---

## Linting Configuration

The `.golangci.yml` uses golangci-lint v2 format with `default: all` and specific exclusions. Key settings:

**Enabled by default:** All linters except those explicitly disabled

**Notable disabled linters (with rationale):**
- `cyclop`, `gocognit`, `gocyclo`, `funlen` - Complexity limits are arbitrary
- `err113` - Would require significant refactor for sentinel errors
- `exhaustruct` - Not all struct fields need explicit initialization
- `forbidigo` - fmt.Print is acceptable for CLI tools
- `wsl`, `wsl_v5` - Too opinionated about whitespace
- `noctx` - exec.Command without context is acceptable for build tools

**gosec exclusions:**
- `G104` - Unhandled errors acceptable for debug/cleanup
- `G204` - Subprocess with variable acceptable for build tools
- `G304` - File path as taint input acceptable for build tools

---

## Related Resources

| Resource | URL |
|----------|-----|
| Issues | https://github.com/actualyze-ai/mage/issues |
| This Repository | https://github.com/actualyze-ai/mage |
| Upstream Mage | https://github.com/magefile/mage |
| Mage Documentation | https://magefile.org |
| Releases | https://github.com/actualyze-ai/mage/releases |
| Go Package Docs | https://pkg.go.dev/github.com/actualyze-ai/mage |

---

## Fork Modifications Summary

See `NOTICE` file for the complete list. Key changes from upstream:

1. **Go Version:** 1.12 → 1.25
2. **Module Path:** `github.com/magefile/mage` → `github.com/actualyze-ai/mage`
3. **Deprecated APIs:** Replaced `ioutil` with `os`/`io` equivalents
4. **Linting:** Added golangci-lint v2 with comprehensive rules
5. **Security:** Added govulncheck, Grype, Dependabot, OSSF Scorecard
6. **Releases:** Added goreleaser automation with SBOM generation
7. **Coverage:** Added 65% threshold enforcement
8. **CI:** Modernized GitHub Actions (v6 actions, Go 1.25)
