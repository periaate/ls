package lfs

import (
	"io/fs"
	"path/filepath"

	"github.com/periaate/ls/files"
)

func IsZipLike(path string) bool {
	return files.ExtToMaskMap[filepath.Ext(path)]&files.MaskZipLike != 0
}

func addModT(fi *Element, info fs.FileInfo) { fi.Mod = info.ModTime().Unix() }
func addSize(fi *Element, info fs.FileInfo) { fi.Size = info.Size() }
