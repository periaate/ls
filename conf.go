package ls

// type OptFn func(*Options)

// func Paths(paths ...string) OptFn {
// 	return func(opts *Options) { opts.Args = append(opts.Args, paths...) }
// }

// func FromDepth(depth int) OptFn {
// 	return func(opts *Options) { opts.FromDepth = depth }
// }

// func ToDepth(depth int) OptFn {
// 	return func(opts *Options) { opts.ToDepth = depth }
// }

// func DepthPattern(pat string) OptFn {
// 	return func(opts *Options) {
// 		from, to, err := slice.ParsePattern(pat, math.MaxInt)
// 		if err != nil {
// 			return
// 		}
// 		opts.FromDepth = from
// 		opts.ToDepth = to
// 	}
// }

// const (
// 	Files = false
// 	Dirs  = true
// )

// func Only(b bool) OptFn {
// 	return func(opts *Options) {
// 		switch b {
// 		case Files:
// 			opts.OnlyFiles = true
// 		case Dirs:
// 			opts.DirOnly = true
// 		}
// 	}
// }

// func Search(queries ...string) OptFn {
// 	return func(opts *Options) {
// 		f := ParseSearch(queries)
// 		if f != nil {
// 			opts.Filters = append(opts.Filters, f)
// 		}
// 	}
// }

// const (
// 	Include = true
// 	Exclude = false
// )

// func Kind(inc bool, kinds ...string) OptFn {
// 	return func(opts *Options) {
// 		f := ParseKind(kinds, inc)
// 		if f != nil {
// 			opts.Filters = append(opts.Filters, f)
// 		}
// 	}
// }

// func Recurse(opts *Options) {
// 	opts.ToDepth = math.MaxInt64
// }

// func NoHide(opt *Options) {
// 	opt.NoHide = true
// }

// func Do(args ...OptFn) []*Element {
// 	opts := &Options{}
// 	for _, fn := range args {
// 		fn(opts)
// 	}

// 	if len(opts.Args) == 0 {
// 		opts.Args = []string{"./"}
// 	}

// 	return Run(opts)
// }
