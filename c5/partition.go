package main

import (
	"container/list"
	"fmt"
	"github.com/mg/i"
)

type entry struct {
	target      int
	pool, share []int
}

type partition struct {
	todo  *list.List
	err   error
	share []int
}

func Partition(target int, treasures []int) i.Forward {
	p := partition{}
	p.todo = list.New()
	if target == 0 {
		p.todo.PushBack(&entry{target: target, pool: make([]int, 0)})
	} else {
		p.todo.PushBack(&entry{target: target, pool: treasures})
	}
	p.Next()
	return &p
}

func (p *partition) AtEnd() bool {
	return p.todo.Len() == 0
}

func (p *partition) Value() interface{} {
	if p.AtEnd() {
		p.err = fmt.Errorf("Next: Beyond end")
		return nil
	}
	return p.share
}

func (p *partition) Error() error {
	return p.err
}

func (p *partition) SetError(err error) {
}

func (p *partition) Next() error {
	if p.AtEnd() {
		p.err = fmt.Errorf("Next: Beyond end")
		return p.err
	}
	for !p.AtEnd() {
		e, _ := p.todo.Front().Value.(*entry)
		p.todo.Remove(p.todo.Front())

		first, rest := e.pool[0], e.pool[1:len(e.pool)]
		if len(rest) > 0 {
			p.todo.PushBack(&entry{e.target, rest, e.share})
		}
		if e.target == first {
			p.share = append(e.share, first)
			return nil
		} else if e.target > first && len(rest) > 0 {
			p.todo.PushBack(&entry{e.target - first, rest, append(e.share, first)})
		}
	}
	return nil
}

func main() {
	treasures := []int{5, 2, 4, 8, 1}
	total := 0
	for _, v := range treasures {
		total += v
	}

	itr := Partition(total/2, treasures)
	for ; !itr.AtEnd(); itr.Next() {
		fmt.Println(itr.Value())
	}
}
