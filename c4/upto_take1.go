package main

import (
	"fmt"
)

type (
	Iter func() (int, bool)
)

func upto(m, n int) Iter {
	return func() (int, bool) {
		if m < n {
			ret := m
			m++
			return ret, true
		}
		return n, false
	}
}

func main() {
	i := upto(3, 5)
	for v, ok := i(); ok; v, ok = i() {
		fmt.Println(v)
	}
}
