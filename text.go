package yahw

import "io"

type Text string

func (t Text) TagRender(w io.Writer) error {
	_, err := w.Write([]byte(t))
	return err
}
