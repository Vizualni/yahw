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

func unwrapNodes(nodes []Node) []Node {
	nn := []Node{}
	for _, n := range nodes {
		switch t := n.(type) {
		case Nodes:
			nn = append(nn, unwrapNodes(t)...)
		default:
			nn = append(nn, n)
		}
	}
	return nn
}

func TagBuilder(tagName string) func(...Node) CommonTag {
	if !isValidTagName(tagName) {
		panic("Invalid tag name: " + tagName)
	}

	return func(nodes ...Node) CommonTag {
		ct := CommonTag{
			tagName: tagName,
		}

		unwrapped := unwrapNodes(nodes)
		for _, n := range unwrapped {
			if n == nil {
				continue
			}

			switch t := n.(type) {
			case attrable:
				ct.attrs = append(ct.attrs, t)
			case taggable:
				ct.tags = append(ct.tags, t)
			default:
				r := n.Node()
				switch t := r.(type) {
				case attrable:
					ct.attrs = append(ct.attrs, t)
				case taggable:
					ct.tags = append(ct.tags, t)
				default:
					panic("invalid node type")
				}
			}
		}

		return ct
	}
}

func SelfClosingTagBuilder(tagName string) func(...attrable) SelfClosingTag {
	if !isValidTagName(tagName) {
		panic("Invalid self closing tag name: " + tagName)
	}
	return func(attrs ...attrable) SelfClosingTag {
		return SelfClosingTag{
			tagName: tagName,
			attrs:   attrs,
		}
	}
}

type SelfClosingTag struct {
	tagName string
	attrs   []attrable
}

func (t SelfClosingTag) tag()             {}
func (t SelfClosingTag) Node() Renderable { return t }

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
		err = attr.Render(w)
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

type Nodes []Node

func (n Nodes) Node() Renderable { panic("Nodes is not a node") }

type CommonTag struct {
	tagName string
	attrs   []attrable
	tags    []taggable
}

func (t CommonTag) tag()             {}
func (t CommonTag) Node() Renderable { return t }

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

	for _, child := range t.tags {
		if child == nil {
			continue
		}
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
	children TagSlice
}

func (t HTML5Doctype) tag()             {}
func (t HTML5Doctype) Node() Renderable { return t }

func (t HTML5Doctype) Render(w io.Writer) error {
	_, err := w.Write([]byte("<!DOCTYPE html>"))
	if err != nil {
		return err
	}
	return t.children.Render(w)
}

type TagSlice []taggable

func (t TagSlice) tag() {}
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
func (t TagSlice) Node() Renderable { return t }

// All known HTML5 tags

func NewHTML5Doctype(cs ...taggable) HTML5Doctype { return HTML5Doctype{children: cs} }
func A(attrs ...Node) CommonTag                   { return TagBuilder("a")(attrs...) }
func Abbr(attrs ...Node) CommonTag                { return TagBuilder("abbr")(attrs...) }
func Address(attrs ...Node) CommonTag             { return TagBuilder("address")(attrs...) }
func Area(attrs ...attrable) SelfClosingTag       { return SelfClosingTagBuilder("area")(attrs...) }
func Article(attrs ...Node) CommonTag             { return TagBuilder("article")(attrs...) }
func Aside(attrs ...Node) CommonTag               { return TagBuilder("aside")(attrs...) }
func Audio(attrs ...Node) CommonTag               { return TagBuilder("audio")(attrs...) }
func B(attrs ...Node) CommonTag                   { return TagBuilder("b")(attrs...) }
func Base(attrs ...attrable) SelfClosingTag       { return SelfClosingTagBuilder("base")(attrs...) }
func Bdi(attrs ...Node) CommonTag                 { return TagBuilder("bdi")(attrs...) }
func Bdo(attrs ...Node) CommonTag                 { return TagBuilder("bdo")(attrs...) }
func Blockquote(attrs ...Node) CommonTag          { return TagBuilder("blockquote")(attrs...) }
func Body(attrs ...Node) CommonTag                { return TagBuilder("body")(attrs...) }
func Br(attrs ...attrable) SelfClosingTag         { return SelfClosingTagBuilder("br")(attrs...) }
func Button(attrs ...Node) CommonTag              { return TagBuilder("button")(attrs...) }
func Canvas(attrs ...Node) CommonTag              { return TagBuilder("canvas")(attrs...) }
func Caption(attrs ...Node) CommonTag             { return TagBuilder("caption")(attrs...) }
func Cite(attrs ...Node) CommonTag                { return TagBuilder("cite")(attrs...) }
func Code(attrs ...Node) CommonTag                { return TagBuilder("code")(attrs...) }
func Col(attrs ...attrable) SelfClosingTag        { return SelfClosingTagBuilder("col")(attrs...) }
func Colgroup(attrs ...Node) CommonTag            { return TagBuilder("colgroup")(attrs...) }
func Data(attrs ...Node) CommonTag                { return TagBuilder("data")(attrs...) }
func Datalist(attrs ...Node) CommonTag            { return TagBuilder("datalist")(attrs...) }
func Dd(attrs ...Node) CommonTag                  { return TagBuilder("dd")(attrs...) }
func Del(attrs ...Node) CommonTag                 { return TagBuilder("del")(attrs...) }
func Details(attrs ...Node) CommonTag             { return TagBuilder("details")(attrs...) }
func Dfn(attrs ...Node) CommonTag                 { return TagBuilder("dfn")(attrs...) }
func Dialog(attrs ...Node) CommonTag              { return TagBuilder("dialog")(attrs...) }
func Div(attrs ...Node) CommonTag                 { return TagBuilder("div")(attrs...) }
func Dl(attrs ...Node) CommonTag                  { return TagBuilder("dl")(attrs...) }
func Dt(attrs ...Node) CommonTag                  { return TagBuilder("dt")(attrs...) }
func Em(attrs ...Node) CommonTag                  { return TagBuilder("em")(attrs...) }
func Embed(attrs ...attrable) SelfClosingTag      { return SelfClosingTagBuilder("embed")(attrs...) }
func Fieldset(attrs ...Node) CommonTag            { return TagBuilder("fieldset")(attrs...) }
func Figcaption(attrs ...Node) CommonTag          { return TagBuilder("figcaption")(attrs...) }
func Figure(attrs ...Node) CommonTag              { return TagBuilder("figure")(attrs...) }
func Footer(attrs ...Node) CommonTag              { return TagBuilder("footer")(attrs...) }
func Form(attrs ...Node) CommonTag                { return TagBuilder("form")(attrs...) }
func H1(attrs ...Node) CommonTag                  { return TagBuilder("h1")(attrs...) }
func H2(attrs ...Node) CommonTag                  { return TagBuilder("h2")(attrs...) }
func H3(attrs ...Node) CommonTag                  { return TagBuilder("h3")(attrs...) }
func H4(attrs ...Node) CommonTag                  { return TagBuilder("h4")(attrs...) }
func H5(attrs ...Node) CommonTag                  { return TagBuilder("h5")(attrs...) }
func H6(attrs ...Node) CommonTag                  { return TagBuilder("h6")(attrs...) }
func Head(attrs ...Node) CommonTag                { return TagBuilder("head")(attrs...) }
func Header(attrs ...Node) CommonTag              { return TagBuilder("header")(attrs...) }
func Hr(attrs ...attrable) SelfClosingTag         { return SelfClosingTagBuilder("hr")(attrs...) }
func HTML(attrs ...Node) CommonTag                { return TagBuilder("html")(attrs...) }
func I(attrs ...Node) CommonTag                   { return TagBuilder("i")(attrs...) }
func Iframe(attrs ...Node) CommonTag              { return TagBuilder("iframe")(attrs...) }
func Img(attrs ...attrable) SelfClosingTag        { return SelfClosingTagBuilder("img")(attrs...) }
func Input(attrs ...attrable) SelfClosingTag      { return SelfClosingTagBuilder("input")(attrs...) }
func Ins(attrs ...Node) CommonTag                 { return TagBuilder("ins")(attrs...) }
func Kbd(attrs ...Node) CommonTag                 { return TagBuilder("kbd")(attrs...) }
func Label(attrs ...Node) CommonTag               { return TagBuilder("label")(attrs...) }
func Legend(attrs ...Node) CommonTag              { return TagBuilder("legend")(attrs...) }
func Li(attrs ...Node) CommonTag                  { return TagBuilder("li")(attrs...) }
func Link(attrs ...attrable) SelfClosingTag       { return SelfClosingTagBuilder("link")(attrs...) }
func Main(attrs ...Node) CommonTag                { return TagBuilder("main")(attrs...) }
func Map(attrs ...Node) CommonTag                 { return TagBuilder("map")(attrs...) }
func Mark(attrs ...Node) CommonTag                { return TagBuilder("mark")(attrs...) }
func Meta(attrs ...attrable) SelfClosingTag       { return SelfClosingTagBuilder("meta")(attrs...) }
func Meter(attrs ...Node) CommonTag               { return TagBuilder("meter")(attrs...) }
func Nav(attrs ...Node) CommonTag                 { return TagBuilder("nav")(attrs...) }
func Noscript(attrs ...Node) CommonTag            { return TagBuilder("noscript")(attrs...) }
func Object(attrs ...Node) CommonTag              { return TagBuilder("object")(attrs...) }
func Ol(attrs ...Node) CommonTag                  { return TagBuilder("ol")(attrs...) }
func Optgroup(attrs ...Node) CommonTag            { return TagBuilder("optgroup")(attrs...) }
func Option(attrs ...Node) CommonTag              { return TagBuilder("option")(attrs...) }
func Output(attrs ...Node) CommonTag              { return TagBuilder("output")(attrs...) }
func P(attrs ...Node) CommonTag                   { return TagBuilder("p")(attrs...) }
func Param(attrs ...attrable) SelfClosingTag      { return SelfClosingTagBuilder("param")(attrs...) }
func Picture(attrs ...Node) CommonTag             { return TagBuilder("picture")(attrs...) }
func Pre(attrs ...Node) CommonTag                 { return TagBuilder("pre")(attrs...) }
func Progress(attrs ...Node) CommonTag            { return TagBuilder("progress")(attrs...) }
func Q(attrs ...Node) CommonTag                   { return TagBuilder("q")(attrs...) }
func Rp(attrs ...Node) CommonTag                  { return TagBuilder("rp")(attrs...) }
func Rt(attrs ...Node) CommonTag                  { return TagBuilder("rt")(attrs...) }
func Ruby(attrs ...Node) CommonTag                { return TagBuilder("ruby")(attrs...) }
func S(attrs ...Node) CommonTag                   { return TagBuilder("s")(attrs...) }
func Samp(attrs ...Node) CommonTag                { return TagBuilder("samp")(attrs...) }
func Script(attrs ...Node) CommonTag              { return TagBuilder("script")(attrs...) }
func Section(attrs ...Node) CommonTag             { return TagBuilder("section")(attrs...) }
func Select(attrs ...Node) CommonTag              { return TagBuilder("select")(attrs...) }
func Slot(attrs ...Node) CommonTag                { return TagBuilder("slot")(attrs...) }
func Small(attrs ...Node) CommonTag               { return TagBuilder("small")(attrs...) }
func Source(attrs ...attrable) SelfClosingTag     { return SelfClosingTagBuilder("source")(attrs...) }
func Span(attrs ...Node) CommonTag                { return TagBuilder("span")(attrs...) }
func Strong(attrs ...Node) CommonTag              { return TagBuilder("strong")(attrs...) }
func Style(attrs ...Node) CommonTag               { return TagBuilder("style")(attrs...) }
func Sub(attrs ...Node) CommonTag                 { return TagBuilder("sub")(attrs...) }
func Summary(attrs ...Node) CommonTag             { return TagBuilder("summary")(attrs...) }
func Sup(attrs ...Node) CommonTag                 { return TagBuilder("sup")(attrs...) }
func Table(attrs ...Node) CommonTag               { return TagBuilder("table")(attrs...) }
func Tbody(attrs ...Node) CommonTag               { return TagBuilder("tbody")(attrs...) }
func Td(attrs ...Node) CommonTag                  { return TagBuilder("td")(attrs...) }
func Template(attrs ...Node) CommonTag            { return TagBuilder("template")(attrs...) }
func Textarea(attrs ...Node) CommonTag            { return TagBuilder("textarea")(attrs...) }
func Tfoot(attrs ...Node) CommonTag               { return TagBuilder("tfoot")(attrs...) }
func Th(attrs ...Node) CommonTag                  { return TagBuilder("th")(attrs...) }
func Thead(attrs ...Node) CommonTag               { return TagBuilder("thead")(attrs...) }
func Time(attrs ...Node) CommonTag                { return TagBuilder("time")(attrs...) }
func Title(attrs ...Node) CommonTag               { return TagBuilder("title")(attrs...) }
func Tr(attrs ...Node) CommonTag                  { return TagBuilder("tr")(attrs...) }
func Track(attrs ...attrable) SelfClosingTag      { return SelfClosingTagBuilder("track")(attrs...) }
func U(attrs ...Node) CommonTag                   { return TagBuilder("u")(attrs...) }
func Ul(attrs ...Node) CommonTag                  { return TagBuilder("ul")(attrs...) }
func Var(attrs ...Node) CommonTag                 { return TagBuilder("var")(attrs...) }
func Video(attrs ...Node) CommonTag               { return TagBuilder("video")(attrs...) }
func Wbr(attrs ...attrable) SelfClosingTag        { return SelfClosingTagBuilder("wbr")(attrs...) }
