package spider

import (
	"code.google.com/p/go.net/html"
	"github.com/mg/i"
	"github.com/mg/i/hoi"
)

// Link iterator
// Filters html nodes, returning only link elements
func links(itr i.Iterator) bool {
	n, _ := itr.Value().(*html.Node)

	if n.Type != html.ElementNode {
		return false
	}
	if n.Data != "a" && n.Data != "img" && n.Data != "link" && n.Data != "style" && n.Data != "script" {
		return false
	}

	if n.Data == "style" || n.Data == "script" {
		src := attr("src", n.Attr)
		return src != ""
	}
	return true
}

func LinkItr(itr i.Forward) i.Forward {
	return hoi.Filter(links, itr)
}

// Url Iterator
// Maps html.Node elements to url strings
func attr(name string, attrs []html.Attribute) string {
	for _, a := range attrs {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

func geturl(itr i.Iterator) interface{} {
	n, _ := itr.Value().(*html.Node)
	var url string
	if n.Data == "a" {
		url = attr("href", n.Attr)
	} else if n.Data == "img" {
		url = attr("src", n.Attr)
	} else if n.Data == "link" {
		url = attr("href", n.Attr)
	} else if n.Data == "style" {
		url = attr("srr", n.Attr)
	} else if n.Data == "script" {
		url = attr("src", n.Attr)
	}
	return url
}

func UrlItr(itr i.Forward) i.Forward {
	return hoi.Map(geturl, itr)
}
