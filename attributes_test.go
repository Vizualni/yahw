package yahw

import (
	"strings"
	"testing"
)

func assertEqual(t *testing.T, r Renderable, expected string) {
	strbuf := &strings.Builder{}
	err := r.Render(strbuf)
	if err != nil {
		t.Errorf("Error rendering: %s", err)
	}

	if strbuf.String() != expected {
		t.Errorf("Expected %s, got %s", expected, strbuf.String())
	}
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, got nil")
		}
	}()
	f()
}

func TestWritingAttributes(t *testing.T) {
	tt := []struct {
		Name string
		Attr Attribute
		Exp  string
	}{
		{Name: "Single", Attr: Attribute{"foo", "bar"}, Exp: `foo="bar"`},
		{Name: "Empty", Attr: Attribute{"foo", ""}, Exp: `foo=""`},
		{Name: "And", Attr: Attribute{"foo", "bar&baz"}, Exp: `foo="bar&amp;baz"`},
		{Name: "Single quote", Attr: Attribute{"foo", "bar'baz"}, Exp: `foo="bar&#39;baz"`},
		{Name: "Forward slash", Attr: Attribute{"foo", "bar/baz"}, Exp: `foo="bar/baz"`},
		{Name: "Less than", Attr: Attribute{"foo", "bar<baz"}, Exp: `foo="bar&lt;baz"`},
		{Name: "Greater than", Attr: Attribute{"foo", "bar>baz"}, Exp: `foo="bar&gt;baz"`},
		{Name: "Double quote", Attr: Attribute{"foo", `bar"baz`}, Exp: `foo="bar&#34;baz"`},

		{Name: "Newline", Attr: Attribute{"foo", "bar\nbaz"}, Exp: "foo=\"bar\nbaz\""},
		{Name: "Carriage return", Attr: Attribute{"foo", "bar\rbaz"}, Exp: "foo=\"bar\rbaz\""},
		{Name: "Tab", Attr: Attribute{"foo", "bar\tbaz"}, Exp: "foo=\"bar\tbaz\""},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			assertEqual(t, tc.Attr, tc.Exp)
		})
	}
}

func TestCreatingInvalidAttrs(t *testing.T) {
	tt := []struct {
		Name     string
		AttrName string
		AttrVal  string
	}{
		{Name: "Empty", AttrName: "", AttrVal: "bar"},
		{Name: "Contains space", AttrName: "foo bar", AttrVal: "baz"},
		{Name: "Contains equals", AttrName: "foo=bar", AttrVal: "baz"},
		{Name: "Contains newline", AttrName: "foo\nbar", AttrVal: "baz"},
		{Name: "Contains carriage return", AttrName: "foo\rbar", AttrVal: "baz"},
		{Name: "Contains tab", AttrName: "foo\tbar", AttrVal: "baz"},
		{Name: "Contains double quote", AttrName: "foo\"bar", AttrVal: "baz"},
		{Name: "Contains single quote", AttrName: "foo'bar", AttrVal: "baz"},
		{Name: "Contains backtick", AttrName: "foo`bar", AttrVal: "baz"},
		{Name: "Contains backslash", AttrName: "foo\\bar", AttrVal: "baz"},
		{Name: "Contains forward slash", AttrName: "foo/bar", AttrVal: "baz"},
		{Name: "Contains less than", AttrName: "foo<bar", AttrVal: "baz"},
		{Name: "Contains greater than", AttrName: "foo>bar", AttrVal: "baz"},
		{Name: "Contains ampersand", AttrName: "foo&bar", AttrVal: "baz"},
		{Name: "Contains pipe", AttrName: "foo|bar", AttrVal: "baz"},
		{Name: "Contains exclamation mark", AttrName: "foo!bar", AttrVal: "baz"},
		{Name: "Contains question mark", AttrName: "foo?bar", AttrVal: "baz"},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			assertPanic(t, func() { Attr(tc.AttrName, tc.AttrVal) })
		})
	}
}
