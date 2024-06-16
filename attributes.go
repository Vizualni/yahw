package yahw

import (
	"html"
	"io"
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

func Attr(key, value string) Attribute {
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
		return Attr(key, value)
	}
}

type Attribute struct {
	key   string
	value string
}

func (a Attribute) AttrRender(w io.Writer) error {
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

func (a NoValAttribute) AttrRender(w io.Writer) error {
	escapedKey := html.EscapeString(a.key)
	_, err := w.Write([]byte(escapedKey))
	if err != nil {
		return err
	}
	return nil
}

type AttrSlice []AttrRenderer

func (a AttrSlice) AttrRender(w io.Writer) error {
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
		err := attr.AttrRender(w)
		if err != nil {
			return err
		}
	}
	return nil
}

type Classes string

func (c Classes) AttrRender(w io.Writer) error {
	clss := strings.Split(string(c), " ")
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
	_, err := w.Write([]byte(`class="` + strings.Join(res, " ") + `"`))
	if err != nil {
		return err
	}
	return nil
}

// All common attributes

func ID(id string) Attribute { return Attr("id", id) }

func OnClick(handler string) Attribute     { return Attr("onclick", handler) }
func OnChange(handler string) Attribute    { return Attr("onchange", handler) }
func OnMouseOver(handler string) Attribute { return Attr("onmouseover", handler) }
func OnMouseOut(handler string) Attribute  { return Attr("onmouseout", handler) }
func OnMouseDown(handler string) Attribute { return Attr("onmousedown", handler) }
func OnMouseUp(handler string) Attribute   { return Attr("onmouseup", handler) }
func OnFocus(handler string) Attribute     { return Attr("onfocus", handler) }
func OnBlur(handler string) Attribute      { return Attr("onblur", handler) }
func OnKeyDown(handler string) Attribute   { return Attr("onkeydown", handler) }
func OnKeyPress(handler string) Attribute  { return Attr("onkeypress", handler) }
func OnKeyUp(handler string) Attribute     { return Attr("onkeyup", handler) }
func OnLoad(handler string) Attribute      { return Attr("onload", handler) }
func OnSubmit(handler string) Attribute    { return Attr("onsubmit", handler) }
func OnReset(handler string) Attribute     { return Attr("onreset", handler) }
func OnSelect(handler string) Attribute    { return Attr("onselect", handler) }
func OnAbort(handler string) Attribute     { return Attr("onabort", handler) }
func OnError(handler string) Attribute     { return Attr("onerror", handler) }
func OnResize(handler string) Attribute    { return Attr("onresize", handler) }
func OnScroll(handler string) Attribute    { return Attr("onscroll", handler) }
func OnUnload(handler string) Attribute    { return Attr("onunload", handler) }

// Global attributes

func Class(class string) Attribute         { return Attr("class", class) }
func Lang(lang string) Attribute           { return Attr("lang", lang) }
func Dir(dir string) Attribute             { return Attr("dir", dir) }
func AccessKey(accessKey string) Attribute { return Attr("accesskey", accessKey) }
func TabIndex(tabIndex string) Attribute   { return Attr("tabindex", tabIndex) }
func ContentEditable(contentEditable string) Attribute {
	return Attr("contenteditable", contentEditable)
}
func ContextMenu(contextMenu string) Attribute       { return Attr("contextmenu", contextMenu) }
func Draggable(draggable string) Attribute           { return Attr("draggable", draggable) }
func DropZone(dropZone string) Attribute             { return Attr("dropzone", dropZone) }
func Hidden(hidden string) Attribute                 { return Attr("hidden", hidden) }
func SpellCheck(spellCheck string) Attribute         { return Attr("spellcheck", spellCheck) }
func Translate(translate string) Attribute           { return Attr("translate", translate) }
func Role(role string) Attribute                     { return Attr("role", role) }
func AutoFocus(autoFocus string) Attribute           { return Attr("autofocus", autoFocus) }
func AutoComplete(autoComplete string) Attribute     { return Attr("autocomplete", autoComplete) }
func AutoSave(autoSave string) Attribute             { return Attr("autosave", autoSave) }
func FormAction(formAction string) Attribute         { return Attr("formaction", formAction) }
func FormEncType(formEncType string) Attribute       { return Attr("formenctype", formEncType) }
func FormMethod(formMethod string) Attribute         { return Attr("formmethod", formMethod) }
func FormNoValidate(formNoValidate string) Attribute { return Attr("formnovalidate", formNoValidate) }
func FormTarget(formTarget string) Attribute         { return Attr("formtarget", formTarget) }
func List(list string) Attribute                     { return Attr("list", list) }
func Max(max string) Attribute                       { return Attr("max", max) }
func Min(min string) Attribute                       { return Attr("min", min) }
func Multiple(multiple string) Attribute             { return Attr("multiple", multiple) }
func Pattern(pattern string) Attribute               { return Attr("pattern", pattern) }
func Placeholder(placeholder string) Attribute       { return Attr("placeholder", placeholder) }
func ReadOnly(readOnly string) Attribute             { return Attr("readonly", readOnly) }
func Required(required string) Attribute             { return Attr("required", required) }
func Size(size string) Attribute                     { return Attr("size", size) }
func Src(src string) Attribute                       { return Attr("src", src) }
func Step(step string) Attribute                     { return Attr("step", step) }
func Width(width string) Attribute                   { return Attr("width", width) }
func Height(height string) Attribute                 { return Attr("height", height) }
func Alt(alt string) Attribute                       { return Attr("alt", alt) }
func UseMap(useMap string) Attribute                 { return Attr("usemap", useMap) }
func IsMap(isMap string) Attribute                   { return Attr("ismap", isMap) }
func LongDesc(longDesc string) Attribute             { return Attr("longdesc", longDesc) }
func SrcSet(srcSet string) Attribute                 { return Attr("srcset", srcSet) }
func Sizes(sizes string) Attribute                   { return Attr("sizes", sizes) }
func CrossOrigin(crossOrigin string) Attribute       { return Attr("crossorigin", crossOrigin) }
func Media(media string) Attribute                   { return Attr("media", media) }
func Type(type_ string) Attribute                    { return Attr("type", type_) }
func Charset(charset string) Attribute               { return Attr("charset", charset) }
func Href(href string) Attribute                     { return Attr("href", href) }
func HrefLang(hrefLang string) Attribute             { return Attr("hreflang", hrefLang) }
func Rel(rel string) Attribute                       { return Attr("rel", rel) }
func Rev(rev string) Attribute                       { return Attr("rev", rev) }
func Target(target string) Attribute                 { return Attr("target", target) }
func Download(download string) Attribute             { return Attr("download", download) }
func Ping(ping string) Attribute                     { return Attr("ping", ping) }
func ReferrerPolicy(referrerPolicy string) Attribute { return Attr("referrerpolicy", referrerPolicy) }
func Integrity(integrity string) Attribute           { return Attr("integrity", integrity) }
func Content(content string) Attribute               { return Attr("content", content) }
func HttpEquiv(httpEquiv string) Attribute           { return Attr("http-equiv", httpEquiv) }
func Name(name string) Attribute                     { return Attr("name", name) }
func Scheme(scheme string) Attribute                 { return Attr("scheme", scheme) }
func Coords(coords string) Attribute                 { return Attr("coords", coords) }
func Shape(shape string) Attribute                   { return Attr("shape", shape) }
func Axis(axis string) Attribute                     { return Attr("axis", axis) }
func Headers(headers string) Attribute               { return Attr("headers", headers) }
func Scope(scope string) Attribute                   { return Attr("scope", scope) }
func ColSpan(colSpan string) Attribute               { return Attr("colspan", colSpan) }
func RowSpan(rowSpan string) Attribute               { return Attr("rowspan", rowSpan) }
func Action(action string) Attribute                 { return Attr("action", action) }
func Method(method string) Attribute                 { return Attr("method", method) }
func NoValidate(noValidate string) Attribute         { return Attr("novalidate", noValidate) }

// Additional common attributes
func TitleAttr(title string) Attribute             { return Attr("title", title) }
func StyleAttr(style string) Attribute             { return Attr("style", style) }
func DataAttr(name string, value string) Attribute { return Attr("data-"+name, value) }
func Aria(name string, value string) Attribute     { return Attr("aria-"+name, value) }
func Disabled() NoValAttribute                     { return NoValAttr("disabled") }
func Checked() NoValAttribute                      { return NoValAttr("checked") }
func Value(value string) Attribute                 { return Attr("value", value) }
func MaxLength(maxLength string) Attribute         { return Attr("maxlength", maxLength) }
func MinLength(minLength string) Attribute         { return Attr("minlength", minLength) }
func SrcLang(srclang string) Attribute             { return Attr("srclang", srclang) }
func Kind(kind string) Attribute                   { return Attr("kind", kind) }
func LabelAttr(label string) Attribute             { return Attr("label", label) }
func Default(default_ string) Attribute            { return Attr("default", default_) }
func KeyType(keytype string) Attribute             { return Attr("keytype", keytype) }
func DateTime(datetime string) Attribute           { return Attr("datetime", datetime) }
func For(for_ string) Attribute                    { return Attr("for", for_) }
func HeadersAttr(headers string) Attribute         { return Attr("headers", headers) }
func High(high string) Attribute                   { return Attr("high", high) }
func Low(low string) Attribute                     { return Attr("low", low) }
func Optimum(optimum string) Attribute             { return Attr("optimum", optimum) }
