package yahw

import (
	"testing"
)

func TestCreatingTags(t *testing.T) {
	foo := TagBuilder("foo")
	bar := TagBuilder("bar")
	single := SelfClosingTagBuilder("single")
	tt := []struct {
		Name string
		Tag  TagRenderer
		Exp  string
	}{
		{Name: "Simple foo tag", Tag: foo(), Exp: "<foo></foo>"},
		{Name: "Simple foo tag with bar content", Tag: foo(Attr("key", "value")), Exp: "<foo key=\"value\"></foo>"},

		{Name: "Tag with two attributes", Tag: foo(Attr("key1", "value1"), Attr("key2", "value2")), Exp: "<foo key1=\"value1\" key2=\"value2\"></foo>"},

		{Name: "Tag with a child", Tag: foo().X(bar()), Exp: "<foo><bar></bar></foo>"},

		{Name: "Tag with a child and attrs", Tag: foo(
			Attr("key1", "value1"),
			Attr("key2", "value2"),
		).X(bar(
			Attr("key3", "value3"),
			Attr("key4", "value4"),
		)), Exp: "<foo key1=\"value1\" key2=\"value2\"><bar key3=\"value3\" key4=\"value4\"></bar></foo>"},
		{Name: "Self-closing tag", Tag: single(), Exp: "<single />"},
		{Name: "Self-closing tag with attrs", Tag: single(Attr("key", "value")), Exp: "<single key=\"value\" />"},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			assertEqual(t, tc.Tag, tc.Exp)
		})
	}
}

func TestCreatingInvalidTags(t *testing.T) {
	tt := []struct {
		Name string
		Func func()
	}{
		{Name: "Empty tag", Func: func() {
			TagBuilder("")
		}},
		{Name: "Tag with space", Func: func() {
			TagBuilder("foo space")
		}},
		{Name: "Tag with newline", Func: func() {
			TagBuilder("foo\nnewline")
		}},
		{Name: "Tag with carriage return", Func: func() {
			TagBuilder("foo\rnewline")
		}},
		{Name: "Tag with tab", Func: func() {
			TagBuilder("foo\ttab")
		}},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			assertPanic(t, tc.Func)
		})
	}
}
