package parser

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/michaeldebetaz/unilike/internal/assert"
	"golang.org/x/net/html"
)

type (
	Node  html.Node
	Nodes []Node
)

func (n *Nodes) first() *Node {
	node := assert.At(*n, 0)

	return &node
}

func (n *Node) attributeValue(attrName string) string {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == attrName {
				return attr.Val
			}
		}
	}
	return ""
}

func (n *Node) id(id string) (*Node, error) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return n, nil
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found, _ := (*Node)(c).id(id); found != nil {
			return found, nil
		}
	}

	err := fmt.Errorf("Element with id %s not found", id)
	return nil, err
}

func (n *Node) classes(class string) *Nodes {
	var nodes Nodes
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, class) {
				nodes = append(nodes, Node(*n))
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		node := Node(*c)
		nodes = append(nodes, *node.classes(class)...)
	}

	return &nodes
}

func (n *Node) tags(tag string) *Nodes {
	var nodes Nodes
	if n.Type == html.ElementNode && n.Data == tag {
		node := Node(*n)
		nodes = append(nodes, node)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		node := Node(*c)
		nodes = append(nodes, *node.tags(tag)...)
	}
	return &nodes
}

func (n *Node) innerText() string {
	var buf bytes.Buffer

	if n == nil {
		return ""
	}

	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	} else if n.Type == html.ElementNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			node := Node(*c)
			buf.WriteString(node.innerText())
		}
	}

	return strings.TrimSpace(buf.String())
}

func (n *Nodes) len() int {
	if n == nil {
		return 0
	}
	return len(*n)
}

func (n *Nodes) string() string {
	var b bytes.Buffer

	for _, node := range *n {
		chunk := bytes.Buffer{}
		if err := html.Render(&chunk, (*html.Node)(&node)); err != nil {
			err := fmt.Errorf("Error while rendering HTML: %v", err)
			log.Fatal(err)
		}
		b.Write(chunk.Bytes())
	}
	return b.String()
}

func parseHtml(htmlStr string) (*Node, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		err := fmt.Errorf("Error while parsing HTML: %v", err)
		return nil, err
	}

	node := Node(*doc)

	return &node, nil
}
