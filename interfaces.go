package yahw

import (
	"context"
	"io"
)

type Renderable interface {
	Render(ctx context.Context, w io.Writer) error
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
	Node(ctx context.Context) Renderable
}
