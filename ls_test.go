package ls

import (
	"log/slog"
	"testing"

	"github.com/periaate/ls/files"
	"github.com/periaate/ls/lfs"
)

func TestDo(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	res := Do(
		Combine(
			Dir(
				Paths("./test/"),
				LogWith(slog.Default()),
				Masks(Include, files.MaskDirectory),
			),
			Dir(
				Paths("./test/subfolder/"),
				LogWith(slog.Default()),
			),
		),
		SliceProcess("[:3]"),
		Sort(lfs.ByName),
	)

	if len(res) != 3 {
		t.Fatal("Traverse failed")
	}

	for _, el := range res {
		t.Log(el.Path)
	}
}
