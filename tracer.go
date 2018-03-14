package genesis

import (
	"fmt"
	"io"
)

type tracer struct {
	w io.Writer
}

func newTracer(w io.Writer) *tracer {
	return &tracer{w: w}
}

func (w *tracer) Write(p []byte) (n int, err error) {
	if w == nil {
		return 0, nil
	}
	return w.w.Write(p)
}

func (w *tracer) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, format, a...)
}
