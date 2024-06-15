package yahw

import "io"

type Text string

func (Text) tag() {}

func (t Text) Render(w io.Writer) error {
	_, err := w.Write([]byte(t))
	return err
}
