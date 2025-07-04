package yahw

import (
	"context"
	"fmt"
	"io"
)

type Raw string

var (
	_ Node       = Raw("")
	_ Renderable = Raw("")
	_ taggable   = Raw("")
)

func (r Raw) tag()                                {}
func (r Raw) Node(ctx context.Context) Renderable { return r }

func (r Raw) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(r))
	return err
}

func RawFormat(format string, args ...any) Raw {
	return Raw(fmt.Sprintf(format, args...))
}
