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

func TagBuilder(tagName string) func(...attrer) Tag {
	if !isValidTagName(tagName) {
		panic("Invalid tag name: " + tagName)
	}

	return func(attrs ...attrer) Tag {
		return Tag{
			tagName: tagName,
			attrs:   attrs,
		}
	}
}

func SelfClosingTagBuilder(tagName string) func(...attrer) SelfClosingTag {
	return func(attrs ...attrer) SelfClosingTag {
		return SelfClosingTag{
			tagName: tagName,
			attrs:   attrs,
		}
	}
}

type SelfClosingTag struct {
	tagName string
	attrs   []attrer
}

func (s SelfClosingTag) tag() {}

type Tag struct {
	tagName  string
	attrs    []attrer
	children []tagger
}

func (t Tag) tag() {}

func (t Tag) X(children ...tagger) Tag {
	clone := t.clone()
	clone.children = append(clone.children, children...)
	return clone
}

func (t Tag) clone() Tag {
	attrs := make([]attrer, len(t.attrs))
	copy(attrs, t.attrs)

	children := make([]tagger, len(t.children))
	copy(children, t.children)

	return Tag{
		tagName:  t.tagName,
		attrs:    attrs,
		children: children,
	}
}

func (t Tag) Render(w io.Writer) error {
	_, err := w.Write([]byte("<" + t.tagName))
	if err != nil {
		return err
	}

	if len(t.attrs) > 0 {
		w.Write([]byte(" "))
	}

	for idx, attr := range t.attrs {
		err = attr.Render(w)
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
		err = child.Render(w)
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
	children []tagger
}

func (HTML5Doctype) tag() {}
func (t HTML5Doctype) Render(w io.Writer) error {
	_, err := w.Write([]byte("<!DOCTYPE html>"))
	if err != nil {
		return err
	}
	for _, child := range t.children {
		err = child.Render(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t HTML5Doctype) X(children ...tagger) HTML5Doctype {
	clone := t.clone()
	clone.children = append(clone.children, children...)
	return clone
}

func (t HTML5Doctype) clone() HTML5Doctype {
	children := make([]tagger, len(t.children))
	copy(children, t.children)

	return HTML5Doctype{
		children: children,
	}
}

// All known HTML5 tags

func NewHTML5Doctype() HTML5Doctype { return HTML5Doctype{} }

func A(attrs ...attrer) Tag                 { return TagBuilder("a")(attrs...) }
func Abbr(attrs ...attrer) Tag              { return TagBuilder("abbr")(attrs...) }
func Address(attrs ...attrer) Tag           { return TagBuilder("address")(attrs...) }
func Area(attrs ...attrer) SelfClosingTag   { return SelfClosingTagBuilder("area")(attrs...) }
func Article(attrs ...attrer) Tag           { return TagBuilder("article")(attrs...) }
func Aside(attrs ...attrer) Tag             { return TagBuilder("aside")(attrs...) }
func Audio(attrs ...attrer) Tag             { return TagBuilder("audio")(attrs...) }
func B(attrs ...attrer) Tag                 { return TagBuilder("b")(attrs...) }
func Base(attrs ...attrer) SelfClosingTag   { return SelfClosingTagBuilder("base")(attrs...) }
func Bdi(attrs ...attrer) Tag               { return TagBuilder("bdi")(attrs...) }
func Bdo(attrs ...attrer) Tag               { return TagBuilder("bdo")(attrs...) }
func Blockquote(attrs ...attrer) Tag        { return TagBuilder("blockquote")(attrs...) }
func Body(attrs ...attrer) Tag              { return TagBuilder("body")(attrs...) }
func Br(attrs ...attrer) SelfClosingTag     { return SelfClosingTagBuilder("br")(attrs...) }
func Button(attrs ...attrer) Tag            { return TagBuilder("button")(attrs...) }
func Canvas(attrs ...attrer) Tag            { return TagBuilder("canvas")(attrs...) }
func Caption(attrs ...attrer) Tag           { return TagBuilder("caption")(attrs...) }
func Cite(attrs ...attrer) Tag              { return TagBuilder("cite")(attrs...) }
func Code(attrs ...attrer) Tag              { return TagBuilder("code")(attrs...) }
func Col(attrs ...attrer) SelfClosingTag    { return SelfClosingTagBuilder("col")(attrs...) }
func Colgroup(attrs ...attrer) Tag          { return TagBuilder("colgroup")(attrs...) }
func Data(attrs ...attrer) Tag              { return TagBuilder("data")(attrs...) }
func Datalist(attrs ...attrer) Tag          { return TagBuilder("datalist")(attrs...) }
func Dd(attrs ...attrer) Tag                { return TagBuilder("dd")(attrs...) }
func Del(attrs ...attrer) Tag               { return TagBuilder("del")(attrs...) }
func Details(attrs ...attrer) Tag           { return TagBuilder("details")(attrs...) }
func Dfn(attrs ...attrer) Tag               { return TagBuilder("dfn")(attrs...) }
func Dialog(attrs ...attrer) Tag            { return TagBuilder("dialog")(attrs...) }
func Div(attrs ...attrer) Tag               { return TagBuilder("div")(attrs...) }
func Dl(attrs ...attrer) Tag                { return TagBuilder("dl")(attrs...) }
func Dt(attrs ...attrer) Tag                { return TagBuilder("dt")(attrs...) }
func Em(attrs ...attrer) Tag                { return TagBuilder("em")(attrs...) }
func Embed(attrs ...attrer) SelfClosingTag  { return SelfClosingTagBuilder("embed")(attrs...) }
func Fieldset(attrs ...attrer) Tag          { return TagBuilder("fieldset")(attrs...) }
func Figcaption(attrs ...attrer) Tag        { return TagBuilder("figcaption")(attrs...) }
func Figure(attrs ...attrer) Tag            { return TagBuilder("figure")(attrs...) }
func Footer(attrs ...attrer) Tag            { return TagBuilder("footer")(attrs...) }
func Form(attrs ...attrer) Tag              { return TagBuilder("form")(attrs...) }
func H1(attrs ...attrer) Tag                { return TagBuilder("h1")(attrs...) }
func H2(attrs ...attrer) Tag                { return TagBuilder("h2")(attrs...) }
func H3(attrs ...attrer) Tag                { return TagBuilder("h3")(attrs...) }
func H4(attrs ...attrer) Tag                { return TagBuilder("h4")(attrs...) }
func H5(attrs ...attrer) Tag                { return TagBuilder("h5")(attrs...) }
func H6(attrs ...attrer) Tag                { return TagBuilder("h6")(attrs...) }
func Head(attrs ...attrer) Tag              { return TagBuilder("head")(attrs...) }
func Header(attrs ...attrer) Tag            { return TagBuilder("header")(attrs...) }
func Hr(attrs ...attrer) SelfClosingTag     { return SelfClosingTagBuilder("hr")(attrs...) }
func HTML(attrs ...attrer) Tag              { return TagBuilder("html")(attrs...) }
func I(attrs ...attrer) Tag                 { return TagBuilder("i")(attrs...) }
func Iframe(attrs ...attrer) Tag            { return TagBuilder("iframe")(attrs...) }
func Img(attrs ...attrer) SelfClosingTag    { return SelfClosingTagBuilder("img")(attrs...) }
func Input(attrs ...attrer) SelfClosingTag  { return SelfClosingTagBuilder("input")(attrs...) }
func Ins(attrs ...attrer) Tag               { return TagBuilder("ins")(attrs...) }
func Kbd(attrs ...attrer) Tag               { return TagBuilder("kbd")(attrs...) }
func Label(attrs ...attrer) Tag             { return TagBuilder("label")(attrs...) }
func Legend(attrs ...attrer) Tag            { return TagBuilder("legend")(attrs...) }
func Li(attrs ...attrer) Tag                { return TagBuilder("li")(attrs...) }
func Link(attrs ...attrer) SelfClosingTag   { return SelfClosingTagBuilder("link")(attrs...) }
func Main(attrs ...attrer) Tag              { return TagBuilder("main")(attrs...) }
func Map(attrs ...attrer) Tag               { return TagBuilder("map")(attrs...) }
func Mark(attrs ...attrer) Tag              { return TagBuilder("mark")(attrs...) }
func Meta(attrs ...attrer) SelfClosingTag   { return SelfClosingTagBuilder("meta")(attrs...) }
func Meter(attrs ...attrer) Tag             { return TagBuilder("meter")(attrs...) }
func Nav(attrs ...attrer) Tag               { return TagBuilder("nav")(attrs...) }
func Noscript(attrs ...attrer) Tag          { return TagBuilder("noscript")(attrs...) }
func Object(attrs ...attrer) Tag            { return TagBuilder("object")(attrs...) }
func Ol(attrs ...attrer) Tag                { return TagBuilder("ol")(attrs...) }
func Optgroup(attrs ...attrer) Tag          { return TagBuilder("optgroup")(attrs...) }
func Option(attrs ...attrer) Tag            { return TagBuilder("option")(attrs...) }
func Output(attrs ...attrer) Tag            { return TagBuilder("output")(attrs...) }
func P(attrs ...attrer) Tag                 { return TagBuilder("p")(attrs...) }
func Param(attrs ...attrer) SelfClosingTag  { return SelfClosingTagBuilder("param")(attrs...) }
func Picture(attrs ...attrer) Tag           { return TagBuilder("picture")(attrs...) }
func Pre(attrs ...attrer) Tag               { return TagBuilder("pre")(attrs...) }
func Progress(attrs ...attrer) Tag          { return TagBuilder("progress")(attrs...) }
func Q(attrs ...attrer) Tag                 { return TagBuilder("q")(attrs...) }
func Rp(attrs ...attrer) Tag                { return TagBuilder("rp")(attrs...) }
func Rt(attrs ...attrer) Tag                { return TagBuilder("rt")(attrs...) }
func Ruby(attrs ...attrer) Tag              { return TagBuilder("ruby")(attrs...) }
func S(attrs ...attrer) Tag                 { return TagBuilder("s")(attrs...) }
func Samp(attrs ...attrer) Tag              { return TagBuilder("samp")(attrs...) }
func Script(attrs ...attrer) Tag            { return TagBuilder("script")(attrs...) }
func Section(attrs ...attrer) Tag           { return TagBuilder("section")(attrs...) }
func Select(attrs ...attrer) Tag            { return TagBuilder("select")(attrs...) }
func Slot(attrs ...attrer) Tag              { return TagBuilder("slot")(attrs...) }
func Small(attrs ...attrer) Tag             { return TagBuilder("small")(attrs...) }
func Source(attrs ...attrer) SelfClosingTag { return SelfClosingTagBuilder("source")(attrs...) }
func Span(attrs ...attrer) Tag              { return TagBuilder("span")(attrs...) }
func Strong(attrs ...attrer) Tag            { return TagBuilder("strong")(attrs...) }
func Style(attrs ...attrer) Tag             { return TagBuilder("style")(attrs...) }
func Sub(attrs ...attrer) Tag               { return TagBuilder("sub")(attrs...) }
func Summary(attrs ...attrer) Tag           { return TagBuilder("summary")(attrs...) }
func Sup(attrs ...attrer) Tag               { return TagBuilder("sup")(attrs...) }
func Table(attrs ...attrer) Tag             { return TagBuilder("table")(attrs...) }
func Tbody(attrs ...attrer) Tag             { return TagBuilder("tbody")(attrs...) }
func Td(attrs ...attrer) Tag                { return TagBuilder("td")(attrs...) }
func Template(attrs ...attrer) Tag          { return TagBuilder("template")(attrs...) }
func Textarea(attrs ...attrer) Tag          { return TagBuilder("textarea")(attrs...) }
func Tfoot(attrs ...attrer) Tag             { return TagBuilder("tfoot")(attrs...) }
func Th(attrs ...attrer) Tag                { return TagBuilder("th")(attrs...) }
func Thead(attrs ...attrer) Tag             { return TagBuilder("thead")(attrs...) }
func Time(attrs ...attrer) Tag              { return TagBuilder("time")(attrs...) }
func Title(attrs ...attrer) Tag             { return TagBuilder("title")(attrs...) }
func Tr(attrs ...attrer) Tag                { return TagBuilder("tr")(attrs...) }
func Track(attrs ...attrer) SelfClosingTag  { return SelfClosingTagBuilder("track")(attrs...) }
func U(attrs ...attrer) Tag                 { return TagBuilder("u")(attrs...) }
func Ul(attrs ...attrer) Tag                { return TagBuilder("ul")(attrs...) }
func Var(attrs ...attrer) Tag               { return TagBuilder("var")(attrs...) }
func Video(attrs ...attrer) Tag             { return TagBuilder("video")(attrs...) }
func Wbr(attrs ...attrer) SelfClosingTag    { return SelfClosingTagBuilder("wbr")(attrs...) }
