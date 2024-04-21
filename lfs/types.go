package lfs

type SortBy uint8

const (
	ByNone SortBy = iota
	ByName
	ByMod
	BySize
	ByCreation // Available only on Windows
)

type Element struct {
	Name     string
	Path     string // includes name, relative path to cwd
	Mod      int64  // Unix time for mod
	Creation int64  // Unix time for creation
	Size     int64  // Size in bytes
	Mask     uint32 // Bitmask of element types, file, dir, content types, hidden, etc.
}

type SourceIter[T, S any] interface {
	Iter() (res []T, curr S, ok bool)
	Seed(seeds []S)
}

type Traverser[T, S any] interface {
	Traverse(SourceIter[T, S]) []T
}
