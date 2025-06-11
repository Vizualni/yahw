package yahw

import "io"

type IfElseTag struct {
	cond bool
	then taggable
	els  taggable
}

var _ Node = IfElseTag{}

func If(cond bool, then taggable) IfElseTag {
	return IfElseTag{cond: cond, then: then}
}

func (ie IfElseTag) Else(els taggable) IfElseTag {
	ie.els = els
	return ie
}

// Node implements Node.
func (ie IfElseTag) Node() Renderable {
	return ie
}

func (t IfElseTag) Render(w io.Writer) error {
	if t.cond {
		if t.then == nil {
			return nil
		}
		return t.then.Render(w)
	}
	if t.els == nil {
		return nil
	}
	return t.els.Render(w)
}

type IfElseAttr struct {
	cond bool
	then attrable
	els  attrable
}

func IfAttr(cond bool, then attrable) IfElseAttr {
	return IfElseAttr{cond: cond, then: then}
}

func (ie IfElseAttr) Else(els attrable) IfElseAttr {
	ie.els = els
	return ie
}

func (t IfElseAttr) Render(w io.Writer) error {
	if t.cond {
		if t.then == nil {
			return nil
		}
		return t.then.Render(w)
	}
	if t.els == nil {
		return nil
	}
	return t.els.Render(w)
}

func (t IfElseAttr) Node() Renderable {
	return t
}
