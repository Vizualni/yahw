package yahw

import (
	"io"
)

type AddAttributter interface {
	AddAttributes(...AttrRenderer) TagRenderer
}

type TagRenderer interface {
	TagRender(w io.Writer) error
}

type AttrRenderer interface {
	AttrRender(w io.Writer) error
}
