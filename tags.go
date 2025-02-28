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
		case c == '-' || c == '_':
		default:
			return false
		}
	}
	return true
}

func NewTag(tagName string) CommonTag {
	return CommonTag{
		tagName: tagName,
	}
}

func TagBuilder(tagName string) func(...Attr) CommonTag {
	if !isValidTagName(tagName) {
		panic("Invalid tag name: " + tagName)
	}

	return func(attrs ...Attr) CommonTag {
		return CommonTag{
			tagName: tagName,
			attrs:   attrs,
		}
	}
}

func (t CommonTag) Attrs(attrs ...Attr) CommonTag {
	clone := t.clone()
	clone.attrs = append(clone.attrs, attrs...)
	return clone
}

func SelfClosingTagBuilder(tagName string) func(...Attr) SelfClosingTag {
	if !isValidTagName(tagName) {
		panic("Invalid self closing tag name: " + tagName)
	}
	return func(attrs ...Attr) SelfClosingTag {
		return SelfClosingTag{
			tagName: tagName,
			attrs:   attrs,
		}
	}
}

type SelfClosingTag struct {
	tagName string
	attrs   []Attr
}

func (t SelfClosingTag) Tag() Renderable { return t }

func mergeClasses(clss AttrSlice) Classes {
	merged := Classes("")
	for _, c := range clss {
		switch t := c.(type) {
		case Classes:
			merged = merged.Merge(t)
		case ClassesMap:
			merged = merged.MergeMap(t)
		case Attribute:
			merged = merged.Add(t.value)
		default:
			panic("can only merge attributes with class")
		}
	}

	return merged
}

func (t SelfClosingTag) Render(w io.Writer) error {
	_, err := w.Write([]byte("<" + t.tagName))
	if err != nil {
		return err
	}

	if len(t.attrs) > 0 {
		w.Write([]byte(" "))
	}

	newAttrs := AttrSlice{}
	toMerge := map[string]AttrSlice{}
	for _, attr := range t.attrs {
		switch t := attr.(type) {
		case Classes, ClassesMap:
			toMerge["class"] = append(toMerge["class"], t)
		case Attribute:
			if t.key == "class" {
				toMerge["class"] = append(toMerge["class"], t)
			} else {
				newAttrs = append(newAttrs, t)
			}
		default:
			newAttrs = append(newAttrs, t)
		}
	}

	for k, attrs := range toMerge {
		switch k {
		case "class":
			newAttrs = append(newAttrs, mergeClasses(attrs))
		}
	}

	for idx, attr := range newAttrs {
		if attr == nil {
			continue
		}
		err = attr.Attr().Render(w)
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

type CommonTag struct {
	tagName  string
	attrs    []Attr
	children []Tag
}

func (t CommonTag) Tag() Renderable { return t }

func (t CommonTag) X(children ...Tag) CommonTag {
	clone := t.clone()
	clone.children = append(clone.children, children...)
	return clone
}

func (t CommonTag) clone() CommonTag {
	attrs := make([]Attr, len(t.attrs))
	copy(attrs, t.attrs)

	children := make([]Tag, len(t.children))
	copy(children, t.children)

	return CommonTag{
		tagName:  t.tagName,
		attrs:    attrs,
		children: children,
	}
}

func (t CommonTag) Render(w io.Writer) error {
	_, err := w.Write([]byte("<" + t.tagName))
	if err != nil {
		return err
	}

	if len(t.attrs) > 0 {
		w.Write([]byte(" "))
	}
	newAttrs := AttrSlice{}
	toMerge := map[string]AttrSlice{}
	for _, attr := range t.attrs {
		switch t := attr.(type) {
		case Classes, ClassesMap:
			toMerge["class"] = append(toMerge["class"], t)
		case Attribute:
			if t.key == "class" {
				toMerge["class"] = append(toMerge["class"], t)
			} else {
				newAttrs = append(newAttrs, t)
			}
		default:
			newAttrs = append(newAttrs, t)
		}
	}

	for k, attrs := range toMerge {
		switch k {
		case "class":
			newAttrs = append(newAttrs, mergeClasses(attrs))
		}
	}

	for idx, attr := range newAttrs {
		if attr == nil {
			continue
		}
		err = attr.Attr().Render(w)
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
		if child == nil {
			continue
		}
		err = child.Tag().Render(w)
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
	children TagSlice
}

func (t HTML5Doctype) Tag() Renderable { return t }

func (t HTML5Doctype) Render(w io.Writer) error {
	_, err := w.Write([]byte("<!DOCTYPE html>"))
	if err != nil {
		return err
	}
	return t.children.Render(w)
}

func (t HTML5Doctype) X(children ...Renderable) HTML5Doctype {
	clone := t.clone()
	clone.children = append(clone.children, children...)
	return clone
}

func (t HTML5Doctype) clone() HTML5Doctype {
	children := make([]Renderable, len(t.children))
	copy(children, t.children)

	return HTML5Doctype{
		children: children,
	}
}

type TagSlice []Renderable

func (t TagSlice) Render(w io.Writer) error {
	for _, tag := range t {
		if tag == nil {
			continue
		}
		err := tag.Render(w)
		if err != nil {
			return err
		}
	}
	return nil
}

// All known HTML5 tags

func NewHTML5Doctype(cs ...Renderable) HTML5Doctype { return HTML5Doctype{children: cs} }

func A(attrs ...Attr) CommonTag           { return TagBuilder("a")(attrs...) }
func Abbr(attrs ...Attr) CommonTag        { return TagBuilder("abbr")(attrs...) }
func Address(attrs ...Attr) CommonTag     { return TagBuilder("address")(attrs...) }
func Area(attrs ...Attr) SelfClosingTag   { return SelfClosingTagBuilder("area")(attrs...) }
func Article(attrs ...Attr) CommonTag     { return TagBuilder("article")(attrs...) }
func Aside(attrs ...Attr) CommonTag       { return TagBuilder("aside")(attrs...) }
func Audio(attrs ...Attr) CommonTag       { return TagBuilder("audio")(attrs...) }
func B(attrs ...Attr) CommonTag           { return TagBuilder("b")(attrs...) }
func Base(attrs ...Attr) SelfClosingTag   { return SelfClosingTagBuilder("base")(attrs...) }
func Bdi(attrs ...Attr) CommonTag         { return TagBuilder("bdi")(attrs...) }
func Bdo(attrs ...Attr) CommonTag         { return TagBuilder("bdo")(attrs...) }
func Blockquote(attrs ...Attr) CommonTag  { return TagBuilder("blockquote")(attrs...) }
func Body(attrs ...Attr) CommonTag        { return TagBuilder("body")(attrs...) }
func Br(attrs ...Attr) SelfClosingTag     { return SelfClosingTagBuilder("br")(attrs...) }
func Button(attrs ...Attr) CommonTag      { return TagBuilder("button")(attrs...) }
func Canvas(attrs ...Attr) CommonTag      { return TagBuilder("canvas")(attrs...) }
func Caption(attrs ...Attr) CommonTag     { return TagBuilder("caption")(attrs...) }
func Cite(attrs ...Attr) CommonTag        { return TagBuilder("cite")(attrs...) }
func Code(attrs ...Attr) CommonTag        { return TagBuilder("code")(attrs...) }
func Col(attrs ...Attr) SelfClosingTag    { return SelfClosingTagBuilder("col")(attrs...) }
func Colgroup(attrs ...Attr) CommonTag    { return TagBuilder("colgroup")(attrs...) }
func Data(attrs ...Attr) CommonTag        { return TagBuilder("data")(attrs...) }
func Datalist(attrs ...Attr) CommonTag    { return TagBuilder("datalist")(attrs...) }
func Dd(attrs ...Attr) CommonTag          { return TagBuilder("dd")(attrs...) }
func Del(attrs ...Attr) CommonTag         { return TagBuilder("del")(attrs...) }
func Details(attrs ...Attr) CommonTag     { return TagBuilder("details")(attrs...) }
func Dfn(attrs ...Attr) CommonTag         { return TagBuilder("dfn")(attrs...) }
func Dialog(attrs ...Attr) CommonTag      { return TagBuilder("dialog")(attrs...) }
func Div(attrs ...Attr) CommonTag         { return TagBuilder("div")(attrs...) }
func Dl(attrs ...Attr) CommonTag          { return TagBuilder("dl")(attrs...) }
func Dt(attrs ...Attr) CommonTag          { return TagBuilder("dt")(attrs...) }
func Em(attrs ...Attr) CommonTag          { return TagBuilder("em")(attrs...) }
func Embed(attrs ...Attr) SelfClosingTag  { return SelfClosingTagBuilder("embed")(attrs...) }
func Fieldset(attrs ...Attr) CommonTag    { return TagBuilder("fieldset")(attrs...) }
func Figcaption(attrs ...Attr) CommonTag  { return TagBuilder("figcaption")(attrs...) }
func Figure(attrs ...Attr) CommonTag      { return TagBuilder("figure")(attrs...) }
func Footer(attrs ...Attr) CommonTag      { return TagBuilder("footer")(attrs...) }
func Form(attrs ...Attr) CommonTag        { return TagBuilder("form")(attrs...) }
func H1(attrs ...Attr) CommonTag          { return TagBuilder("h1")(attrs...) }
func H2(attrs ...Attr) CommonTag          { return TagBuilder("h2")(attrs...) }
func H3(attrs ...Attr) CommonTag          { return TagBuilder("h3")(attrs...) }
func H4(attrs ...Attr) CommonTag          { return TagBuilder("h4")(attrs...) }
func H5(attrs ...Attr) CommonTag          { return TagBuilder("h5")(attrs...) }
func H6(attrs ...Attr) CommonTag          { return TagBuilder("h6")(attrs...) }
func Head(attrs ...Attr) CommonTag        { return TagBuilder("head")(attrs...) }
func Header(attrs ...Attr) CommonTag      { return TagBuilder("header")(attrs...) }
func Hr(attrs ...Attr) SelfClosingTag     { return SelfClosingTagBuilder("hr")(attrs...) }
func HTML(attrs ...Attr) CommonTag        { return TagBuilder("html")(attrs...) }
func I(attrs ...Attr) CommonTag           { return TagBuilder("i")(attrs...) }
func Iframe(attrs ...Attr) CommonTag      { return TagBuilder("iframe")(attrs...) }
func Img(attrs ...Attr) SelfClosingTag    { return SelfClosingTagBuilder("img")(attrs...) }
func Input(attrs ...Attr) SelfClosingTag  { return SelfClosingTagBuilder("input")(attrs...) }
func Ins(attrs ...Attr) CommonTag         { return TagBuilder("ins")(attrs...) }
func Kbd(attrs ...Attr) CommonTag         { return TagBuilder("kbd")(attrs...) }
func Label(attrs ...Attr) CommonTag       { return TagBuilder("label")(attrs...) }
func Legend(attrs ...Attr) CommonTag      { return TagBuilder("legend")(attrs...) }
func Li(attrs ...Attr) CommonTag          { return TagBuilder("li")(attrs...) }
func Link(attrs ...Attr) SelfClosingTag   { return SelfClosingTagBuilder("link")(attrs...) }
func Main(attrs ...Attr) CommonTag        { return TagBuilder("main")(attrs...) }
func Map(attrs ...Attr) CommonTag         { return TagBuilder("map")(attrs...) }
func Mark(attrs ...Attr) CommonTag        { return TagBuilder("mark")(attrs...) }
func Meta(attrs ...Attr) SelfClosingTag   { return SelfClosingTagBuilder("meta")(attrs...) }
func Meter(attrs ...Attr) CommonTag       { return TagBuilder("meter")(attrs...) }
func Nav(attrs ...Attr) CommonTag         { return TagBuilder("nav")(attrs...) }
func Noscript(attrs ...Attr) CommonTag    { return TagBuilder("noscript")(attrs...) }
func Object(attrs ...Attr) CommonTag      { return TagBuilder("object")(attrs...) }
func Ol(attrs ...Attr) CommonTag          { return TagBuilder("ol")(attrs...) }
func Optgroup(attrs ...Attr) CommonTag    { return TagBuilder("optgroup")(attrs...) }
func Option(attrs ...Attr) CommonTag      { return TagBuilder("option")(attrs...) }
func Output(attrs ...Attr) CommonTag      { return TagBuilder("output")(attrs...) }
func P(attrs ...Attr) CommonTag           { return TagBuilder("p")(attrs...) }
func Param(attrs ...Attr) SelfClosingTag  { return SelfClosingTagBuilder("param")(attrs...) }
func Picture(attrs ...Attr) CommonTag     { return TagBuilder("picture")(attrs...) }
func Pre(attrs ...Attr) CommonTag         { return TagBuilder("pre")(attrs...) }
func Progress(attrs ...Attr) CommonTag    { return TagBuilder("progress")(attrs...) }
func Q(attrs ...Attr) CommonTag           { return TagBuilder("q")(attrs...) }
func Rp(attrs ...Attr) CommonTag          { return TagBuilder("rp")(attrs...) }
func Rt(attrs ...Attr) CommonTag          { return TagBuilder("rt")(attrs...) }
func Ruby(attrs ...Attr) CommonTag        { return TagBuilder("ruby")(attrs...) }
func S(attrs ...Attr) CommonTag           { return TagBuilder("s")(attrs...) }
func Samp(attrs ...Attr) CommonTag        { return TagBuilder("samp")(attrs...) }
func Script(attrs ...Attr) CommonTag      { return TagBuilder("script")(attrs...) }
func Section(attrs ...Attr) CommonTag     { return TagBuilder("section")(attrs...) }
func Select(attrs ...Attr) CommonTag      { return TagBuilder("select")(attrs...) }
func Slot(attrs ...Attr) CommonTag        { return TagBuilder("slot")(attrs...) }
func Small(attrs ...Attr) CommonTag       { return TagBuilder("small")(attrs...) }
func Source(attrs ...Attr) SelfClosingTag { return SelfClosingTagBuilder("source")(attrs...) }
func Span(attrs ...Attr) CommonTag        { return TagBuilder("span")(attrs...) }
func Strong(attrs ...Attr) CommonTag      { return TagBuilder("strong")(attrs...) }
func Style(attrs ...Attr) CommonTag       { return TagBuilder("style")(attrs...) }
func Sub(attrs ...Attr) CommonTag         { return TagBuilder("sub")(attrs...) }
func Summary(attrs ...Attr) CommonTag     { return TagBuilder("summary")(attrs...) }
func Sup(attrs ...Attr) CommonTag         { return TagBuilder("sup")(attrs...) }
func Table(attrs ...Attr) CommonTag       { return TagBuilder("table")(attrs...) }
func Tbody(attrs ...Attr) CommonTag       { return TagBuilder("tbody")(attrs...) }
func Td(attrs ...Attr) CommonTag          { return TagBuilder("td")(attrs...) }
func Template(attrs ...Attr) CommonTag    { return TagBuilder("template")(attrs...) }
func Textarea(attrs ...Attr) CommonTag    { return TagBuilder("textarea")(attrs...) }
func Tfoot(attrs ...Attr) CommonTag       { return TagBuilder("tfoot")(attrs...) }
func Th(attrs ...Attr) CommonTag          { return TagBuilder("th")(attrs...) }
func Thead(attrs ...Attr) CommonTag       { return TagBuilder("thead")(attrs...) }
func Time(attrs ...Attr) CommonTag        { return TagBuilder("time")(attrs...) }
func Title(attrs ...Attr) CommonTag       { return TagBuilder("title")(attrs...) }
func Tr(attrs ...Attr) CommonTag          { return TagBuilder("tr")(attrs...) }
func Track(attrs ...Attr) SelfClosingTag  { return SelfClosingTagBuilder("track")(attrs...) }
func U(attrs ...Attr) CommonTag           { return TagBuilder("u")(attrs...) }
func Ul(attrs ...Attr) CommonTag          { return TagBuilder("ul")(attrs...) }
func Var(attrs ...Attr) CommonTag         { return TagBuilder("var")(attrs...) }
func Video(attrs ...Attr) CommonTag       { return TagBuilder("video")(attrs...) }
func Wbr(attrs ...Attr) SelfClosingTag    { return SelfClosingTagBuilder("wbr")(attrs...) }
