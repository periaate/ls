package ls

import (
	"log/slog"
	"strings"

	"github.com/periaate/common"
	"github.com/periaate/ls/files"
	"github.com/periaate/ls/lfs"
	"github.com/periaate/slice"
)

var glog common.Logger = common.DummyLogger{}

func SetGlobalLogger(l common.Logger) {
	glog = l
}

func QueryAsFilter(qr ...string) Filter {
	scorer := common.GetScoringFunction(qr, 3)
	return func(e *lfs.Element) bool {
		score := scorer(e.Name)
		slog.Debug("query filter", "name", e.Name, "score", score)
		return score != 0
	}
}

func ParseSearch(args []string) Filter {
	glog.Debug("parsing search filter", "cnt", len(args), "args", args)
	filters := []func(*lfs.Element) bool{}

	for _, arg := range args {
		q := Query{Include: true}
		switch {
		case len(arg) < 2:
			q.Kind = Substring
			q.Value = arg
		case arg[:2] == "!=":
			arg = arg[1:]
			q.Include = false
			fallthrough
		case arg[0] == '=':
			arg = arg[1:]
			q.Value = arg
			q.Kind = Exact
		case arg[0] == '!':
			arg = arg[1:]
			q.Include = false
			fallthrough
		case arg[0] == '_':
			fallthrough
		default:
			q.Kind = Substring
			q.Value = arg
		}
		glog.Debug("search substring arg", "arg", q.Value, "kind", q.Kind, "include", q.Include)
		filters = append(filters, q.GetFilter())
	}

	return common.Any(true, filters...)
}

func ParseKind(args []string, inc bool) Filter {
	q := Query{
		Kind:    MaskK,
		Include: inc,
	}
	for _, arg := range args {
		q.Mask |= files.StrToMask(arg)
	}
	if q.Mask == 0 {
		glog.Debug("no mask found", "args", args)
		return NoneFilter
	}
	return q.GetFilter()
}

func MaskFilter(mask uint32) Filter {
	return func(e *lfs.Element) bool {
		return e.Mask&mask != 0
	}
}

type QueryKind [2]bool

var (
	Substring = QueryKind{false, false}
	Fuzzy     = QueryKind{true, false}
	Exact     = QueryKind{false, true}
	MaskK     = QueryKind{true, true}
)

type Query struct {
	Value   string
	Include bool
	Mask    uint32
	Kind    QueryKind
}

func (q Query) GetFilter() (f Filter) {
	switch q.Kind {
	case Fuzzy:
		f = QueryAsFilter(q.Value)
	case Exact:
		f = ExactFilter(q.Value)
	case MaskK:
		f = MaskFilter(q.Mask)
	case Substring:
		fallthrough
	default:
		f = SubstringFilter(q.Value)
	}

	if !q.Include {
		f = common.Negate(f)
	}
	return f
}

func ExactFilter(search string) Filter {
	return func(e *lfs.Element) bool { return search == e.Name }
}
func SubstringFilter(search string) Filter {
	search = strings.ToLower(search)
	return func(e *lfs.Element) bool {
		r := strings.Contains(strings.ToLower(e.Name), search)
		glog.Debug("substring filter", "name", e.Name, "search", search, "result", r)
		return r
	}
}

func SliceProcess(pattern string) Process {
	exp := slice.NewExpression[*lfs.Element]()
	exp.Parse(pattern)

	return func(filenames []*lfs.Element) (res []*lfs.Element) {
		res, err := exp.Eval(filenames)
		if err != nil {
			glog.Error("error in Slice", "error", err)
		}
		return
	}
}
