package main

import (
	"fmt"
	"github.com/mg/i"
	"github.com/mg/i/hoi"
	"github.com/mg/i/icon"
)

type eachlike struct {
	itr i.Forward
	f   hoi.MapFunc
}

func EachLike(f hoi.MapFunc, itr i.Forward) *eachlike {
	return &eachlike{f: f, itr: itr}
}

func (i *eachlike) Value() interface{} {
	ret := make([]interface{}, 2)
	ret[0] = i.itr.Value()
	ret[1] = i.f(i.itr)
	return ret
}

func (i *eachlike) ValuePair() (interface{}, interface{}) {
	res := i.Value().([]interface{})
	return res[0], res[1]
}

func (i *eachlike) Error() error {
	return i.itr.Error()
}

func (i *eachlike) Next() error {
	return i.itr.Next()
}

func (i *eachlike) AtEnd() bool {
	return i.itr.AtEnd()
}

func multiply(n int) hoi.MapFunc {
	return func(itr i.Iterator) interface{} {
		return n * itr.Value().(int)
	}
}

func main() {
	itr := EachLike(multiply(10), icon.List(1, 2, 3, 4, 5, 6, 7, 8, 9))
	for ; !itr.AtEnd(); itr.Next() {
		val, res := itr.ValuePair()
		fmt.Printf("Value: %v, result: %v\n", val, res)
	}
}
