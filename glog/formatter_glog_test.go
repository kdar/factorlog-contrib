package glog

import (
	"testing"
	"time"

	log "github.com/kdar/factorlog"
)

var fmtTestsContext = log.LogContext{
	Time:     time.Unix(0, 1389223634123456789).In(time.UTC),
	Severity: log.PANIC,
	File:     "path/to/testing.go",
	Line:     391,
	Format:   nil,
	Args:     []interface{}{"hello there!"},
	Function: "some crazy/path.path/pkg.(*Type).Function",
	Pid:      1234,
}

func TestGlogFormatter(t *testing.T) {
	f := NewGlogFormatter()
	expect := "P0108 23:27:14.123456 01234 testing.go:391] hello there!\n"
	out := string(f.Format(fmtTestsContext))
	if expect != out {
		t.Fatalf("\nexpected: %#v\ngot:      %#v", expect, out)
	}
}

func BenchmarkGlogFormatter(b *testing.B) {
	f := NewGlogFormatter()

	b.ResetTimer()
	for x := 0; x < b.N; x++ {
		f.Format(fmtTestsContext)
	}
}
