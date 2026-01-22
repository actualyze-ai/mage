//go:build mage
// +build mage

// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright (c) 2026 Actualyze AI
//
// NOTE: This file has been modified by Actualyze AI from the original upstream
// version (magefile/mage). See git history for details.

// Compiled package description.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/actualyze-ai/mage/mg"
)

var Default = Deploy

// This is very verbose.
func TestVerbose() {
	log.Println("hi!")
}

// PrintVerboseFlag prints the value of mg.Verbose() to stdout.
func PrintVerboseFlag() {
	fmt.Printf("mg.Verbose()==%v", mg.Verbose())
}

// This is the synopsis for Deploy. This part shouldn't show up.
func Deploy() {
	mg.Deps(f)
}

// Sleep sleeps 5 seconds.
func Sleep() {
	time.Sleep(5 * time.Second)
}

func f() {
	log.Println("i am independent -- not")
}
