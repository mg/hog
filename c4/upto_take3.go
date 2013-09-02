package main

import (
	"fmt"
)

type upto struct {
	m, n int
}

func Upto(m, n int) *upto {
	return &upto{m: m, n: n}
}

func (i *upto) Next() {
	i.m++
}

func (i *upto) AtEnd() bool {
	return i.m >= i.n
}

func (i *upto) Value() int {
	return i.m
}

func main() {
	for i := Upto(3, 5); !i.AtEnd(); i.Next() {
		fmt.Println(i.Value())
	}
}
