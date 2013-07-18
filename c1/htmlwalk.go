package main

import (
	"bufio"
	"bytes"
	"code.google.com/p/go.net/html"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type (
	ResultType  interface{}
	TextFunc    func(n *html.Node) ResultType
	ElementFunc func(n *html.Node, results []ResultType) ResultType
)

func htmlwalk(n *html.Node, textf TextFunc, elementf ElementFunc) ResultType {
	if n.Type == html.TextNode {
		return textf(n)
	}
	results := make([]ResultType, 0)
	child := n.FirstChild
	for {
		if child == nil {
			break
		}
		result := htmlwalk(child, textf, elementf)
		if multiresults, ok := result.([]ResultType); ok {
			results = append(results, multiresults...)
		} else {
			results = append(results, result)
		}
		child = child.NextSibling
	}
	return elementf(n, results)
}

func stringValue(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func untagText(n *html.Node) ResultType {
	return n.Data
}

func untagElement(n *html.Node, results []ResultType) ResultType {
	var b bytes.Buffer
	for _, v := range results {
		s := stringValue(v)
		_, err := b.WriteString(s)
		if err != nil {
			panic(err)
		}
	}
	return b.String()
}

type TagType int16

const (
	Maybe TagType = iota
	Keep
)

type tag struct {
	Type  TagType
	Value string
}

func tagValue(v interface{}) (*tag, []ResultType) {
	if t, ok := v.(*tag); ok {
		return t, nil
	} else if t, ok := v.([]ResultType); ok {
		return nil, t
	}
	panic(fmt.Sprintf("Unkown value %v", v))
}

func promoteText(n *html.Node) ResultType {
	return &tag{Maybe, n.Data}
}

func promoteElementIf(tagname string) func(n *html.Node, results []ResultType) ResultType {
	return func(n *html.Node, results []ResultType) ResultType {
		if n.Data == tagname {
			var b bytes.Buffer
			for _, v := range results {
				t, _ := tagValue(v)
				_, err := b.WriteString(t.Value)
				if err != nil {
					panic(err)
				}
			}
			return &tag{Keep, b.String()}
		}
		return results
	}
}

func extractPromoted(n *html.Node) string {
	_, results := tagValue(htmlwalk(n, promoteText, promoteElementIf("h1")))

	if results != nil {
		var b bytes.Buffer
		for _, r := range results {
			tag, _ := tagValue(r)
			if tag.Type == Keep {
				_, err := b.WriteString(tag.Value + " ")
				if err != nil {
					panic(err)
				}
			}
		}
		return b.String()
	}
	return ""
}

var url = flag.Bool("u", false, "Use -u for url")

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Printf("Usage %s [-u] NAME\n", os.Args[0])
		os.Exit(0)
	}

	var in io.Reader
	if *url {
		resp, err := http.Get(os.Args[2])
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		in = resp.Body
	} else {
		fi, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fi.Close()
		in = bufio.NewReader(fi)
	}

	n, err := html.Parse(in)
	if err != nil {
		panic(err)
	}

	fmt.Println(stringValue(htmlwalk(n, untagText, untagElement)))
	fmt.Println(extractPromoted(n))
}
