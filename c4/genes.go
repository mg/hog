package main

import (
	"bytes"
	"fmt"
	"github.com/mg/i"
	"math"
	"os"
	"regexp"
)

type gene struct {
	atEnd bool
	data  []interface{}
	cur   string
}

func Gene(pattern string) i.Forward {
	r := regexp.MustCompile("[()]")
	tokens := r.Split(pattern, -1)
	g := gene{atEnd: false, data: make([]interface{}, 0, len(tokens))}
	for i, v := range tokens {
		if math.Mod(float64(i), 2) == 0 {
			// constant
			if v != "" {
				g.data = append(g.data, v)
			}
		} else {
			// option
			if v != "" {
				option := make([]interface{}, 0, len(v)+1)
				option = append(option, 0)
				for _, c := range v {
					option = append(option, string(c))
				}
				g.data = append(g.data, option)
			}
		}
	}
	g.Next()
	return &g
}

func (g *gene) Error() error {
	return nil
}

func (g *gene) SetError(err error) {
}

func (g *gene) Value() interface{} {
	return g.cur
}

func (g *gene) AtEnd() bool {
	return g.cur == ""
}

func (g *gene) Next() error {
	if g.atEnd {
		g.cur = ""
		return nil
	}
	finishedIncr := false
	var result bytes.Buffer
	for _, t := range g.data {
		if v, ok := t.(string); ok {
			result.WriteString(v)
		} else if v, ok := t.([]interface{}); ok {
			n := v[0].(int)
			result.WriteString(v[n+1].(string))
			if !finishedIncr {
				if n == len(v)-2 {
					v[0] = 0
				} else {
					v[0] = v[0].(int) + 1
					finishedIncr = true
				}
			}
		}
	}
	g.atEnd = !finishedIncr
	g.cur = result.String()
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage %s PATTERN\n", os.Args[0])
		os.Exit(0)
	}

	i.Each(
		Gene(os.Args[1]),
		func(itr i.Iterator) bool {
			fmt.Println(itr.Value())
			return true
		})
}
