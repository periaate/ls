package lfs

import (
	"archive/zip"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/periaate/ls/files"
)

type FSWorker struct {
	Sort     SortBy
	Hide     bool
	Archives bool
	// Format directory paths to end with "/". Used for internal logic, turning it
	// off will remove file|directory selection functionality
	WebStyle bool
	Logger
}

func NewFSWorker() *FSWorker {
	return &FSWorker{
		Sort:     ByNone,
		Hide:     true,
		Archives: false,
		WebStyle: true,
		Logger:   DummyLogger{},
	}
}

func (fsw *FSWorker) Parser() func(string) ([]*Element, error) {
	return func(path string) (res []*Element, err error) {
		path, err = filepath.Abs(path)
		if err != nil {
			return
		}
		fsw.Debug("parsing", "path", path)
		if fsw.Hide {
			base := filepath.Base(path)
			if len(base) > 1 {
				if base[0] == '.' {
					return
				}
			}
		}

		var finfos []fs.FileInfo
		switch {
		case fsw.Archives && IsZipLike(path):
			finfos, err = fsw.Zip(path)
			if err != nil {
				return
			}
		default:
			finfos, err = fsw.Dir(path)
			if err != nil {
				fsw.Error("encountered error reading a directory", "path", path, "err", err)
				return
			}

			fsw.Debug("reading a directory", "path", path, "files", len(finfos))
		}

		for _, fi := range finfos {
			p := filepath.Join(path, fi.Name())
			if fsw.Hide && files.ShouldIgnore(p) {
				continue
			}

			var isDir, isArchive bool
			isDir = fi.IsDir()
			if fsw.Archives && IsZipLike(p) {
				// if archives are included, they are considered to be directories
				isArchive = true
				isDir = true
			}

			if isDir && fsw.WebStyle {
				p += "/"
			}

			el := fsw.Parse(p)

			if isArchive {
				el.Mask |= files.MaskDirectory
				el.Mask &= ^files.MaskFile
			}

			switch fsw.Sort {
			case ByMod:
				addModT(el, fi)
			case BySize:
				addSize(el, fi)
			case ByCreation:
				addCreationT(el, fi)
			}

			res = append(res, el)
		}
		return
	}
}

func (fsw *FSWorker) Dir(path string) (files []fs.FileInfo, err error) {
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

func (fsw *FSWorker) Zip(path string) (files []fs.FileInfo, err error) {
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

func (fsw *FSWorker) Parse(path string) *Element {
	isDir := path[len(path)-1] == '/'
	path = filepath.ToSlash(path)
	name := filepath.Base(path)
	fi := &Element{
		Name: name,
		Path: path,
	}

	fi.Mask |= files.ExtToMaskMap[filepath.Ext(fi.Name)]
	if isDir {
		fi.Mask |= files.MaskDirectory
	}

	return fi
}
