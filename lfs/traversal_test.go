package lfs

import (
	"log/slog"
	"testing"
)

func TestTraversal(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ftsr := &FSTraverser{}
	ftsr.Init()

	ftsr.Logger = slog.Default()

	src := NewFSSource()
	src.Seed([]string{"../test/"})
	fsw := NewFSWorker()
	fsw.Logger = slog.Default()
	src.Parser = fsw.Parser()
	res := ftsr.Traverse(src)

	if len(res) == 0 {
		t.Fatal("Traverse failed")
	}

	if len(res) != 7 {
		t.Fatalf("Expected 7 elements, got %d", len(res))
	}

	for _, el := range res {
		t.Log(el.Path)
	}
}
