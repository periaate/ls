package ls

import (
	"log/slog"
	"math"
	"sort"

	"github.com/facette/natsort"
	"github.com/periaate/common"
	"github.com/periaate/ls/lfs"
	"github.com/periaate/slice"
)

func NoneFilter(*lfs.Element) bool                  { return true }
func NoneProcess(els []*lfs.Element) []*lfs.Element { return els }

type Process func(inp []*lfs.Element) (out []*lfs.Element)
type Filter func(*lfs.Element) bool

type Option func(*lfs.FSTraverser) *lfs.FSTraverser

func FromDepth(depth int) Option {
	return func(opts *lfs.FSTraverser) *lfs.FSTraverser {
		opts.MinDepth = depth
		return opts
	}
}

func ToDepth(depth int) Option {
	return func(opts *lfs.FSTraverser) *lfs.FSTraverser {
		opts.MaxDepth = depth
		return opts
	}
}

func DepthPattern(pat string) Option {
	return func(opts *lfs.FSTraverser) *lfs.FSTraverser {
		from, to, err := slice.ParsePattern(pat, math.MaxInt)
		if err != nil {
			return opts
		}
		opts.MinDepth = from
		opts.MaxDepth = to
		return opts
	}
}

func Search(queries ...string) Option {
	return func(opts *lfs.FSTraverser) *lfs.FSTraverser {
		f := ParseSearch(queries)
		if f != nil {
			opts.ResFilter = common.All(true, opts.ResFilter, f)
		}
		return opts
	}
}

func LogWith(logger *slog.Logger) Option {
	return func(opts *lfs.FSTraverser) *lfs.FSTraverser {
		opts.Logger = logger
		return opts
	}
}

const (
	Include = true
	Exclude = false
)

func Recurse(opts *lfs.FSTraverser) *lfs.FSTraverser {
	opts.MaxDepth = math.MaxInt64
	return opts
}

func NoHide(opts *lfs.FSTraverser) *lfs.FSTraverser {
	opts.FSW.Hide = false
	return opts
}

func Dir(opts ...Option) (tr *lfs.FSTraverser) {
	for _, opt := range opts {
		tr = opt(tr)
		if tr == nil {
			return nil
		}
	}
	return
}

func Paths(paths ...string) Option {
	return func(opts *lfs.FSTraverser) *lfs.FSTraverser {
		if opts == nil {
			opts = lfs.GetDefault()
		}
		opts.Src.Seed(paths)
		return opts
	}
}

func Masks(inc bool, masks ...uint32) Option {
	return func(opts *lfs.FSTraverser) *lfs.FSTraverser {
		var mask uint32
		for _, m := range masks {
			mask |= m
		}
		f := MaskFilter(mask)
		if f != nil {
			opts.ResFilter = common.All(true, opts.ResFilter, f)
		}
		return opts
	}
}

func Combine(trs ...*lfs.FSTraverser) Process {
	return func(inp []*lfs.Element) (out []*lfs.Element) {
		for _, tr := range trs {
			res := tr.Traverse(tr.Src)
			inp = append(inp, res...)
		}
		return inp
	}
}

func Do(prs ...Process) (res []*lfs.Element) {
	for _, proc := range prs {
		res = proc(res)
	}
	return res
}

func Sort(by lfs.SortBy) Process {
	return func(inp []*lfs.Element) (out []*lfs.Element) {
		switch by {
		case lfs.ByName:
			sort.Slice(inp, func(i, j int) bool {
				return !natsort.Compare(inp[i].Name, inp[j].Name)
			})
		case lfs.ByMod:
			sort.Slice(inp, func(i, j int) bool {
				return inp[i].Mod > inp[j].Mod
			})
		case lfs.BySize:
			sort.Slice(inp, func(i, j int) bool {
				return inp[i].Size > inp[j].Size
			})
		case lfs.ByCreation:
			sort.Slice(inp, func(i, j int) bool {
				return inp[i].Creation > inp[j].Creation
			})
		}
		return inp
	}
}

func TimeSlice(pattern string) Process {
	return func(inp []*lfs.Element) (out []*lfs.Element) {
		from, to, err := slice.ParseTimeSlice(pattern)
		if err != nil {
			return inp
		}
		for _, el := range inp {
			if from <= el.Mod && el.Mod <= to {
				out = append(out, el)
			}
		}
		return out
	}
}

/*
ls.Do(
	ls.Combine(
		ls.Dir(
			ls.Paths("./"),
			ls.Masks(ls.Include, files.Image, files.Video),
		),
		ls.Dir(
			ls.Paths("./"),
			ls.Masks(ls.Include, files.Image, files.Video),
			ls.Slice(0, 200),
		),
	),
	ls.Sort(ls.TimeCreated),
	ls.Slice("[:200]"),
	ls.FuzzySearch("_p01"),
)
*/
