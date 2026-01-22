<!--
SPDX-License-Identifier: Apache-2.0
Modifications Copyright (c) 2026 Actualyze AI

NOTE: This file has been modified by Actualyze AI from the original upstream
version (magefile/mage). See git history for details.
-->

[![Built with Mage](https://magefile.org/badge.svg)](https://magefile.org)
[![CI](https://github.com/actualyze-ai/mage/actions/workflows/ci.yml/badge.svg)](https://github.com/actualyze-ai/mage/actions/workflows/ci.yml)

<p align="center"><img src="https://user-images.githubusercontent.com/3185864/32058716-5ee9b512-ba38-11e7-978a-287eb2a62743.png"/></p>

## Fork Notice

This is a fork of [magefile/mage](https://github.com/magefile/mage), maintained by
[Actualyze AI](https://github.com/actualyze-ai). This fork includes modernization
updates, CI improvements, and is used as part of the Actualyze AI infrastructure.

For issues specific to this fork, please open an issue at
[actualyze-ai/mage](https://github.com/actualyze-ai/mage/issues).

For upstream mage issues and general mage discussion, see the
[upstream repository](https://github.com/magefile/mage).

## About

Mage is a make-like build tool using Go. You write plain-old go functions,
and Mage automatically uses them as Makefile-like runnable targets.

## Installation

Mage has no dependencies outside the Go standard library and requires Go 1.25
or above.

**Using Go Install (Recommended)**

```
go install github.com/actualyze-ai/mage@latest
mage -init
```

**Using Go Modules**

```
git clone https://github.com/actualyze-ai/mage
cd mage
go run bootstrap.go
```

This will download the code and then run the bootstrap script to build mage with
version information embedded in it. A normal `go install` will build the binary
correctly, but no version info will be embedded. If you've done this, no worries,
just go to the cloned directory and run `mage install` or `go run bootstrap.go`
and a new binary will be created with the correct version information.

You may also install a binary release from our
[releases](https://github.com/actualyze-ai/mage/releases) page.

## Demo

[![Mage Demo](https://img.youtube.com/vi/GOqbD0lF-iA/maxresdefault.jpg)](https://www.youtube.com/watch?v=GOqbD0lF-iA)

## Discussion

For upstream mage community discussion, join the `#mage` channel on
[gophers slack](https://gophers.slack.com/messages/general/) or post on the
[magefile google group](https://groups.google.com/forum/#!forum/magefile).

# Documentation

See [magefile.org](https://magefile.org) for full documentation.

See [pkg.go.dev/github.com/actualyze-ai/mage/mage](https://pkg.go.dev/github.com/actualyze-ai/mage/mage)
for instructions on how to use Mage as a library.

# Why?

Makefiles are hard to read and hard to write. Mostly because makefiles are
essentially fancy bash scripts with significant white space and additional
make-related syntax.

Mage lets you have multiple magefiles, name your magefiles whatever you want,
and they're easy to customize for multiple operating systems. Mage has no
dependencies (aside from go) and runs just fine on all major operating systems,
whereas make generally uses bash which is not well supported on Windows. Go is
superior to bash for any non-trivial task involving branching, looping, anything
that's not just straight line execution of commands. And if your project is
written in Go, why introduce another language as idiosyncratic as bash? Why not
use the language your contributors are already comfortable with?
