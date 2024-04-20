package ls

import (
	"github.com/periaate/ls/lfs"
)

// type Options struct {
// 	ToDepth   int
// 	FromDepth int
// 	Queries   []Query

// 	Args []string

// 	Filters   []func(*Element) bool
// 	Processes []func(els []*Element) []*Element

// 	Select []string
// }

func NoneFilter(*lfs.Element) bool                  { return true }
func NoneProcess(els []*lfs.Element) []*lfs.Element { return els }

type Process func(inp []*lfs.Element) (out []*lfs.Element)
type Filter func(*lfs.Element) bool

type Result struct{ Files []*lfs.Element }
