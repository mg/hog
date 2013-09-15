package main

import (
	"fmt"
	"github.com/mg/i"
	"github.com/mg/i/hoi"
	"github.com/mg/i/igen"
)

type iappend struct {
	itrs []i.Forward
	pos  int
	err  error
}

func Append(itrs ...i.Forward) i.Forward {
	return &iappend{itrs: itrs}
}

func (i *iappend) AtEnd() bool {
	if i.pos < len(i.itrs)-1 {
		return false
	} else if i.pos >= len(i.itrs) {
		return true
	}
	return i.itrs[i.pos].AtEnd()
}

func (i *iappend) Next() error {
	if i.pos >= len(i.itrs) {
		i.err = fmt.Errorf("Appending beyond last iterator")
		return i.err
	}
	i.itrs[i.pos].Next()
	for i.itrs[i.pos].AtEnd() {
		i.pos++
		if i.pos >= len(i.itrs) {
			return nil
		}
	}
	return nil
}

func (i *iappend) Value() interface{} {
	return i.itrs[i.pos].Value()
}

func (i *iappend) Error() error {
	if i.err != nil {
		return i.err
	}
	for _, v := range i.itrs {
		err := v.Error()
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *iappend) SetError(err error) {
	i.err = err
}

func main() {
	for itr := Append(hoi.List(1, 2, 3), igen.Range(10, 20)); !itr.AtEnd(); itr.Next() {
		fmt.Println(itr.Value())
	}

}
