package html_to_markdown

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
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

// getAttribute returns an attribute of a node, or blank strink if not found.
func getAttribute(n *html.Node, key string) (val string) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// extract returns the recursive concatenation of the raw text contents of an html node, with Markdown tags.
func extract(n *html.Node) (out string) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			out += c.Data
		} else if c.Type == html.ElementNode && c.DataAtom == atom.A {
			out += "[" + extract(c) + "](" + getAttribute(c, "href") + ")"
		} else if c.Type == html.ElementNode && c.DataAtom == atom.Img {
			out += "![" + getAttribute(c, "title") + "](" + getAttribute(c, "src") + ")"
		} else if c.Type == html.ElementNode && c.DataAtom == atom.Blockquote {
			out += "> " + extract(c)
		} else if c.Type == html.ElementNode && (c.DataAtom == atom.Ul || c.DataAtom == atom.Ol) {
			out += extractList(c, c.DataAtom)
		} else if c.Type == html.ElementNode && (c.DataAtom == atom.B || c.DataAtom == atom.Strong) {
			out += "**" + extract(c) + "**"
		} else {
			out += extract(c)
		}
	}
	return out
}

func extractList(n *html.Node, listType atom.Atom) (out string) {
	firstItem := true

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.DataAtom == atom.Li {
			if firstItem {
				firstItem = false
			} else {
				out += "\n"
			}

			switch listType {
			case atom.Ul:
				out += "- "
			case atom.Ol:
				out += "1. "
			}
			out += extract(c)
		}
	}
	return out
}

// Paragraph renders a single paragraph html node to Markdown.
func Paragraph(paragraph *html.Node) (markdown string) {
	return clean(extract(paragraph))
}

func Unnamed1(x *html.Node) (markdown string) {
	return extract(x)
}
