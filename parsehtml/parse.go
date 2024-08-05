package parsehtml

import (
	"fmt"
	"go/format"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func GenerateGo(in io.Reader) string {
	nodes, err := html.ParseFragment(in, &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	if err != nil {
		panic(err)
	}

	buf := strings.Builder{}
	buf.WriteString("return yahw.TagSlice{\n")
	for _, node := range nodes {
		if node.Type != html.ElementNode {
			continue
		}
		writeNode(&buf, node)
		buf.WriteString(",\n")
	}
	buf.WriteString("}")
	code := buf.String()
	fmt.Println(code)
	bz, err := format.Source([]byte(code))
	if err != nil {
		panic(err)
	}
	formatted := string(bz)
	return formatted
}

func writeNode(w io.Writer, node *html.Node) {
	if node == nil {
		return
	}
	if node.Type != html.ElementNode {
		panic("can only create a tag from the element node")
	}
	tag := node.Data
	code := fmt.Sprintf(`yahw.NewTag("%s")`, tag)
	w.Write([]byte(code))
	if len(node.Attr) > 0 {
		w.Write([]byte(".Attrs(\n"))
		writeAttrs(w, node)
		w.Write([]byte("\n)"))
	}
	w.Write([]byte(".X(\n"))
	for next := node.FirstChild; next != nil; next = next.NextSibling {
		switch next.Type {
		case html.ElementNode:
			writeNode(w, next)
			w.Write([]byte(",\n"))
		case html.TextNode:
			if len(strings.TrimSpace(next.Data)) > 0 {
				fmt.Fprintf(w, "yahw.Text(`%s`),", next.Data)
			}
		}
	}
	w.Write([]byte("\n)"))
}

func writeAttrs(w io.Writer, node *html.Node) {
	if node.Type != html.ElementNode {
		panic("can only get tags from the element node")
	}

	if len(node.Attr) == 0 {
		return
	}

	for _, attr := range node.Attr {
		code := fmt.Sprintf("yahw.Attr(\"%s\", `%s`),", attr.Key, attr.Val)
		w.Write([]byte(code))
		w.Write([]byte("\n"))
	}
}
