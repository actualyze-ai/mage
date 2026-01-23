// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright (c) 2026 Actualyze AI
//
// NOTE: This file has been modified by Actualyze AI from the original upstream
// version (magefile/mage). See git history for details.

package mg

import (
	"errors"
	"fmt"
)

type fatalErr struct {
	error

	code int
}

func (f fatalErr) ExitStatus() int {
	return f.code
}

type exitStatus interface {
	ExitStatus() int
}

// Fatal returns an error that will cause mage to print out the
// given args and exit with the given exit code.
func Fatal(code int, args ...interface{}) error {
	return fatalErr{
		code:  code,
		error: errors.New(fmt.Sprint(args...)),
	}
}

// Fatalf returns an error that will cause mage to print out the
// given message and exit with the given exit code.
func Fatalf(code int, format string, args ...interface{}) error {
	return fatalErr{
		code:  code,
		error: fmt.Errorf(format, args...),
	}
}

// ExitStatus queries the error for an exit status.  If the error is nil, it
// returns 0.  If the error does not implement ExitStatus() int, it returns 1.
// Otherwise it retiurns the value from ExitStatus().
func ExitStatus(err error) int {
	if err == nil {
		return 0
	}
	exit, ok := err.(exitStatus)
	if !ok {
		return 1
	}
	return exit.ExitStatus()
}
