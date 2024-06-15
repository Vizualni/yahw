package yahw

import (
	"io"
)

type Renderable interface {
	Render(w io.Writer) error
}

type tagger interface {
	Renderable
	tag()
}
type attrer interface {
	Renderable
	attr()
}
