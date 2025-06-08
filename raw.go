package yahw

import "io"

type Raw string

var (
	_ Node       = Raw("")
	_ Renderable = Raw("")
)

func (r Raw) Node() Renderable { return r }

func (r Raw) Render(w io.Writer) error {
	_, err := w.Write([]byte(r))
	return err
}
