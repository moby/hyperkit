// Package html_table provides WalkRows to walk the rows of an html table.
package html_table

import (
	"golang.org/x/net/html"

	"github.com/shurcooL/go/html_to_markdown"
)

// WalkRows walks the rows of an html table, calling walkFunc on each row.
func WalkRows(htmlTable *html.Node, walkFunc func(columns ...string)) {
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {

			var columns []string

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					columns = append(columns, html_to_markdown.Paragraph(c))
				}
			}

			walkFunc(columns...)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(htmlTable)
}
