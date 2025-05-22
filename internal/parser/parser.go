package parser

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/michaeldebetaz/unilike/internal/assert"
	"golang.org/x/net/html"
)

func parseHtml(htmlStr string) (*html.Node, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		err := fmt.Errorf("Error while parsing HTML: %v", err)
		return nil, err
	}

	return doc, nil
}

func getElementsByTag(n *html.Node, tag string) []*html.Node {
	var nodes []*html.Node
	if n.Type == html.ElementNode && n.Data == tag {
		nodes = append(nodes, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, getElementsByTag(c, tag)...)
	}
	return nodes
}

func getElementById(n *html.Node, id string) (*html.Node, error) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return n, nil
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found, _ := getElementById(c, id); found != nil {
			return found, nil
		}
	}

	err := fmt.Errorf("Element with id %s not found", id)
	return nil, err
}

func getElementsByClass(n *html.Node, class string) []*html.Node {
	var nodes []*html.Node
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, class) {
				nodes = append(nodes, n)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, getElementsByClass(c, class)...)
	}

	return nodes
}

func getAttributeValue(n *html.Node, attrName string) string {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == attrName {
				return attr.Val
			}
		}
	}
	return ""
}

func getInnerText(n *html.Node) string {
	var buf bytes.Buffer

	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	} else if n.Type == html.ElementNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			buf.WriteString(getInnerText(c))
		}
	}

	return strings.TrimSpace(buf.String())
}

func getFirstNode(nodes []*html.Node, msg string) *html.Node {
	return assert.At(nodes, 0, msg)
}

func toString(nodes ...*html.Node) string {
	var b bytes.Buffer

	for _, n := range nodes {
		chunk := bytes.Buffer{}
		if err := html.Render(&chunk, n); err != nil {
			err := fmt.Errorf("Error while rendering HTML: %v", err)
			log.Fatal(err)
		}
		b.Write(chunk.Bytes())
	}
	return b.String()
}
