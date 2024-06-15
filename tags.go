package yahw

import (
	"io"
)

func isValidTagName(tagName string) bool {
	if len(tagName) == 0 {
		return false
	}
	for _, c := range tagName {
		switch {
		case '0' <= c && c <= '9':
		case 'a' <= c && c <= 'z':
		case 'A' <= c && c <= 'Z':
		default:
			return false
		}
	}
	return true
}

func TagBuilder(tagName string) func(...AttrRenderer) Tag {
	if !isValidTagName(tagName) {
		panic("Invalid tag name: " + tagName)
	}

	return func(attrs ...AttrRenderer) Tag {
		return Tag{
			tagName: tagName,
			attrs:   attrs,
		}
	}
}

func SelfClosingTagBuilder(tagName string) func(...AttrRenderer) SelfClosingTag {
	if !isValidTagName(tagName) {
		panic("Invalid self closing tag name: " + tagName)
	}
	return func(attrs ...AttrRenderer) SelfClosingTag {
		return SelfClosingTag{
			tagName: tagName,
			attrs:   attrs,
		}
	}
}

type SelfClosingTag struct {
	tagName string
	attrs   []AttrRenderer
}

func (t SelfClosingTag) TagRender(w io.Writer) error {
	_, err := w.Write([]byte("<" + t.tagName))
	if err != nil {
		return err
	}

	if len(t.attrs) > 0 {
		w.Write([]byte(" "))
	}

	for idx, attr := range t.attrs {
		err = attr.AttrRender(w)
		if err != nil {
			return err
		}

		if idx < len(t.attrs)-1 {
			w.Write([]byte(" "))
		}
	}

	_, err = w.Write([]byte(" />"))
	if err != nil {
		return err
	}

	return nil
}

type Tag struct {
	tagName  string
	attrs    []AttrRenderer
	children []TagRenderer
}

func (t Tag) tag()      {}
func (t Tag) yahwNode() {}

func (t Tag) X(children ...TagRenderer) Tag {
	clone := t.clone()
	clone.children = append(clone.children, children...)
	return clone
}

func (t Tag) clone() Tag {
	attrs := make([]AttrRenderer, len(t.attrs))
	copy(attrs, t.attrs)

	children := make([]TagRenderer, len(t.children))
	copy(children, t.children)

	return Tag{
		tagName:  t.tagName,
		attrs:    attrs,
		children: children,
	}
}

func (t Tag) TagRender(w io.Writer) error {
	_, err := w.Write([]byte("<" + t.tagName))
	if err != nil {
		return err
	}

	if len(t.attrs) > 0 {
		w.Write([]byte(" "))
	}

	for idx, attr := range t.attrs {
		err = attr.AttrRender(w)
		if err != nil {
			return err
		}

		if idx < len(t.attrs)-1 {
			w.Write([]byte(" "))
		}
	}

	_, err = w.Write([]byte(">"))
	if err != nil {
		return err
	}

	for _, child := range t.children {
		err = child.TagRender(w)
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte("</" + t.tagName + ">"))
	if err != nil {
		return err
	}

	return nil
}

type HTML5Doctype struct {
	children []TagRenderer
}

func (HTML5Doctype) tag() {}
func (t HTML5Doctype) Render(w io.Writer) error {
	_, err := w.Write([]byte("<!DOCTYPE html>"))
	if err != nil {
		return err
	}
	for _, child := range t.children {
		err = child.TagRender(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t HTML5Doctype) X(children ...TagRenderer) HTML5Doctype {
	clone := t.clone()
	clone.children = append(clone.children, children...)
	return clone
}

func (t HTML5Doctype) clone() HTML5Doctype {
	children := make([]TagRenderer, len(t.children))
	copy(children, t.children)

	return HTML5Doctype{
		children: children,
	}
}

// All known HTML5 tags

func NewHTML5Doctype() HTML5Doctype { return HTML5Doctype{} }

func A(attrs ...AttrRenderer) Tag                 { return TagBuilder("a")(attrs...) }
func Abbr(attrs ...AttrRenderer) Tag              { return TagBuilder("abbr")(attrs...) }
func Address(attrs ...AttrRenderer) Tag           { return TagBuilder("address")(attrs...) }
func Area(attrs ...AttrRenderer) SelfClosingTag   { return SelfClosingTagBuilder("area")(attrs...) }
func Article(attrs ...AttrRenderer) Tag           { return TagBuilder("article")(attrs...) }
func Aside(attrs ...AttrRenderer) Tag             { return TagBuilder("aside")(attrs...) }
func Audio(attrs ...AttrRenderer) Tag             { return TagBuilder("audio")(attrs...) }
func B(attrs ...AttrRenderer) Tag                 { return TagBuilder("b")(attrs...) }
func Base(attrs ...AttrRenderer) SelfClosingTag   { return SelfClosingTagBuilder("base")(attrs...) }
func Bdi(attrs ...AttrRenderer) Tag               { return TagBuilder("bdi")(attrs...) }
func Bdo(attrs ...AttrRenderer) Tag               { return TagBuilder("bdo")(attrs...) }
func Blockquote(attrs ...AttrRenderer) Tag        { return TagBuilder("blockquote")(attrs...) }
func Body(attrs ...AttrRenderer) Tag              { return TagBuilder("body")(attrs...) }
func Br(attrs ...AttrRenderer) SelfClosingTag     { return SelfClosingTagBuilder("br")(attrs...) }
func Button(attrs ...AttrRenderer) Tag            { return TagBuilder("button")(attrs...) }
func Canvas(attrs ...AttrRenderer) Tag            { return TagBuilder("canvas")(attrs...) }
func Caption(attrs ...AttrRenderer) Tag           { return TagBuilder("caption")(attrs...) }
func Cite(attrs ...AttrRenderer) Tag              { return TagBuilder("cite")(attrs...) }
func Code(attrs ...AttrRenderer) Tag              { return TagBuilder("code")(attrs...) }
func Col(attrs ...AttrRenderer) SelfClosingTag    { return SelfClosingTagBuilder("col")(attrs...) }
func Colgroup(attrs ...AttrRenderer) Tag          { return TagBuilder("colgroup")(attrs...) }
func Data(attrs ...AttrRenderer) Tag              { return TagBuilder("data")(attrs...) }
func Datalist(attrs ...AttrRenderer) Tag          { return TagBuilder("datalist")(attrs...) }
func Dd(attrs ...AttrRenderer) Tag                { return TagBuilder("dd")(attrs...) }
func Del(attrs ...AttrRenderer) Tag               { return TagBuilder("del")(attrs...) }
func Details(attrs ...AttrRenderer) Tag           { return TagBuilder("details")(attrs...) }
func Dfn(attrs ...AttrRenderer) Tag               { return TagBuilder("dfn")(attrs...) }
func Dialog(attrs ...AttrRenderer) Tag            { return TagBuilder("dialog")(attrs...) }
func Div(attrs ...AttrRenderer) Tag               { return TagBuilder("div")(attrs...) }
func Dl(attrs ...AttrRenderer) Tag                { return TagBuilder("dl")(attrs...) }
func Dt(attrs ...AttrRenderer) Tag                { return TagBuilder("dt")(attrs...) }
func Em(attrs ...AttrRenderer) Tag                { return TagBuilder("em")(attrs...) }
func Embed(attrs ...AttrRenderer) SelfClosingTag  { return SelfClosingTagBuilder("embed")(attrs...) }
func Fieldset(attrs ...AttrRenderer) Tag          { return TagBuilder("fieldset")(attrs...) }
func Figcaption(attrs ...AttrRenderer) Tag        { return TagBuilder("figcaption")(attrs...) }
func Figure(attrs ...AttrRenderer) Tag            { return TagBuilder("figure")(attrs...) }
func Footer(attrs ...AttrRenderer) Tag            { return TagBuilder("footer")(attrs...) }
func Form(attrs ...AttrRenderer) Tag              { return TagBuilder("form")(attrs...) }
func H1(attrs ...AttrRenderer) Tag                { return TagBuilder("h1")(attrs...) }
func H2(attrs ...AttrRenderer) Tag                { return TagBuilder("h2")(attrs...) }
func H3(attrs ...AttrRenderer) Tag                { return TagBuilder("h3")(attrs...) }
func H4(attrs ...AttrRenderer) Tag                { return TagBuilder("h4")(attrs...) }
func H5(attrs ...AttrRenderer) Tag                { return TagBuilder("h5")(attrs...) }
func H6(attrs ...AttrRenderer) Tag                { return TagBuilder("h6")(attrs...) }
func Head(attrs ...AttrRenderer) Tag              { return TagBuilder("head")(attrs...) }
func Header(attrs ...AttrRenderer) Tag            { return TagBuilder("header")(attrs...) }
func Hr(attrs ...AttrRenderer) SelfClosingTag     { return SelfClosingTagBuilder("hr")(attrs...) }
func HTML(attrs ...AttrRenderer) Tag              { return TagBuilder("html")(attrs...) }
func I(attrs ...AttrRenderer) Tag                 { return TagBuilder("i")(attrs...) }
func Iframe(attrs ...AttrRenderer) Tag            { return TagBuilder("iframe")(attrs...) }
func Img(attrs ...AttrRenderer) SelfClosingTag    { return SelfClosingTagBuilder("img")(attrs...) }
func Input(attrs ...AttrRenderer) SelfClosingTag  { return SelfClosingTagBuilder("input")(attrs...) }
func Ins(attrs ...AttrRenderer) Tag               { return TagBuilder("ins")(attrs...) }
func Kbd(attrs ...AttrRenderer) Tag               { return TagBuilder("kbd")(attrs...) }
func Label(attrs ...AttrRenderer) Tag             { return TagBuilder("label")(attrs...) }
func Legend(attrs ...AttrRenderer) Tag            { return TagBuilder("legend")(attrs...) }
func Li(attrs ...AttrRenderer) Tag                { return TagBuilder("li")(attrs...) }
func Link(attrs ...AttrRenderer) SelfClosingTag   { return SelfClosingTagBuilder("link")(attrs...) }
func Main(attrs ...AttrRenderer) Tag              { return TagBuilder("main")(attrs...) }
func Map(attrs ...AttrRenderer) Tag               { return TagBuilder("map")(attrs...) }
func Mark(attrs ...AttrRenderer) Tag              { return TagBuilder("mark")(attrs...) }
func Meta(attrs ...AttrRenderer) SelfClosingTag   { return SelfClosingTagBuilder("meta")(attrs...) }
func Meter(attrs ...AttrRenderer) Tag             { return TagBuilder("meter")(attrs...) }
func Nav(attrs ...AttrRenderer) Tag               { return TagBuilder("nav")(attrs...) }
func Noscript(attrs ...AttrRenderer) Tag          { return TagBuilder("noscript")(attrs...) }
func Object(attrs ...AttrRenderer) Tag            { return TagBuilder("object")(attrs...) }
func Ol(attrs ...AttrRenderer) Tag                { return TagBuilder("ol")(attrs...) }
func Optgroup(attrs ...AttrRenderer) Tag          { return TagBuilder("optgroup")(attrs...) }
func Option(attrs ...AttrRenderer) Tag            { return TagBuilder("option")(attrs...) }
func Output(attrs ...AttrRenderer) Tag            { return TagBuilder("output")(attrs...) }
func P(attrs ...AttrRenderer) Tag                 { return TagBuilder("p")(attrs...) }
func Param(attrs ...AttrRenderer) SelfClosingTag  { return SelfClosingTagBuilder("param")(attrs...) }
func Picture(attrs ...AttrRenderer) Tag           { return TagBuilder("picture")(attrs...) }
func Pre(attrs ...AttrRenderer) Tag               { return TagBuilder("pre")(attrs...) }
func Progress(attrs ...AttrRenderer) Tag          { return TagBuilder("progress")(attrs...) }
func Q(attrs ...AttrRenderer) Tag                 { return TagBuilder("q")(attrs...) }
func Rp(attrs ...AttrRenderer) Tag                { return TagBuilder("rp")(attrs...) }
func Rt(attrs ...AttrRenderer) Tag                { return TagBuilder("rt")(attrs...) }
func Ruby(attrs ...AttrRenderer) Tag              { return TagBuilder("ruby")(attrs...) }
func S(attrs ...AttrRenderer) Tag                 { return TagBuilder("s")(attrs...) }
func Samp(attrs ...AttrRenderer) Tag              { return TagBuilder("samp")(attrs...) }
func Script(attrs ...AttrRenderer) Tag            { return TagBuilder("script")(attrs...) }
func Section(attrs ...AttrRenderer) Tag           { return TagBuilder("section")(attrs...) }
func Select(attrs ...AttrRenderer) Tag            { return TagBuilder("select")(attrs...) }
func Slot(attrs ...AttrRenderer) Tag              { return TagBuilder("slot")(attrs...) }
func Small(attrs ...AttrRenderer) Tag             { return TagBuilder("small")(attrs...) }
func Source(attrs ...AttrRenderer) SelfClosingTag { return SelfClosingTagBuilder("source")(attrs...) }
func Span(attrs ...AttrRenderer) Tag              { return TagBuilder("span")(attrs...) }
func Strong(attrs ...AttrRenderer) Tag            { return TagBuilder("strong")(attrs...) }
func Style(attrs ...AttrRenderer) Tag             { return TagBuilder("style")(attrs...) }
func Sub(attrs ...AttrRenderer) Tag               { return TagBuilder("sub")(attrs...) }
func Summary(attrs ...AttrRenderer) Tag           { return TagBuilder("summary")(attrs...) }
func Sup(attrs ...AttrRenderer) Tag               { return TagBuilder("sup")(attrs...) }
func Table(attrs ...AttrRenderer) Tag             { return TagBuilder("table")(attrs...) }
func Tbody(attrs ...AttrRenderer) Tag             { return TagBuilder("tbody")(attrs...) }
func Td(attrs ...AttrRenderer) Tag                { return TagBuilder("td")(attrs...) }
func Template(attrs ...AttrRenderer) Tag          { return TagBuilder("template")(attrs...) }
func Textarea(attrs ...AttrRenderer) Tag          { return TagBuilder("textarea")(attrs...) }
func Tfoot(attrs ...AttrRenderer) Tag             { return TagBuilder("tfoot")(attrs...) }
func Th(attrs ...AttrRenderer) Tag                { return TagBuilder("th")(attrs...) }
func Thead(attrs ...AttrRenderer) Tag             { return TagBuilder("thead")(attrs...) }
func Time(attrs ...AttrRenderer) Tag              { return TagBuilder("time")(attrs...) }
func Title(attrs ...AttrRenderer) Tag             { return TagBuilder("title")(attrs...) }
func Tr(attrs ...AttrRenderer) Tag                { return TagBuilder("tr")(attrs...) }
func Track(attrs ...AttrRenderer) SelfClosingTag  { return SelfClosingTagBuilder("track")(attrs...) }
func U(attrs ...AttrRenderer) Tag                 { return TagBuilder("u")(attrs...) }
func Ul(attrs ...AttrRenderer) Tag                { return TagBuilder("ul")(attrs...) }
func Var(attrs ...AttrRenderer) Tag               { return TagBuilder("var")(attrs...) }
func Video(attrs ...AttrRenderer) Tag             { return TagBuilder("video")(attrs...) }
func Wbr(attrs ...AttrRenderer) SelfClosingTag    { return SelfClosingTagBuilder("wbr")(attrs...) }
