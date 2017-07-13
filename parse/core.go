package parse

import (
	"strings"

	"errors"

	"golang.org/x/net/html"
)

var (
	// ErrNoSuchElement ...
	ErrNoSuchElement = errors.New("No such element")
)

// GetAttr ...
func GetAttr(n *html.Node, attrName string) string {
	if n == nil {
		return ""
	}
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func getElementWithClass(n *html.Node, tag, class string) (res *html.Node) {
	if n.Type == html.ElementNode && n.Data == tag {
		for _, attr := range n.Attr {
			if attr.Key == "class" && containsAllClasses(attr.Val, class) {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res = getElementWithClass(c, tag, class)
		if res != nil {
			return res
		}
	}
	return res
}

// GetTextFromTag ...
func GetTextFromTag(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			return c.Data
		}
	}
	return ""
}

// GetElementByID ...
func GetElementByID(n *html.Node, id string) (*html.Node, error) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return n, nil
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res, _ := GetElementByID(c, id)
		if res != nil {
			return res, nil
		}
	}
	return nil, ErrNoSuchElement
}

func containsAllClasses(val, class string) bool {
	valAr := strings.Split(val, " ")
	classAr := strings.Split(class, " ")
	m := make(map[string]bool)
	for _, v := range valAr {
		m[v] = true
	}
	for _, c := range classAr {
		_, ok := m[c]
		if !ok {
			return false
		}
	}
	return true
}
