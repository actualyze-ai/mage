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

# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.16.x  | :white_check_mark: |
| < 1.16  | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue,
please report it responsibly.

### Preferred Method: GitHub Security Advisories

1. Go to the [Security tab](https://github.com/actualyze-ai/mage/security/advisories)
2. Click "Report a vulnerability"
3. Provide details about the vulnerability

This allows for private disclosure and coordinated fixes.

### Alternative: Email

For sensitive issues, you can also email: security@actualyze.ai

### What to Include

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Any suggested fixes (optional)

### Response Timeline

- **Acknowledgment:** Within 48-72 hours
- **Initial Assessment:** Within 1 week
- **Fix Timeline:** Depends on severity
  - Critical: As soon as possible
  - High: Within 2 weeks
  - Medium/Low: Next scheduled release
- **Public Disclosure:** 90 days after report, or when fix is released (whichever is first)

### What to Expect

1. We'll acknowledge your report promptly
2. We'll investigate and keep you informed of progress
3. We'll credit you in the security advisory (unless you prefer anonymity)
4. We'll coordinate disclosure timing with you

## Security Best Practices

When using Mage in your projects:

- Keep Mage updated to the latest supported version
- Review magefiles before running them, especially from untrusted sources
- Use `go mod verify` to ensure dependency integrity
- Consider running mage in a sandboxed environment for untrusted builds
