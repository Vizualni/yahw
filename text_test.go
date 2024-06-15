package yahw

import "testing"

func TestTextRendering(t *testing.T) {
	assertEqual(t, Text("foo"), "foo")
	assertEqual(t, Text("bar"), "bar")
	assertEqual(t, Text("baz"), "baz")
}

func TestTextRenderingWithinTags(t *testing.T) {
	T1 := TagBuilder("T1")

	assertEqual(t, T1().X(Text("foo")), "<T1>foo</T1>")
	assertEqual(t, T1().X(Text("foo\nbar")), "<T1>foo\nbar</T1>")
}
