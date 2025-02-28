package yahw_test

import (
	"strings"
	"testing"

	. "github.com/vizualni/yahw"
)

func TestSomewhatRealworldExample(t *testing.T) {
	root := NewHTML5Doctype(
		HTML(
			Title(Text("Hello, World!")),
			Style(Text("body { background-color: #f0f0f0; }")),
			Body(
				H1(Text("Hello, World!")),
				P(StyleAttr("color: red;"), Text("This is a paragraph.")),
			),
		),
	)

	strbuf := &strings.Builder{}
	err := root.Render(strbuf)
	if err != nil {
		t.Errorf("Error rendering: %s", err)
	}

	expected := `<!DOCTYPE html><html><title>Hello, World!</title><style>body { background-color: #f0f0f0; }</style><body><h1>Hello, World!</h1><p style="color: red;">This is a paragraph.</p></body></html>`

	if strbuf.String() != expected {
		t.Errorf("Expected %s, got %s", expected, strbuf.String())
	}
}
