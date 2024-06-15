package yahw

import (
	"io"
)

type TagRenderer interface {
	TagRender(w io.Writer) error
}

type AttrRenderer interface {
	AttrRender(w io.Writer) error
}
