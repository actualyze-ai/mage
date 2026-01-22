// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright (c) 2026 Actualyze AI
//
// NOTE: This file has been modified by Actualyze AI from the original upstream
// version (magefile/mage). See git history for details.

package deps

import (
	"fmt"

	"github.com/actualyze-ai/mage/mg"
)

// All code in this package belongs to @na4ma4 in GitHub https://github.com/na4ma4/magefile-test-import
// reproduced here for ease of testing regression on bug 508

type Docker mg.Namespace

func (Docker) Test() {
	fmt.Println("docker")
}

func Test() {
	fmt.Println("test")
}
