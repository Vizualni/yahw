package main

import (
	"context"
	"net/http"

	. "github.com/vizualni/yahw"
)

type MyCustomButton struct {
	Text            string
	BackgroundColor string
}

func (m MyCustomButton) Node(ctx context.Context) Renderable {
	return Button(
		BuildAttr("style", "background-color: "+m.BackgroundColor),
		Text(m.Text),
	)
}

func MyCustomInput(name, placeholder string) Node {
	return Input(
		BuildAttr("name", name),
		BuildAttr("placeholder", placeholder),
	)
}

func MyCommonAttributes(link string) Node {
	return AttrSlice{BuildAttr("id", "my-id"), Classes("my-1 my-2 my-1"), BuildAttr("href", link)}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		root := NewHTML5Doctype(
			HTML(
				Head(
					Title((Text("My Custom Button Example")),
						Style((Text("button { padding: 10px; border: none; }"))),
						Body(
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
							A(MyCommonAttributes("https://example1.com"), Text("Click me!")),
							Br(),
							A(MyCommonAttributes("https://example2.com"), Text("No, click me!")),
						),
					),
				),
			),
		)

		err := root.Render(w)
		if err != nil {
			panic(err)
		}
	})

	if err := http.ListenAndServe("127.0.0.1:8585", nil); err != nil {
		panic(err)
	}
}
