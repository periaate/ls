//go:build windows
// +build windows

package lfs

import (
	"io/fs"
	"syscall"
)

func addCreationT(fi *Element, info fs.FileInfo) {
	winFileInfo := info.Sys().(*syscall.Win32FileAttributeData)

	fi.Vany = winFileInfo.CreationTime.Nanoseconds()
}
