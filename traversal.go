package ls

type FSTraverser struct {
	MaxDepth int
	MinDpeth int

	// Selection filter, determines which paths are and aren't traversed
	SelFilter func(*Element) bool
	// Result filter, determines which elments are and aren't included in result
	ResFilter func(*Element) bool

	Parser func(string) (res []*Element, err error)

	Logger
}

func (t *FSTraverser) init() {
	t.MaxDepth = 0
	t.MinDpeth = -1
	t.SelFilter = func(*Element) bool { return true }
	t.ResFilter = func(*Element) bool { return true }
	t.Logger = DummyLogger{}
}

func (tr *FSTraverser) Traverse(src SourceIter[*Element, string]) (res []*Element) {
	var depth int
	var lastEl string
	var lastPath string
	defer func() {
		tr.Info("traversal finished", "depth", depth, "last path", lastPath, "last el", lastEl, "total found", len(res))
	}()

	var seeds []string
	temp, curr, ok := src.Iter()
	tr.Info("starting traversal", "path", curr)
	for ok {
		lastPath = curr
		for _, el := range temp {
			lastEl = el.Path
			if el.Mask&MaskDirectory != 0 {
				if !tr.SelFilter(el) {
					tr.Debug("selection filtered out", "path", el.Path)
					continue
				}
				seeds = append(seeds, el.Path)
				continue
			}

			if !tr.ResFilter(el) {
				tr.Debug("result filtered out", "path", el.Path)
				continue
			}

			res = append(res, el)
		}

		temp, curr, ok = src.Iter()
		for len(temp) == 0 {
			tr.Info("no elements received", "path", curr)
			switch {
			case !ok && len(seeds) == 0:
				return
			case !ok && len(seeds) > 0:
				src.Seed(seeds)
				seeds = make([]string, 0)
				depth++
			}
			temp, curr, ok = src.Iter()
		}
	}
	return
}

type FSSource struct {
	Ind    int
	Paths  []string
	Parser func(string) (res []*Element, err error)
}

func (s *FSSource) Iter() (res []*Element, curr string, ok bool) {
	if s.Ind >= len(s.Paths) {
		return nil, "", false
	}

	curr = s.Paths[s.Ind]
	s.Ind++
	res, err := s.Parser(curr)
	if err != nil {
		return nil, curr, false
	}
	return res, curr, true
}

func (s *FSSource) Seed(seeds []string) {
	s.Paths = seeds
}

// func Traverse(opts *Options, yfn Yield, rfn ResultFn) {
// 	var depth int
// 	dirPaths := opts.Args
// 	for i, path := range dirPaths {
// 		dirPaths[i] = ResolveHome(path)
// 	}
// 	slog.Debug("traversing", "dirs", dirPaths)

// 	for els, ok := yfn(dirPaths); ok; els, ok = yfn(dirPaths) {
// 		dirPaths = make([]string, 0)

// 		for _, el := range els {
// 			if el.IsDir {
// 				dirPaths = append(dirPaths, el.Path)
// 				if opts.OnlyFiles {
// 					continue
// 				}
// 				if opts.DirOnly {
// 					rfn(el)
// 					continue
// 				}
// 			}

// 			if depth < opts.FromDepth {
// 				slog.Debug("skipping", "element", el.Path, "depth", depth)
// 				continue
// 			}

// 			rfn(el)
// 		}

// 		depth++
// 		if depth > opts.ToDepth {
// 			slog.Debug("reached max depth", "depth", depth, "todepth", opts.ToDepth)
// 			return
// 		}
// 	}
// }

// func GetYieldFs(opts *Options) Yield {
// 	parser := InitFileParser(opts)
// 	return func(paths []string) (els []*Element, ok bool) {
// 		slog.Debug("yielding", "paths", len(paths))
// 		if len(paths) == 0 {
// 			return nil, false
// 		}
// 		for _, path := range paths {
// 			var err error
// 			var finfos []fs.FileInfo
// 			switch {
// 			case opts.Archive && IsZipLike(path):
// 				finfos, err = TraverseZip(path)
// 			default:
// 				finfos, err = TraverseDir(path)
// 			}
// 			if err != nil {
// 				slog.Debug("error during traversal", "err", err)
// 				continue
// 			}

// 			for _, finfo := range finfos {
// 				el := parser(path, finfo)
// 				if el != nil {
// 					els = append(els, el)
// 				}
// 			}

// 		}
// 		if len(els) == 0 {
// 			slog.Debug("no elements found")
// 			return nil, false
// 		}
// 		ok = true
// 		slog.Debug("yielding", "elements", len(els))
// 		return
// 	}
// }
