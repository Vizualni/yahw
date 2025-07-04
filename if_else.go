package yahw

import (
	"context"
)

type IfElse struct {
	cond bool
	then Node
	els  Node
}

var _ Node = IfElse{}

func If(cond bool, then Node) IfElse {
	return IfElse{cond: cond, then: then}
}

func (ie IfElse) Else(els Node) IfElse {
	ie.els = els
	return ie
}

// Node implements Node.
func (ie IfElse) Node(ctx context.Context) Renderable {
	eval := ie.evaluate(ctx)
	if eval == nil {
		return nil
	}
	return eval
}

func (t IfElse) evaluate(ctx context.Context) Renderable {
	if t.cond {
		return t.then.Node(ctx)
	}
	if t.els != nil {
		return t.els.Node(ctx)
	}
	return nil
}
