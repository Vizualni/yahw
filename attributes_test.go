package yahw

import (
	"strings"
	"testing"
)

func assertEqual(t *testing.T, r any, expected string) {
	strbuf := &strings.Builder{}
	var err error
	switch r := r.(type) {
	case Renderable:
		err = r.Render(strbuf)
	default:
		t.Errorf("Unknown type: %T", r)
	}
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
			assertPanic(t, func() { BuildAttr(tc.AttrName, tc.AttrVal) })
		})
	}
}

func TestClassNames(t *testing.T) {
	tt := []struct {
		Name    string
		Classes string
		Expect  string
	}{
		{Name: "Single", Classes: "foo", Expect: `class="foo"`},
		{Name: "Multiple", Classes: "foo bar", Expect: `class="foo bar"`},
		{Name: "Empty", Classes: "", Expect: `class=""`},
		{Name: "Duplicate", Classes: "foo foo", Expect: `class="foo"`},
		{Name: "Duplicate with space", Classes: "foo  foo", Expect: `class="foo"`},
		{Name: "Duplicate with other classes", Classes: "foo bar foo", Expect: `class="foo bar"`},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			assertEqual(t, Classes(tc.Classes), tc.Expect)
		})
	}
}
