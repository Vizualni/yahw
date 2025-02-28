package yahw

import (
	"io"
)

type Renderable interface {
	Render(w io.Writer) error
}

type Tag interface {
	Tag() Renderable
}

type Attr interface {
	Attr() Renderable
}
