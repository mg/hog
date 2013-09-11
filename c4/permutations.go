package main

import (
	"fmt"
	"github.com/mg/i"
	"github.com/mg/i/igen"
	"math"
	"os"
)

type permute struct {
	items []rune
	n     int
	atEnd bool
}

func Permute(items []rune) i.Forward {
	return &permute{items: items, n: 1, atEnd: false}
}

func (p *permute) Error() error {
	return nil
}

func (p *permute) SetError(err error) {
}

func (p *permute) Value() interface{} {
	return p.items
}

func (p *permute) AtEnd() bool {
	return p.atEnd
}

func (per *permute) Next() error {
	var i int
	p := per.n
	for i = 1; i <= len(per.items) && math.Mod(float64(p), float64(i)) == 0; i++ {
		p /= i
	}
	d := int(math.Mod(float64(p), float64(i)))
	j := len(per.items) - i
	if j < 0 {
		per.atEnd = true
		return nil
	}

	copy(per.items[j+1:len(per.items)], reverse(per.items[j+1:len(per.items)]))
	per.items[j], per.items[j+d] = per.items[j+d], per.items[j]
	per.n++
	return nil
}

func reverse(in []rune) []rune {
	out := make([]rune, len(in))
	for i, v := range in {
		out[len(out)-i-1] = v
	}
	return out
}

func generate(from, to rune) []rune {
	list := make([]rune, 0, to-from+1)
	for itr := igen.Range(int(from), int(to)+1); !itr.AtEnd(); itr.Next() {
		list = append(list, rune(itr.Int()))
	}
	return list
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage %s FROM TO\n", os.Args[0])
		os.Exit(0)
	}
	from, to := rune(os.Args[1][0]), rune(os.Args[2][0])

	i.Each(
		Permute(generate(from, to)),
		func(itr i.Iterator) bool {
			r, _ := itr.Value().([]rune)
			fmt.Println(string(r))
			return true
		})
}
