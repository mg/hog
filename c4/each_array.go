package main

import (
	"fmt"
	"github.com/mg/i"
)

type stop int

const (
	StopAtMin stop = iota
	StopAtMax
)

type eacharray struct {
	itrs   []i.Forward
	err    error
	stopAt stop
}

func EachArray(stopAt stop, itrs ...i.Forward) i.Forward {
	return &eacharray{itrs: itrs, stopAt: stopAt}
}

func (i *eacharray) Value() interface{} {
	ret := make([]interface{}, len(i.itrs))
	for idx, itr := range i.itrs {
		if !itr.AtEnd() {
			ret[idx] = itr.Value()
		} else {
			ret[idx] = nil
		}
	}
	return ret
}

func (i *eacharray) Error() error {
	return i.err
}

func (i *eacharray) SetError(err error) {
	i.err = err
}

func (i *eacharray) Next() error {
	for _, itr := range i.itrs {
		if !itr.AtEnd() {
			err := itr.Next()
			if err != nil {
				i.err = err
				break
			}
		}
	}
	return i.err
}

func (i *eacharray) AtEnd() bool {
	atEndCount := 0
	for _, itr := range i.itrs {
		if itr.AtEnd() {
			if i.stopAt == StopAtMin {
				return true
			}
			atEndCount++
		}
	}
	return atEndCount == len(i.itrs)
}

func main() {
	itr := EachArray(
		StopAtMax,
		i.List(1, 2, 3, 4, 5, 6),
		i.List(6.4, 7.1, 8.2, 9.9),
		i.List("A", "B", "C", "D", "E"))
	for ; !itr.AtEnd(); itr.Next() {
		fmt.Println(itr.Value())
	}
}
