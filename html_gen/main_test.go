package html_gen_test

import (
	"fmt"

	"github.com/shurcooL/go/html_gen"
)

func ExampleRenderNodes() {
	// Context-aware escaping is done just like in html/template.
	html, err := html_gen.RenderNodes(
		html_gen.Text("Hi & how are you, "),
		html_gen.A("Gophers", "https://golang.org/"),
		html_gen.Text("? <script> is a cool gopher."),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(html)

	// Output:
	//Hi &amp; how are you, <a href="https://golang.org/">Gophers</a>? &lt;script&gt; is a cool gopher.
}
