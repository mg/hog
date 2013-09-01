package main

import (
	"fmt"
)

type upto struct {
	m, n int
}

func Upto(m, n int) *upto {
	return &upto{m: m - 1, n: n}
}

func (i *upto) Next() bool {
	i.m++
	return i.m < i.n
}

func (i *upto) Value() int {
	return i.m
}

func main() {
	i := Upto(3, 5)
	for i.Next() {
		fmt.Println(i.Value())
	}
}
