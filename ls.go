package ls

// type Options struct {
// 	ToDepth   int
// 	FromDepth int
// 	Queries   []Query

// 	Args []string

// 	Filters   []func(*Element) bool
// 	Processes []func(els []*Element) []*Element

// 	Select []string
// }

const (
	ByNone SortBy = iota
	ByName
	ByMod
	BySize
	ByCreation // Available only on Windows
)

type SortBy uint8

type Element struct {
	Name string
	Path string // includes name, relative path to cwd
	Vany int64  // Unix time for mod|creation
	Mask uint32 // Bitmask of element types, file, dir, content types, hidden, etc.
}

func NoneFilter(*Element) bool              { return true }
func NoneProcess(els []*Element) []*Element { return els }

type Process func(inp []*Element) (out []*Element)
type Filter func(*Element) bool

type Result struct{ Files []*Element }

type SourceIter[T, S any] interface {
	Iter() (res []T, curr S, ok bool)
	Seed(seeds []S)
}

type Traverser[T, S any] interface {
	Traverse(SourceIter[T, S]) []T
}

type Logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
}

type DummyLogger struct{}

func (DummyLogger) Error(_ string, _ ...any) {}
func (DummyLogger) Info(_ string, _ ...any)  {}
func (DummyLogger) Warn(_ string, _ ...any)  {}
func (DummyLogger) Debug(_ string, _ ...any) {}

type Options struct {
	Sort     SortBy
	Hide     bool
	Archives bool
	// Format directory paths to end with "/". Used for internal logic, turning it
	// off will remove file|directory selection functionality
	WebStyle bool
}

func NewOptions() *Options {
	return &Options{
		Sort:     ByNone,
		Hide:     true,
		Archives: false,
		WebStyle: true,
	}
}
