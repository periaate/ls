package lfs

import (
	"io/fs"
	"path/filepath"

	"github.com/periaate/ls/files"
)

func IsZipLike(path string) bool {
	return files.ExtToMaskMap[filepath.Ext(path)]&files.MaskZipLike != 0
}

func addModT(fi *Element, info fs.FileInfo) { fi.Vany = info.ModTime().Unix() }
func addSize(fi *Element, info fs.FileInfo) { fi.Vany = info.Size() }

func ResolveHome(home, path string) string {
	if len(path) == 0 || len(home) == 0 {
		return path
	}
	return filepath.Join(home, path[1:])
}
