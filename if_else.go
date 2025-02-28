package yahw

import "io"

type IfElseTag struct {
	cond bool
	then Tag
	els  Tag
}

func If(cond bool, then Tag) IfElseTag {
	return IfElseTag{cond: cond, then: then}
}

func (ie IfElseTag) Else(els Tag) IfElseTag {
	ie.els = els
	return ie
}

func (t IfElseTag) Render(w io.Writer) error {
	if t.cond {
		if t.then == nil {
			return nil
		}
		return t.then.Tag().Render(w)
	}
	if t.els == nil {
		return nil
	}
	return t.els.Tag().Render(w)
}

type IfElseAttr struct {
	cond bool
	then Attr
	els  Attr
}

func IfAttr(cond bool, then Attr) IfElseAttr {
	return IfElseAttr{cond: cond, then: then}
}

func (ie IfElseAttr) Else(els Attr) IfElseAttr {
	ie.els = els
	return ie
}

func (t IfElseAttr) Render(w io.Writer) error {
	if t.cond {
		if t.then == nil {
			return nil
		}
		return t.then.Attr().Render(w)
	}
	if t.els == nil {
		return nil
	}
	return t.els.Attr().Render(w)
}
