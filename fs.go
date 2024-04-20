package ls

import (
	"archive/zip"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

func FSParser(opts *Options) func(string) ([]*Element, error) {
	home, _ := os.UserHomeDir()

	return func(path string) (res []*Element, err error) {
		path = ResolveHome(home, path)
		if opts.Hide {
			base := filepath.Base(path)
			if len(base) > 1 {
				if base[0] == '.' {
					return
				}
			}
		}

		var finfos []fs.FileInfo
		switch {
		case opts.Archives && IsZipLike(path):
			finfos, err = TraverseZip(path)
			if err != nil {
				return
			}
		default:
			finfos, err = TraverseDir(path)
			if err != nil {
				return
			}
		}

		for _, fi := range finfos {
			p := filepath.Join(path, fi.Name())
			if opts.Hide && ShouldIgnore(p) {
				continue
			}

			var isDir, isArchive bool
			isDir = fi.IsDir()
			if opts.Archives && IsZipLike(p) {
				// if archives are included, they are considered to be directories
				isArchive = true
				isDir = true
			}

			if isDir && opts.WebStyle {
				p += "/"
			}

			el := FileParser(p)

			if isArchive {
				el.Mask |= MaskDirectory
				el.Mask &= ^MaskFile
			}

			switch opts.Sort {
			case ByMod:
				addModT(el, fi)
			case BySize:
				addSize(el, fi)
			case ByCreation:
				addCreationT(el, fi)
			}
		}
		return
	}
}

func TraverseDir(path string) (files []fs.FileInfo, err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			slog.Debug("error reading file info", "file", entry.Name(), "error", err)
			continue
		}
		files = append(files, info)
	}
	return
}

func TraverseZip(path string) (files []fs.FileInfo, err error) {
	// TODO: tar, 7z, xz, etc., support
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		info := f.FileInfo()
		if info.IsDir() {
			continue
		}

		files = append(files, info)
	}

	return
}

func FileParser(path string) *Element {
	isDir := path[len(path)-1] == '/'
	path = filepath.ToSlash(path)
	name := filepath.Base(path)
	fi := &Element{
		Name: name,
		Path: path,
	}

	fi.Mask |= ExtToMaskMap[filepath.Ext(fi.Name)]
	if isDir {
		fi.Mask |= MaskDirectory
	}

	return fi
}
