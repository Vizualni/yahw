package yahw

import (
	"html"
	"io"
	"maps"
	"strings"
)

func isValidAttrName(name string) bool {
	if len(name) == 0 {
		return false
	}
	for _, c := range name {
		switch {
		case 'a' <= c && c <= 'z':
		case 'A' <= c && c <= 'Z':
		case '0' <= c && c <= '9':
		case strings.ContainsRune("_.-:", c):
		default:
			return false
		}
	}
	return true
}

func BuildAttr(key, value string) Attribute {
	if !isValidAttrName(key) {
		panic("Invalid attribute name: " + key)
	}
	return Attribute{
		key:   key,
		value: value,
	}
}

func NoValAttr(key string) NoValAttribute {
	if !isValidAttrName(key) {
		panic("Invalid attribute name: " + key)
	}
	return NoValAttribute{
		key: key,
	}
}

func AttrBuilder(key string) func(string) Attribute {
	return func(value string) Attribute {
		return BuildAttr(key, value)
	}
}

type Attribute struct {
	key   string
	value string
}

func (a Attribute) attr()            {}
func (a Attribute) Node() Renderable { return a }

func (a Attribute) Render(w io.Writer) error {
	escapedKey := html.EscapeString(a.key)
	escapedValue := html.EscapeString(a.value)
	_, err := w.Write([]byte(escapedKey + `="` + escapedValue + `"`))
	if err != nil {
		return err
	}
	return nil
}

type NoValAttribute struct {
	key string
}

func (a NoValAttribute) attr()            {}
func (a NoValAttribute) Node() Renderable { return a }

func (a NoValAttribute) Render(w io.Writer) error {
	escapedKey := html.EscapeString(a.key)
	_, err := w.Write([]byte(escapedKey))
	if err != nil {
		return err
	}
	return nil
}

type AttrSlice []attrable

func (a AttrSlice) attr()            {}
func (a AttrSlice) Node() Renderable { return a }

func (a AttrSlice) Render(w io.Writer) error {
	for i, attr := range a {
		if i > 0 {
			_, err := w.Write([]byte(" "))
			if err != nil {
				return err
			}
		}
		if attr == nil {
			continue
		}
		err := attr.Render(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a AttrSlice) Add(i attrable) AttrSlice {
	return append(a[:], i)
}

func (a AttrSlice) Merge(mrg AttrSlice) AttrSlice {
	return append(a[:], mrg...)
}

func extractClasses(cls string) []string {
	clss := strings.Split(string(cls), " ")
	res := make([]string, 0, len(clss))
	for _, cls := range clss {
		cls = strings.TrimSpace(cls)
		if cls == "" {
			continue
		}

		// Check for duplicates.
		// Not the most efficient way, but it's good enough for now.
		found := false
		for _, existing := range res {
			if existing == cls {
				found = true
				break
			}
		}
		if found {
			continue
		}
		res = append(res, cls)
	}

	return res
}

type Classes string

func (c Classes) attr()            {}
func (c Classes) Node() Renderable { return c }

func (c Classes) Render(w io.Writer) error {
	res := extractClasses(string(c))
	_, err := w.Write([]byte(`class="` + strings.Join(res, " ") + `"`))
	if err != nil {
		return err
	}
	return nil
}

func (c Classes) Add(s string) Classes {
	return c + Classes(" ") + Classes(s)
}

func (c Classes) Merge(oth Classes) Classes {
	return c.Add(string(oth))
}

func (c Classes) MergeMap(m ClassesMap) Classes {
	cmap := ClassesMap{}
	for _, cls := range extractClasses(string(c)) {
		cmap[cls] = true
	}

	for cls, ok := range m {
		cmap[cls] = ok
	}

	return Classes(cmap.extract())
}

type ClassesMap map[string]bool

func (c ClassesMap) attr()            {}
func (c ClassesMap) Node() Renderable { return c }

func (c ClassesMap) extract() string {
	var sb strings.Builder
	cnt := len(c)
	ind := -1
	for k, v := range c {
		ind++
		if !v {
			continue
		}
		sb.WriteString(k)
		if ind < cnt-1 {
			sb.WriteString(" ")
		}

	}
	return sb.String()
}

func (c ClassesMap) Render(w io.Writer) error {
	s := c.extract()
	_, err := w.Write([]byte(`class="` + s + `"`))
	if err != nil {
		return err
	}
	return nil
}

func (c ClassesMap) Add(s string) ClassesMap {
	res := extractClasses(s)

	newMap := maps.Clone(c)
	for _, cls := range res {
		newMap[cls] = true
	}
	return newMap
}

// All common attributes

func ID(id string) Attribute { return BuildAttr("id", id) }

func OnClick(handler string) Attribute     { return BuildAttr("onclick", handler) }
func OnChange(handler string) Attribute    { return BuildAttr("onchange", handler) }
func OnMouseOver(handler string) Attribute { return BuildAttr("onmouseover", handler) }
func OnMouseOut(handler string) Attribute  { return BuildAttr("onmouseout", handler) }
func OnMouseDown(handler string) Attribute { return BuildAttr("onmousedown", handler) }
func OnMouseUp(handler string) Attribute   { return BuildAttr("onmouseup", handler) }
func OnFocus(handler string) Attribute     { return BuildAttr("onfocus", handler) }
func OnBlur(handler string) Attribute      { return BuildAttr("onblur", handler) }
func OnKeyDown(handler string) Attribute   { return BuildAttr("onkeydown", handler) }
func OnKeyPress(handler string) Attribute  { return BuildAttr("onkeypress", handler) }
func OnKeyUp(handler string) Attribute     { return BuildAttr("onkeyup", handler) }
func OnLoad(handler string) Attribute      { return BuildAttr("onload", handler) }
func OnSubmit(handler string) Attribute    { return BuildAttr("onsubmit", handler) }
func OnReset(handler string) Attribute     { return BuildAttr("onreset", handler) }
func OnSelect(handler string) Attribute    { return BuildAttr("onselect", handler) }
func OnAbort(handler string) Attribute     { return BuildAttr("onabort", handler) }
func OnError(handler string) Attribute     { return BuildAttr("onerror", handler) }
func OnResize(handler string) Attribute    { return BuildAttr("onresize", handler) }
func OnScroll(handler string) Attribute    { return BuildAttr("onscroll", handler) }
func OnUnload(handler string) Attribute    { return BuildAttr("onunload", handler) }

// Global attributes

func Class(class string) Attribute         { return BuildAttr("class", class) }
func Lang(lang string) Attribute           { return BuildAttr("lang", lang) }
func Dir(dir string) Attribute             { return BuildAttr("dir", dir) }
func AccessKey(accessKey string) Attribute { return BuildAttr("accesskey", accessKey) }
func TabIndex(tabIndex string) Attribute   { return BuildAttr("tabindex", tabIndex) }
func ContentEditable(contentEditable string) Attribute {
	return BuildAttr("contenteditable", contentEditable)
}
func ContextMenu(contextMenu string) Attribute   { return BuildAttr("contextmenu", contextMenu) }
func Draggable(draggable string) Attribute       { return BuildAttr("draggable", draggable) }
func DropZone(dropZone string) Attribute         { return BuildAttr("dropzone", dropZone) }
func Hidden(hidden string) Attribute             { return BuildAttr("hidden", hidden) }
func SpellCheck(spellCheck string) Attribute     { return BuildAttr("spellcheck", spellCheck) }
func Translate(translate string) Attribute       { return BuildAttr("translate", translate) }
func Role(role string) Attribute                 { return BuildAttr("role", role) }
func AutoFocus(autoFocus string) Attribute       { return BuildAttr("autofocus", autoFocus) }
func AutoComplete(autoComplete string) Attribute { return BuildAttr("autocomplete", autoComplete) }
func AutoSave(autoSave string) Attribute         { return BuildAttr("autosave", autoSave) }
func FormAttr(form string) Attribute             { return BuildAttr("form", form) }
func EncType(encType string) Attribute           { return BuildAttr("enctype", encType) }
func Accept(accept string) Attribute             { return BuildAttr("accept", accept) }
func FormAction(formAction string) Attribute     { return BuildAttr("formaction", formAction) }
func FormEncType(formEncType string) Attribute   { return BuildAttr("formenctype", formEncType) }
func FormMethod(formMethod string) Attribute     { return BuildAttr("formmethod", formMethod) }
func FormNoValidate(formNoValidate string) Attribute {
	return BuildAttr("formnovalidate", formNoValidate)
}
func FormTarget(formTarget string) Attribute   { return BuildAttr("formtarget", formTarget) }
func List(list string) Attribute               { return BuildAttr("list", list) }
func Max(max string) Attribute                 { return BuildAttr("max", max) }
func Min(min string) Attribute                 { return BuildAttr("min", min) }
func Multiple(multiple string) Attribute       { return BuildAttr("multiple", multiple) }
func Pattern(pattern string) Attribute         { return BuildAttr("pattern", pattern) }
func Placeholder(placeholder string) Attribute { return BuildAttr("placeholder", placeholder) }
func ReadOnly() NoValAttribute                 { return NoValAttr("readonly") }
func Required() NoValAttribute                 { return NoValAttr("required") }
func Size(size string) Attribute               { return BuildAttr("size", size) }
func Src(src string) Attribute                 { return BuildAttr("src", src) }
func Step(step string) Attribute               { return BuildAttr("step", step) }
func Width(width string) Attribute             { return BuildAttr("width", width) }
func Height(height string) Attribute           { return BuildAttr("height", height) }
func Alt(alt string) Attribute                 { return BuildAttr("alt", alt) }
func UseMap(useMap string) Attribute           { return BuildAttr("usemap", useMap) }
func IsMap(isMap string) Attribute             { return BuildAttr("ismap", isMap) }
func LongDesc(longDesc string) Attribute       { return BuildAttr("longdesc", longDesc) }
func SrcSet(srcSet string) Attribute           { return BuildAttr("srcset", srcSet) }
func Sizes(sizes string) Attribute             { return BuildAttr("sizes", sizes) }
func CrossOrigin(crossOrigin string) Attribute { return BuildAttr("crossorigin", crossOrigin) }
func Media(media string) Attribute             { return BuildAttr("media", media) }
func Type(type_ string) Attribute              { return BuildAttr("type", type_) }
func Charset(charset string) Attribute         { return BuildAttr("charset", charset) }
func Href(href string) Attribute               { return BuildAttr("href", href) }
func HrefLang(hrefLang string) Attribute       { return BuildAttr("hreflang", hrefLang) }
func Rel(rel string) Attribute                 { return BuildAttr("rel", rel) }
func Rev(rev string) Attribute                 { return BuildAttr("rev", rev) }
func Target(target string) Attribute           { return BuildAttr("target", target) }
func Download(download string) Attribute       { return BuildAttr("download", download) }
func Ping(ping string) Attribute               { return BuildAttr("ping", ping) }
func ReferrerPolicy(referrerPolicy string) Attribute {
	return BuildAttr("referrerpolicy", referrerPolicy)
}
func Integrity(integrity string) Attribute { return BuildAttr("integrity", integrity) }
func Content(content string) Attribute     { return BuildAttr("content", content) }
func HttpEquiv(httpEquiv string) Attribute { return BuildAttr("http-equiv", httpEquiv) }
func Name(name string) Attribute           { return BuildAttr("name", name) }
func Scheme(scheme string) Attribute       { return BuildAttr("scheme", scheme) }
func Coords(coords string) Attribute       { return BuildAttr("coords", coords) }
func Shape(shape string) Attribute         { return BuildAttr("shape", shape) }
func Axis(axis string) Attribute           { return BuildAttr("axis", axis) }
func Headers(headers string) Attribute     { return BuildAttr("headers", headers) }
func Scope(scope string) Attribute         { return BuildAttr("scope", scope) }
func ColSpan(colSpan string) Attribute     { return BuildAttr("colspan", colSpan) }
func RowSpan(rowSpan string) Attribute     { return BuildAttr("rowspan", rowSpan) }
func Action(action string) Attribute       { return BuildAttr("action", action) }
func Method(method string) Attribute       { return BuildAttr("method", method) }
func NoValidate() NoValAttribute           { return NoValAttr("novalidate") }

// Additional common attributes
func TitleAttr(title string) Attribute             { return BuildAttr("title", title) }
func StyleAttr(style string) Attribute             { return BuildAttr("style", style) }
func DataAttr(name string, value string) Attribute { return BuildAttr("data-"+name, value) }
func Aria(name string, value string) Attribute     { return BuildAttr("aria-"+name, value) }
func Disabled() NoValAttribute                     { return NoValAttr("disabled") }
func Checked() NoValAttribute                      { return NoValAttr("checked") }
func Value(value string) Attribute                 { return BuildAttr("value", value) }
func MaxLength(maxLength string) Attribute         { return BuildAttr("maxlength", maxLength) }
func MinLength(minLength string) Attribute         { return BuildAttr("minlength", minLength) }
func SrcLang(srclang string) Attribute             { return BuildAttr("srclang", srclang) }
func Kind(kind string) Attribute                   { return BuildAttr("kind", kind) }
func LabelAttr(label string) Attribute             { return BuildAttr("label", label) }
func Default(default_ string) Attribute            { return BuildAttr("default", default_) }
func KeyType(keytype string) Attribute             { return BuildAttr("keytype", keytype) }
func DateTime(datetime string) Attribute           { return BuildAttr("datetime", datetime) }
func For(for_ string) Attribute                    { return BuildAttr("for", for_) }
func HeadersAttr(headers string) Attribute         { return BuildAttr("headers", headers) }
func High(high string) Attribute                   { return BuildAttr("high", high) }
func Low(low string) Attribute                     { return BuildAttr("low", low) }
func Optimum(optimum string) Attribute             { return BuildAttr("optimum", optimum) }
