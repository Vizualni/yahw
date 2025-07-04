package yahw

import (
	"context"
	"io"
)

type Text string

func (t Text) tag()                                {}
func (t Text) Node(ctx context.Context) Renderable { return t }

func (t Text) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(t))
	return err
}
