package main

import (
	"fmt"
	"os"

	"github.com/vizualni/yahw/parsehtml"
)

func main() {
	code := parsehtml.GenerateGo(os.Stdin)
	fmt.Println(code)
}
