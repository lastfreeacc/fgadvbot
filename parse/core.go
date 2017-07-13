package parse

import (
	"strings"
	
	"golang.org/x/net/html"
	"errors"
)

var (
	// ErrNoSuchElement ...
	ErrNoSuchElement = errors.New("No such element")
)

func getAtr(n *html.Node, atrName string) string {
	if n == nil {
		return ""
	}
	for _, atr := range n.Attr {
		if atr.Key == atrName {
			return atr.Val
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
		res, err := GetElementByID(c, id)
		if err != nil {
			return nil, err
		}
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