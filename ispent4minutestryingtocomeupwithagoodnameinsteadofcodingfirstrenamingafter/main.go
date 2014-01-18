package html_table

import (
	"code.google.com/p/go.net/html"
)

// clean replaces each sequence of space, \n, \r, or \t characters
// with a single space and removes any trailing and leading spaces.
func clean(s string) string {
	var b []byte
	p := byte(' ')
	for i := 0; i < len(s); i++ {
		q := s[i]
		if q == '\n' || q == '\r' || q == '\t' {
			q = ' '
		}
		if q != ' ' || p != ' ' {
			b = append(b, q)
			p = q
		}
	}
	// Remove trailing blank, if any.
	if n := len(b); n > 0 && p == ' ' {
		b = b[0 : n-1]
	}
	return string(b)
}

// extract returns the recursive concatenation of the raw text contents of an html node.
func extract(n *html.Node) (out string) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			out += c.Data
		} else {
			out += extract(c)
		}
	}
	return out
}

// htmlNodeToPlainText renders an html node to plain text.
func htmlNodeToPlainText(htmlNode *html.Node) (plainText string) {
	return clean(extract(htmlNode))
}

// WalkRows walks the rows of an html table, calling walkFunc on each row.
func WalkRows(htmlTable *html.Node, walkFunc func(columns ...string)) {
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {

			var columns []string

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					columns = append(columns, htmlNodeToPlainText(c))
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
