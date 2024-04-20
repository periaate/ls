package ls

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/periaate/common"
	"github.com/periaate/ls/files"
	"github.com/periaate/ls/lfs"
)

func QueryAsFilter(qr ...string) Filter {
	scorer := common.GetScoringFunction(qr, 3)
	return func(e *lfs.Element) bool {
		score := scorer(e.Name)
		slog.Debug("query filter", "name", e.Name, "score", score)
		return score != 0
	}
}

func ParseSearch(args []string) Filter {
	slog.Debug("search args", "args", args)
	filters := []func(*lfs.Element) bool{}

	for _, arg := range args {
		q := Query{Include: true}
		switch {
		case len(arg) < 2:
			continue
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
		default:
			q.Kind = Substring
			q.Value = arg
		}
		filters = append(filters, q.GetFilter())
	}

	return common.All(true, filters...)
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
		slog.Debug("no mask found", "args", args)
		return NoneFilter
	}
	return q.GetFilter()
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

func MaskFilter(mask uint32) Filter {
	return func(e *lfs.Element) bool {
		return files.ExtToMaskMap[filepath.Ext(e.Name)]&mask != 0
	}
}

func ExactFilter(search string) Filter {
	return func(e *lfs.Element) bool { return search == e.Name }
}
func SubstringFilter(search string) Filter {
	search = strings.ToLower(search)
	return func(e *lfs.Element) bool {
		r := strings.Contains(strings.ToLower(e.Name), search)
		slog.Debug("substring filter", "name", e.Name, "search", search, "result", r)
		return r
	}
}
