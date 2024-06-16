# Yet Another HTML Wrapper

## Example

Check the example in [`example/`](https://github.com/Vizualni/yahw/blob/main/example/main.go) folder.

```go
package main

import (
	"io"
	"net/http"

	. "github.com/vizualni/yahw"
)

type MyCustomButton struct {
	Text            string
	BackgroundColor string
}

func (m MyCustomButton) TagRender(w io.Writer) error {
	return Button(
		Attr("style", "background-color: "+m.BackgroundColor),
	).X(
		Text(m.Text),
	).TagRender(w)
}

func MyCustomInput(name, placeholder string) TagRenderer {
	return Input(
		Attr("name", name),
		Attr("placeholder", placeholder),
	)
}

func MyCommonAttributes(link string) AttrRenderer {
	return AttrSlice{Attr("id", "my-id"), Classes("my-1 my-2 my-1"), Attr("href", link)}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		root := NewHTML5Doctype().X(
			HTML().X(
				Head().X(
					Title().X(Text("My Custom Button Example")),
					Style().X(Text("button { padding: 10px; border: none; }")),
				),
				Body().X(
					MyCustomButton{
						Text:            "Click me!",
						BackgroundColor: "red",
					},
					Br(),
					MyCustomButton{
						Text:            "No, click me!",
						BackgroundColor: "green",
					},
					Br(),
					MyCustomInput("name", "Enter your name"),
					Br(),
					MyCustomInput("email", "Enter your email"),
					Br(),
					A(MyCommonAttributes("https://example1.com")).X(Text("Click me!")),
					Br(),
					A(MyCommonAttributes("https://example2.com")).X(Text("No, click me!")),
				),
			),
		)

		root.TagRender(w)
	})

	if err := http.ListenAndServe("127.0.0.1:8585", nil); err != nil {
		panic(err)
	}
}

```

Run `go run ./example` and visit http://localhost:8585/


## Why?

It's a common pattern I am using in my pet projects. I thought I might share it as it could be useful for others.
