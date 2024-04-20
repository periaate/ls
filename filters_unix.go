//go:build !windows
// +build !windows

package ls

import (
	"io/fs"
)

// Implementation or stub of addCreationT for Unix
func addCreationT(_ *Element, _ fs.FileInfo) {}
