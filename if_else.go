package yahw

import "io"

type IfElseTag struct {
	cond bool
	then Node
	els  Node
}

var _ Node = IfElseTag{}

func If(cond bool, then Node) IfElseTag {
	return IfElseTag{cond: cond, then: then}
}

func (ie IfElseTag) Else(els Node) IfElseTag {
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
		return t.then.Node().Render(w)
	}
	if t.els == nil {
		return nil
	}
	return t.els.Node().Render(w)
}
