package yahw

import (
	"io"
)

type Renderable interface {
	Render(w io.Writer) error
}

type taggable interface {
	Renderable
	tag()
}

type attrable interface {
	Renderable
	attr()
}

type Node interface {
	Node() Renderable
}
