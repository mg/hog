package main

import (
	"fmt"
	"github.com/mg/i"
	"github.com/mg/i/igen"
	"os"
	"strconv"
)

type ChildrenFunc func(interface{}) []interface{}

type depthfirst struct {
	children ChildrenFunc
	agenda   []interface{}
	node     interface{}
	err      error
}

func Dfs(root interface{}, children ChildrenFunc) i.Forward {
	dfs := depthfirst{children: children}
	dfs.agenda = append(dfs.agenda, root)
	dfs.Next()
	return &dfs
}

func (dfs *depthfirst) AtEnd() bool {
	return len(dfs.agenda) == 0
}

func (dfs *depthfirst) Value() interface{} {
	if dfs.AtEnd() {
		dfs.err = fmt.Errorf("Next: Beyond end")
		return nil
	}
	return dfs.node
}

func (dfs *depthfirst) Error() error {
	return dfs.err
}

func (dfs *depthfirst) SetError(err error) {
}

func (dfs *depthfirst) Next() error {
	if dfs.AtEnd() {
		dfs.err = fmt.Errorf("Next: Beyond end")
		return dfs.err
	}
	dfs.node, dfs.agenda = dfs.agenda[0], dfs.agenda[1:len(dfs.agenda)]
	dfs.agenda = append(dfs.agenda, dfs.children(dfs.node)...)
	return nil
}

func IntPartition(n int) i.Forward {
	return Dfs([]int{n}, func(node interface{}) []interface{} {
		items, _ := node.([]int)
		largest, rest := items[0], items[1:len(items)]
		min, max := 1, largest/2
		if len(rest) > 0 {
			min = rest[0]
		}
		next := make([]interface{}, 0)
		for r := igen.Range(min, max+1); !r.AtEnd(); r.Next() {
			next = append(next, (append([]int{largest - r.Int(), r.Int()}, rest...)))
		}
		return next
	})
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage %s N\n", os.Args[0])
		os.Exit(0)
	}
	n, _ := strconv.Atoi(os.Args[1])

	itr := IntPartition(n)
	for ; !itr.AtEnd(); itr.Next() {
		fmt.Println(itr.Value())
	}

}
