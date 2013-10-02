package main

import (
	"fmt"
	"github.com/mg/i"
	"github.com/mg/i/hoi"
	"github.com/mg/i/igen"
	"os"
	"sort"
	"strconv"
)

func IntPartitionRecursive(item []int) {
	fmt.Println(item)
	largest, rest := item[0], item[1:len(item)]
	min, max := 1, largest/2
	if len(rest) > 0 {
		min = rest[0]
	}
	for r := igen.Range(min, max+1); !r.AtEnd(); r.Next() {
		IntPartitionRecursive(append([]int{largest - r.Int(), r.Int()}, rest...))
	}
}

type items [][]int

func (it items) Len() int { return len(it) }
func (it items) Swap(i, j int) {
	it[i], it[j] = it[j], it[i]
}
func (it items) Less(i, j int) bool {
	a, b := it[i], it[j]
	for k := 0; k < len(a); k++ {
		if a[k] < b[k] {
			return false
		} else if a[k] > b[k] {
			return true
		}
	}
	// we will never get here since we've alredy pruned the
	// other half of the solution tree
	return true
}

type intpartition struct {
	agenda items
	item   []int
	err    error
}

func IntPartition(n int) i.Forward {
	var p intpartition
	item := make([]int, 1)
	item[0] = n
	p.agenda = append(p.agenda, item)
	p.Next()
	return &p
}

func (p *intpartition) AtEnd() bool {
	return len(p.agenda) == 0
}

func (p *intpartition) Value() interface{} {
	if p.AtEnd() {
		p.err = fmt.Errorf("Value: Beyond end")
		return nil
	}
	return p.item
}

func (p *intpartition) Error() error {
	return p.err
}

func (p *intpartition) SetError(err error) {
}

func (p *intpartition) Next() error {
	if p.AtEnd() {
		p.err = fmt.Errorf("Calling next AtEnd")
		return p.err
	}

	p.item, p.agenda = p.agenda[0], p.agenda[1:len(p.agenda)]

	largest, rest := p.item[0], p.item[1:len(p.item)]
	min, max := 1, largest/2
	if len(rest) > 0 {
		min = rest[0]
	}
	for r := igen.Range(min, max+1); !r.AtEnd(); r.Next() {
		p.agenda = append(p.agenda, (append([]int{largest - r.Int(), r.Int()}, rest...)))
	}
	sort.Sort(p.agenda)
	return nil
}

func distinct(itr i.Iterator) bool {
	item, _ := itr.Value().([]int)
	for i := 1; i < len(item); i++ {
		if item[i-1] == item[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage %s N\n", os.Args[0])
		os.Exit(0)
	}
	n, _ := strconv.Atoi(os.Args[1])

	fmt.Println("Partitions, recursive version, unsorted")
	IntPartitionRecursive([]int{n})

	fmt.Println("\nPartitions, itertor version, sorted")
	itr := IntPartition(n)
	for ; !itr.AtEnd(); itr.Next() {
		fmt.Println(itr.Value())
	}

	fmt.Println("\nDistinct only")
	itr = hoi.Filter(distinct, IntPartition(n))
	for ; !itr.AtEnd(); itr.Next() {
		fmt.Println(itr.Value())
	}
}
