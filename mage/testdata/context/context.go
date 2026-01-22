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
	"errors"
	"fmt"
	"time"

	"github.com/actualyze-ai/mage/mg"
)

// Returns a non-nil error.
func TakesContextNoError(ctx context.Context) {
	deadline, _ := ctx.Deadline()
	fmt.Printf("Context timeout: %v\n", deadline)
}

func Timeout(ctx context.Context) {
	time.Sleep(200 * time.Millisecond)
}

func TakesContextWithError(ctx context.Context) error {
	return errors.New("Something went sideways")
}

func CtxDeps(ctx context.Context) {
	mg.CtxDeps(ctx, TakesContextNoError)
}
