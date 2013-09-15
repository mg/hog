package main

import (
	"fmt"
	"github.com/mg/i"
	"github.com/mg/i/hoi"
	"math"
)

type intslice struct {
	slice []int
	pos   int
}

func IntSlice(vals []int) i.Forward {
	return &intslice{slice: vals}
}

func (i *intslice) AtEnd() bool {
	return i.pos >= len(i.slice)
}

func (i *intslice) Next() error {
	i.pos++
	return nil
}

func (i *intslice) Value() interface{} {
	return i.slice[i.pos]
}

func (i *intslice) Error() error {
	return nil
}

func (i *intslice) SetError(err error) {
}

type list struct {
	slice []interface{}
	pos   int
}

func List(vals ...interface{}) i.Forward {
	return &list{slice: vals}
}

func (i *list) AtEnd() bool {
	return i.pos >= len(i.slice)
}

func (i *list) Next() error {
	i.pos++
	return nil
}

func (i *list) Value() interface{} {
	return i.slice[i.pos]
}

func (i *list) Error() error {
	return nil
}

func (i *list) SetError(err error) {
}

func main() {
	itr := IntSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
	itr = hoi.Filter(
		func(itr i.Iterator) bool {
			n, _ := itr.Value().(int)
			return math.Mod(float64(n), 3) == 0
		}, itr)

	for !itr.AtEnd() {
		fmt.Println(itr.Value())
		itr.Next()
	}

	itr = List(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	itr = hoi.Filter(func(itr i.Iterator) bool {
		n, _ := itr.Value().(int)
		return math.Mod(float64(n), 2) == 0
	}, itr)
	for !itr.AtEnd() {
		fmt.Println(itr.Value())
		itr.Next()
	}
}
