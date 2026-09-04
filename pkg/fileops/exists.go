// Copyright 2026 OVGU Magdeburg
// SPDX-License-Identifier: Apache-2.0

package fileops

import "os"

func FileOrFolderExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
