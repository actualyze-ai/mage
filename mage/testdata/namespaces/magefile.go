//go:build mage
// +build mage

// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright (c) 2026 Actualyze AI
//
// NOTE: This file has been modified by Actualyze AI from the original upstream
// version (magefile/mage). See git history for details.

package main

import (
	"context"
	"fmt"

	"github.com/actualyze-ai/mage/mg"
)

var Default = NS.Error

func TestNamespaceDep() {
	mg.Deps(NS.Error, NS.Bare, NS.BareCtx, NS.CtxErr)
}

type NS mg.Namespace

func (NS) Error() error {
	fmt.Println("hi!")
	return nil
}

func (NS) Bare() {
}

func (NS) BareCtx(ctx context.Context) {
}
func (NS) CtxErr(ctx context.Context) error {
	return nil
}
