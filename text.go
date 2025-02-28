package yahw

import "io"

type Text string

func (t Text) tag()             {}
func (t Text) Node() Renderable { return t }

func (t Text) Render(w io.Writer) error {
	_, err := w.Write([]byte(t))
	return err
}
