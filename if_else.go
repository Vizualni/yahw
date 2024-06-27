package yahw

import "io"

type IfElseTag struct {
	cond bool
	then TagRenderer
	els  TagRenderer
}

func If(cond bool, then TagRenderer) IfElseTag {
	return IfElseTag{cond: cond, then: then}
}

func (ie IfElseTag) Else(els TagRenderer) IfElseTag {
	ie.els = els
	return ie
}

func (t IfElseTag) TagRender(w io.Writer) error {
	if t.cond {
		return t.then.TagRender(w)
	}
	return t.els.TagRender(w)
}

type IfElseAttr struct {
	cond bool
	then AttrRenderer
	els  AttrRenderer
}

func IfAttr(cond bool, then AttrRenderer) IfElseAttr {
	return IfElseAttr{cond: cond, then: then}
}

func (ie IfElseAttr) Else(els AttrRenderer) IfElseAttr {
	ie.els = els
	return ie
}

func (t IfElseAttr) AttrRender(w io.Writer) error {
	if t.cond {
		return t.then.AttrRender(w)
	}
	return t.els.AttrRender(w)
}
