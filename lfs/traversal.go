package lfs

import (
	"github.com/periaate/common"
	"github.com/periaate/ls/files"
)

func GetDefault() *FSTraverser {
	fst := NewFSTraverser()
	fst.FSW = NewFSWorker()
	fst.Src = NewFSSource()
	fst.Src.Parser = fst.FSW.Parser()
	return fst
}

func NewFSTraverser() *FSTraverser {
	fst := &FSTraverser{}
	fst.Init()
	return fst
}

type FSTraverser struct {
	MaxDepth int
	MinDepth int

	// Selection filter, determines which paths are and aren't traversed
	SelFilter func(*Element) bool
	// Result filter, determines which elments are and aren't included in result
	ResFilter func(*Element) bool

	FSW *FSWorker
	Src *FSSource

	common.Logger
}

func BaseFilter(el *Element) bool {
	return !files.ShouldIgnore(el.Name)
}

func (t *FSTraverser) Init() {
	t.MaxDepth = 0
	t.MinDepth = -1
	t.SelFilter = BaseFilter
	t.ResFilter = BaseFilter
	t.Logger = common.DummyLogger{}
}

func (tr *FSTraverser) Traverse(src SourceIter[*Element, string]) (res []*Element) {
	var depth int
	defer func() {
		tr.Info("traversal finished", "depth", depth, "total found", len(res))
	}()

	var seeds []string
	temp, curr, ok := src.Iter()
	for ok {
		tr.Info("traversing", "PATH", curr, "DEPTH", depth, "ELS", len(temp))
		if depth >= tr.MaxDepth && depth != 0 {
			tr.Info("max depth reached", "path", curr, "depth", depth)
			return
		}

		for _, el := range temp {
			if el.Mask&files.MaskDirectory != 0 {
				// if !tr.SelFilter(el) {
				// 	tr.Debug("selection filtered out", "path", el.Path)
				// 	continue
				// }
				seeds = append(seeds, el.Path)
			}

			if !tr.ResFilter(el) {
				tr.Debug("result filtered out", "path", el.Path)
				continue
			}

			if depth >= tr.MinDepth {
				res = append(res, el)
			}
		}

		temp, curr, ok = src.Iter()
		for len(temp) == 0 {
			tr.Info("no elements received")
			switch {
			case !ok && len(seeds) == 0:
				return
			case !ok && len(seeds) > 0:
				depth++
				if depth >= tr.MaxDepth && tr.MaxDepth != 0 {
					tr.Info("max depth reached", "path", curr, "depth", depth)
					return
				}
				src.Seed(seeds)
				seeds = make([]string, 0)
			}
			temp, curr, ok = src.Iter()
		}
	}
	return
}

type Parser func(string) (res []*Element, err error)

func NewFSSource() *FSSource {
	src := &FSSource{
		Paths: []string{},
	}
	return src
}

func (s *FSSource) SeedFromElements(els []*Element) {
	for _, el := range els {
		if el.Mask&files.MaskDirectory != 0 {
			s.Paths = append(s.Paths, el.Path)
		} else {
			s.Els = append(s.Els, el)
		}
	}
}

// FSSource is a BFS source for traversing filesystem
type FSSource struct {
	Paths  []string
	Els    []*Element
	Parser Parser
}

func (s *FSSource) Iter() (res []*Element, curr string, ok bool) {
	if len(s.Els) > 0 {
		res = s.Els
		s.Els = nil
		return res, "FSSource inherited", true
	}

	if len(s.Paths) == 0 {
		return nil, "", false
	}

	curr = s.Paths[0]
	s.Paths = s.Paths[1:]
	res, err := s.Parser(curr)
	if err != nil {
		return nil, curr, false
	}
	return res, curr, true
}

func (s *FSSource) Seed(seeds []string) {
	s.Paths = append(s.Paths, seeds...)
}
