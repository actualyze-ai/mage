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

# GitHub Workflows & Configuration Guide

This document provides comprehensive guidance for configuring GitHub Actions workflows, rulesets, and code scanning for Actualyze AI projects. It serves as a reference for both humans and AI coding assistants.

---

## Table of Contents

- [Quick Reference Checklist](#quick-reference-checklist)
- [Part 1: Workflow Triggers](#part-1-workflow-triggers)
- [Part 2: Rulesets & Branch Protection](#part-2-rulesets--branch-protection)
- [Part 3: Code Scanning](#part-3-code-scanning)
- [Part 4: Workflow Templates](#part-4-workflow-templates)
- [Part 5: Troubleshooting](#part-5-troubleshooting)
- [Appendix: gh api Commands](#appendix-gh-api-commands)

---

## Quick Reference Checklist

Use this checklist when creating or reviewing GitHub workflow configurations.

### Workflow Triggers

**MUST DO:**

- [ ] Include `merge_group:` trigger if repository uses merge queue
- [ ] Include `workflow_dispatch:` for manual triggering from any branch
- [ ] Use `pull_request: branches: [main]` for CI workflows
- [ ] Add `concurrency` settings to cancel stale runs
- [ ] Use `schedule:` with cron for periodic security scans

**MUST NOT DO:**

- [ ] Do NOT use both `push:` and `pull_request:` on same branches (causes duplicate runs)
- [ ] Do NOT omit `merge_group:` if merge queue is enabled (queue will hang indefinitely)
- [ ] Do NOT use `push:` trigger for CI workflows (use `pull_request:` instead)

### Rulesets & Branch Protection

**MUST DO:**

- [ ] Create "All Jobs" summary job in each workflow for easier ruleset configuration
- [ ] Name summary jobs consistently: `<Workflow> / All Jobs` (e.g., `CI / All Jobs`)
- [ ] Only require status checks from tools that run on PRs

**MUST NOT DO:**

- [ ] Do NOT require code scanning results from scheduled-only tools (e.g., Scorecard)
- [ ] Do NOT require individual job names when a summary job exists

### Code Scanning

**MUST DO:**

- [ ] Upload SARIF only for tools that run on pull requests
- [ ] Delete stale analyses when changing code scanning configuration
- [ ] Use consistent `category` values in SARIF uploads

**MUST NOT DO:**

- [ ] Do NOT upload SARIF to code-scanning for scheduled-only tools (blocks PRs)
- [ ] Do NOT mix GitHub's "Default Setup" with custom CodeQL workflows
- [ ] Do NOT require Scorecard in code scanning (it only runs on default branch)

---

## Part 1: Workflow Triggers

### 1.1 Trigger Types Explained

| Trigger | When It Runs | Use Case |
|---------|--------------|----------|
| `pull_request` | PR opened, updated, or reopened | CI checks, code review automation |
| `push` | Commits pushed to branch | Post-merge tasks, deployments |
| `merge_group` | PR added to merge queue | Required for merge queue support |
| `workflow_dispatch` | Manual trigger via UI/API | Testing, manual runs |
| `schedule` | Cron-based schedule | Security scans, dependency updates |

### 1.2 The Merge Queue Trigger

**Critical:** If your repository uses merge queue, ALL required workflows MUST include the `merge_group` trigger.

```yaml
on:
  pull_request:
    branches: [main]
  merge_group:        # Required for merge queue!
  workflow_dispatch:
```

**Why:** When merge queue is enabled, GitHub creates a temporary merge commit and runs workflows using the `merge_group` event. Without this trigger, workflows won't run and the PR will be stuck in the queue indefinitely.

**Bootstrap Problem:** When first adding `merge_group` to workflows, you face a chicken-and-egg situation:
1. The PR adding `merge_group` gets stuck in the queue
2. Because `merge_group` trigger doesn't exist on the base branch yet

**Solution:** Temporarily disable merge queue, merge the PR, then re-enable merge queue.

### 1.3 Avoiding Duplicate Runs

**Problem:** Using both `push` and `pull_request` triggers on the same branch causes workflows to run twice for every PR commit.

**Anti-pattern (causes duplicates):**
```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

**Correct pattern for CI workflows:**
```yaml
on:
  pull_request:
    branches: [main]
  merge_group:
  workflow_dispatch:
```

**When to use `push` trigger:**
- Post-merge tasks (release automation, deployments)
- Workflows that should ONLY run after merge to main
- Scorecard and similar tools that evaluate repository state

### 1.4 Concurrency Settings

Always include concurrency settings to cancel outdated runs:

```yaml
concurrency:
  group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
  cancel-in-progress: true
```

This ensures that when you push new commits to a PR, the old workflow runs are cancelled.

### 1.5 Trigger Patterns by Workflow Type

| Workflow Type | Recommended Triggers |
|---------------|---------------------|
| CI (lint, test, build) | `pull_request`, `merge_group`, `workflow_dispatch` |
| Security (vulncheck, Grype) | `pull_request`, `merge_group`, `push`, `schedule`, `workflow_dispatch` |
| CodeQL | `pull_request`, `merge_group`, `push`, `schedule`, `workflow_dispatch` |
| Scorecard | `push` (main only), `schedule`, `workflow_dispatch` |
| Release | `push: tags: ['v*']` |

---

## Part 2: Rulesets & Branch Protection

### 2.1 The "All Jobs" Summary Pattern

**Problem:** Configuring rulesets to require individual job names is fragile:
- Job names may change
- Matrix jobs have dynamic names
- Hard to maintain as workflows evolve

**Solution:** Create a summary job that depends on all other jobs:

```yaml
jobs:
  lint:
    # ... lint job
  
  test:
    # ... test job
  
  all-jobs:
    name: CI / All Jobs          # This is the status check name
    runs-on: ubuntu-latest
    needs: [lint, test]          # Depends on all jobs
    if: always()                 # Runs even if dependencies fail
    steps:
      - name: Check all jobs passed
        run: |
          if [[ "${{ needs.lint.result }}" != "success" ]] || \
             [[ "${{ needs.test.result }}" != "success" ]]; then
            echo "One or more jobs failed"
            exit 1
          fi
          echo "All jobs passed"
```

**Benefits:**
- Single status check to configure in rulesets
- Automatically adapts when jobs are added/removed
- Clear naming convention: `<Workflow> / All Jobs`

### 2.2 Configuring Required Status Checks

In your repository's ruleset configuration:

1. Go to **Settings > Rules > Rulesets**
2. Create or edit a ruleset for your default branch
3. Enable **Require status checks to pass**
4. Add these status checks:
   - `CI / All Jobs`
   - `Security / All Jobs`
   - `CodeQL / All Jobs`

**Do NOT add:**
- Individual job names (e.g., `lint`, `test`)
- Scorecard checks (doesn't run on PRs)
- Any tool that only runs on schedule or push

### 2.3 Code Scanning Requirements

**How it works:** When you enable "Require code scanning results" in a ruleset, GitHub compares the PR's code scanning results against the base branch.

**The trap:** If a tool has uploaded SARIF to the base branch but doesn't run on PRs, GitHub will block the PR waiting for results that will never arrive.

**Rule:** Only require code scanning results from tools that:
1. Run on `pull_request` events
2. Upload SARIF to code-scanning

**Safe to require:**
- CodeQL (runs on PRs)
- Grype (runs on PRs)

**Do NOT require:**
- Scorecard (only runs on push/schedule)
- Any tool configured with only `push` or `schedule` triggers

### 2.4 Merge Queue Configuration

**Prerequisites:**
1. All required workflows have `merge_group` trigger
2. Ruleset is configured with status checks
3. Merge queue is enabled in repository settings

**Enable merge queue:**
1. Go to **Settings > General > Pull Requests**
2. Check **Require merge queue**
3. Configure queue settings (batch size, timeout, etc.)

**Troubleshooting:** If PRs get stuck in the queue showing "Waiting for status to be reported":
1. Verify all required workflows have `merge_group` trigger
2. Check that the workflows exist on the base branch (not just in the PR)
3. See [Troubleshooting](#52-merge-queue-stuck) section

---

## Part 3: Code Scanning

### 3.1 CodeQL Configuration

**Recommendation:** Use a custom workflow instead of GitHub's "Default Setup" for better control.

**Why avoid Default Setup:**
- Scans languages you may not need (JavaScript in Go projects)
- Creates analyses that are hard to clean up
- Less control over when scans run

**Single-language CodeQL workflow:**
```yaml
- name: Initialize CodeQL
  uses: github/codeql-action/init@v4
  with:
    languages: go              # Only scan Go
    build-mode: autobuild

- name: Perform CodeQL Analysis
  uses: github/codeql-action/analyze@v4
  with:
    category: "/language:go"   # Consistent category naming
```

### 3.2 SARIF Upload Rules

**When to upload SARIF to code-scanning:**
- Tool runs on `pull_request` events
- You want results to appear in GitHub Security tab
- You want PR comparisons (new vs existing alerts)

**When NOT to upload SARIF:**
- Tool only runs on `push` or `schedule`
- Tool evaluates repository state, not code (e.g., Scorecard)

**Scorecard example (correct - no SARIF upload):**
```yaml
- name: Run analysis
  uses: ossf/scorecard-action@v2
  with:
    results_file: results.sarif
    publish_results: true       # Publishes to OSSF dashboard

- name: Upload artifact         # Upload as artifact for manual review
  uses: actions/upload-artifact@v4
  with:
    name: scorecard-results
    path: results.sarif

# Do NOT add upload-sarif step!
```

### 3.3 Managing Stale Analyses

**Problem:** When you change code scanning configuration (disable tools, change languages), stale analyses may remain and block PRs.

**Symptoms:**
- "Code scanning is waiting for results from X"
- PR blocked even though X no longer runs

**Solution:** Delete stale analyses using `gh api`. See [Appendix](#appendix-gh-api-commands) for commands.

---

## Part 4: Workflow Templates

> **Note:** The templates below use `main` as the default branch name. If your
> repository uses a different default branch (e.g., `master`), replace `main`
> with your branch name in the `branches:` arrays.

### 4.1 CI Workflow Template

```yaml
# .github/workflows/ci.yml
#
# Continuous Integration workflow
# Runs: lint, test, and other CI checks on every PR

name: ci

on:
  pull_request:
    branches: [main]
  merge_group:
  workflow_dispatch:

permissions:
  contents: read

concurrency:
  group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
  cancel-in-progress: true

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest

  test:
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        go-version: ['stable', 'oldstable']
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Run tests
        run: go test -race -coverprofile=coverage.out ./...

      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{gsub(/%/, ""); print int($3)}')
          THRESHOLD=65
          if [ "$COVERAGE" -lt "$THRESHOLD" ]; then
            echo "::error::Coverage ${COVERAGE}% is below threshold ${THRESHOLD}%"
            exit 1
          fi

  # Summary job for ruleset configuration
  all-jobs:
    name: CI / All Jobs
    runs-on: ubuntu-latest
    needs: [lint, test]
    if: always()
    steps:
      - name: Check all jobs passed
        run: |
          if [[ "${{ needs.lint.result }}" != "success" ]] || \
             [[ "${{ needs.test.result }}" != "success" ]]; then
            echo "One or more jobs failed"
            exit 1
          fi
          echo "All jobs passed"
```

### 4.2 Security Workflow Template

```yaml
# .github/workflows/security.yml
#
# Security scanning workflow
# Runs: vulnerability checks on PRs and periodically

name: security

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  merge_group:
  schedule:
    - cron: '0 6 * * 1'  # Weekly on Mondays at 6am UTC
  workflow_dispatch:

permissions:
  contents: read
  security-events: write  # Required for SARIF upload

jobs:
  govulncheck:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'

      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest

      - name: Run govulncheck
        run: govulncheck ./...

  grype:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run Grype vulnerability scanner
        uses: anchore/scan-action@v4
        id: scan
        with:
          path: "."
          fail-build: false
          output-format: sarif

      - name: Upload SARIF to GitHub Security
        uses: github/codeql-action/upload-sarif@v4
        with:
          sarif_file: ${{ steps.scan.outputs.sarif }}
          category: "grype"

  # Summary job for ruleset configuration
  all-jobs:
    name: Security / All Jobs
    runs-on: ubuntu-latest
    needs: [govulncheck, grype]
    if: always()
    steps:
      - name: Check all jobs passed
        run: |
          if [[ "${{ needs.govulncheck.result }}" != "success" ]] || \
             [[ "${{ needs.grype.result }}" != "success" ]]; then
            echo "One or more jobs failed"
            exit 1
          fi
          echo "All jobs passed"
```

### 4.3 CodeQL Workflow Template

```yaml
# .github/workflows/codeql.yml
#
# CodeQL static analysis workflow
# Runs: on PRs and periodically for security analysis

name: codeql

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  merge_group:
  schedule:
    - cron: '0 6 * * 1'  # Weekly on Mondays at 6am UTC
  workflow_dispatch:

permissions:
  contents: read
  security-events: write

jobs:
  analyze:
    name: Analyze
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Initialize CodeQL
        uses: github/codeql-action/init@v4
        with:
          languages: go           # Specify your language(s)
          build-mode: autobuild

      - name: Perform CodeQL Analysis
        uses: github/codeql-action/analyze@v4
        with:
          category: "/language:go"

  # Summary job for ruleset configuration
  all-jobs:
    name: CodeQL / All Jobs
    runs-on: ubuntu-latest
    needs: [analyze]
    if: always()
    steps:
      - name: Check all jobs passed
        run: |
          if [[ "${{ needs.analyze.result }}" != "success" ]]; then
            echo "CodeQL analysis failed"
            exit 1
          fi
          echo "All jobs passed"
```

### 4.4 Scorecard Workflow Template

```yaml
# .github/workflows/scorecard.yml
#
# OSSF Scorecard supply chain security analysis
# Runs: ONLY on push to main and schedule (NOT on PRs)
#
# IMPORTANT: This workflow does NOT upload SARIF to code-scanning
# because it doesn't run on PRs. Uploading SARIF would cause PRs
# to be blocked waiting for Scorecard results.

name: scorecard

on:
  push:
    branches: [main]
  schedule:
    - cron: '0 6 * * 1'  # Weekly on Mondays at 6am UTC
  workflow_dispatch:

# Note: No pull_request or merge_group triggers!

permissions: read-all

jobs:
  analysis:
    name: Scorecard analysis
    runs-on: ubuntu-latest
    permissions:
      security-events: write
      id-token: write

    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          persist-credentials: false

      - name: Run analysis
        uses: ossf/scorecard-action@v2
        with:
          results_file: results.sarif
          results_format: sarif
          publish_results: true   # Publishes to OSSF Scorecard dashboard

      # Upload as artifact for manual review
      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: scorecard-results
          path: results.sarif
          retention-days: 5

      # Do NOT upload to code-scanning!
      # This would cause PRs to be blocked waiting for Scorecard results.

  # Summary job (optional, but keeps consistency)
  all-jobs:
    name: Scorecard / All Jobs
    runs-on: ubuntu-latest
    needs: [analysis]
    if: always()
    steps:
      - name: Check all jobs passed
        run: |
          if [[ "${{ needs.analysis.result }}" != "success" ]]; then
            echo "Scorecard analysis failed"
            exit 1
          fi
          echo "All jobs passed"
```

### 4.5 Release Workflow Template

```yaml
# .github/workflows/release.yml
#
# Release automation workflow
# Runs: when version tags are pushed (v*)

name: release

on:
  push:
    tags:
      - "v*"

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Part 5: Troubleshooting

### 5.1 PR Blocked by Code Scanning

**Symptom:** PR shows "Code scanning is waiting for results from X for commits..."

**Causes:**
1. Tool X only runs on push/schedule, not pull_request
2. Tool X was previously uploading SARIF but has been removed/disabled
3. Stale analyses exist from a tool that no longer runs

**Diagnosis:**
```bash
# List all code scanning analyses on the base branch
gh api repos/{owner}/{repo}/code-scanning/analyses --paginate | \
  jq -r '.[] | select(.ref == "refs/heads/main") | "\(.tool.name) \(.category)"' | \
  sort | uniq
```

**Solutions:**

1. **If tool shouldn't be required:** Remove it from ruleset's code scanning requirements

2. **If tool should run on PRs:** Add `pull_request` and `merge_group` triggers to the workflow

3. **If stale analyses exist:** Delete them (see [Appendix](#delete-analyses-for-a-specific-toolcategory))

### 5.2 Merge Queue Stuck

**Symptom:** PR is in merge queue, showing "Waiting for status to be reported"

**Causes:**
1. Workflows missing `merge_group` trigger
2. Required workflows don't exist on base branch
3. GitHub hasn't picked up the workflow triggers yet

**Diagnosis:**
```bash
# Check if workflows are running for merge queue
gh run list --repo {owner}/{repo} --limit 10

# Look for runs with "merge_group" event type
```

**Solutions:**

1. **Missing `merge_group` trigger:**
   - Add `merge_group:` to all required workflows
   - Bootstrap problem: temporarily disable merge queue to merge the fix

2. **Workflows not on base branch:**
   - Ensure workflows exist on main/master branch
   - Merge queue runs workflows from the merge commit, but triggers must exist on base

3. **GitHub delay:**
   - Wait a minute - sometimes there's a delay before workflows start
   - Check GitHub Actions status page for outages

### 5.3 Duplicate Workflow Runs

**Symptom:** Every PR commit triggers the same workflow twice

**Cause:** Using both `push` and `pull_request` triggers on the same branch

**Diagnosis:**
```yaml
# Check your workflow triggers
on:
  push:
    branches: [main]      # This triggers on PR merges AND PR commits
  pull_request:
    branches: [main]      # This also triggers on PR commits
```

**Solution:**
```yaml
# For CI workflows, use only pull_request:
on:
  pull_request:
    branches: [main]
  merge_group:
  workflow_dispatch:

# For post-merge workflows, use only push:
on:
  push:
    branches: [main]
  workflow_dispatch:
```

---

## Appendix: gh api Commands

> **Note:** In the commands below, replace these placeholders with actual values:
> - `{owner}` - Repository owner (e.g., `actualyze-ai`)
> - `{repo}` - Repository name (e.g., `mage`)
> - `{analysis_id}` - Numeric analysis ID from the list commands

### List All Code Scanning Analyses

```bash
# List all analyses for a repository
gh api repos/{owner}/{repo}/code-scanning/analyses --paginate | \
  jq -r '.[] | "\(.id) \(.ref) \(.tool.name) \(.category) \(.created_at)"'
```

### List Analyses on Default Branch

```bash
# Filter to only default branch (main or master)
gh api repos/{owner}/{repo}/code-scanning/analyses --paginate | \
  jq -r '.[] | select(.ref == "refs/heads/main") | "\(.id) \(.tool.name) \(.category)"'
```

### Check if Analysis is Deletable

```bash
# Get details for a specific analysis
gh api repos/{owner}/{repo}/code-scanning/analyses/{analysis_id} | \
  jq '{id, category, deletable, created_at}'
```

### Delete a Single Analysis

```bash
# Delete an analysis (must be deletable: true)
gh api -X DELETE "repos/{owner}/{repo}/code-scanning/analyses/{analysis_id}?confirm_delete=true"
```

**Note:** The `confirm_delete=true` parameter is required when deleting the last analysis for a tool/category.

### Delete Analyses for a Specific Tool/Category

```bash
# Delete all analyses for a specific category (e.g., Scorecard)
# Must delete newest first - only the most recent is deletable at any time

OWNER="your-org"
REPO="your-repo"
CATEGORY="supply-chain/online-scm"  # or "/language:javascript-typescript", etc.

while true; do
  ANALYSIS_ID=$(gh api "repos/${OWNER}/${REPO}/code-scanning/analyses" --paginate | \
    jq -r "[.[] | select(.ref == \"refs/heads/main\" and .category == \"${CATEGORY}\" and .deletable == true)] | .[0].id // empty")
  
  if [ -z "$ANALYSIS_ID" ]; then
    echo "No more deletable analyses for ${CATEGORY}"
    break
  fi
  
  echo "Deleting analysis ${ANALYSIS_ID}..."
  gh api -X DELETE "repos/${OWNER}/${REPO}/code-scanning/analyses/${ANALYSIS_ID}?confirm_delete=true"
done
```

### Full Cleanup Script

```bash
#!/bin/bash
#
# cleanup-code-scanning.sh
#
# Deletes all code scanning analyses for specified tools/categories
# Usage: ./cleanup-code-scanning.sh owner repo category1 [category2 ...]
#
# Example: ./cleanup-code-scanning.sh actualyze-ai mage \
#   "/language:javascript-typescript" \
#   "/language:actions" \
#   "supply-chain/online-scm" \
#   "supply-chain/local" \
#   "supply-chain/branch-protection"

set -e

OWNER="${1:?Usage: $0 owner repo category1 [category2 ...]}"
REPO="${2:?Usage: $0 owner repo category1 [category2 ...]}"
shift 2

for CATEGORY in "$@"; do
  echo "=== Deleting analyses for category: ${CATEGORY} ==="
  
  while true; do
    ANALYSIS_ID=$(gh api "repos/${OWNER}/${REPO}/code-scanning/analyses" --paginate 2>/dev/null | \
      jq -r "[.[] | select(.ref == \"refs/heads/main\" and .category == \"${CATEGORY}\" and .deletable == true)] | .[0].id // empty")
    
    if [ -z "$ANALYSIS_ID" ]; then
      echo "No more deletable analyses for ${CATEGORY}"
      break
    fi
    
    echo "  Deleting analysis ${ANALYSIS_ID}..."
    gh api -X DELETE "repos/${OWNER}/${REPO}/code-scanning/analyses/${ANALYSIS_ID}?confirm_delete=true" 2>/dev/null | \
      jq -r '.next_analysis_url // "deleted"'
  done
done

echo ""
echo "=== Remaining analyses on main branch ==="
gh api "repos/${OWNER}/${REPO}/code-scanning/analyses" --paginate 2>/dev/null | \
  jq -r '.[] | select(.ref == "refs/heads/main") | "\(.tool.name) \(.category)"' | \
  sort | uniq
```

### Required Token Permissions

To delete code scanning analyses, your GitHub token needs:
- `repo` scope (full control of private repositories), OR
- `security_events` scope (read/write security events)

Check your token permissions:
```bash
gh auth status
```

---

## Summary

### Key Principles

1. **Always include `merge_group` trigger** when using merge queue
2. **Use `pull_request` for CI**, not `push`
3. **Create "All Jobs" summary jobs** for easier ruleset configuration
4. **Only upload SARIF for tools that run on PRs**
5. **Delete stale analyses** when changing code scanning configuration
6. **Don't require Scorecard** in code scanning (it doesn't run on PRs)

### Workflow Trigger Quick Reference

| Workflow Type | `pull_request` | `merge_group` | `push` | `schedule` |
|---------------|----------------|---------------|--------|------------|
| CI            | Yes            | Yes           | No     | No         |
| Security      | Yes            | Yes           | Yes    | Yes        |
| CodeQL        | Yes            | Yes           | Yes    | Yes        |
| Scorecard     | No             | No            | Yes    | Yes        |
| Release       | No             | No            | Tags   | No         |

### Status Check Names

Configure these in your ruleset:
- `CI / All Jobs`
- `Security / All Jobs`
- `CodeQL / All Jobs`

Do NOT configure:
- `Scorecard / All Jobs` (doesn't run on PRs)
- Individual job names
