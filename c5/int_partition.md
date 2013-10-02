5.2: How to Convert a Recursive Function to an Iterator

The trick to turn a recursive function is to maintain an agenda of todo items, where each item is the state of each invocation of the recursive version. The agenda replaces the stack management we get when we use recursion. 

The integer partition problem is to break an integer into all possible integer components that sum to the original integer. E.g *6* breaks in to:

    6
    5 1
    4 1 1
    3 1 1 1
    2 1 1 1 1
    1 1 1 1 1 1
    2 2 1 1
    3 2 1
    4 2
    2 2 2
    3 3

A recursive solution to this problem looks a bit like this:

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

For every slice, starting with, in this case, *[6]*, we print out the slice. We when split it into a first element (*6*) and rest (*[]*) and set the min (*1*) and max (*3*) values. We then loop through the range *[1,4)* and call the function again 3 times with the values *[5,1]*, *[4,2]* and *[3,3]*. The call with *[5,1]* will in turn results to calls with *[4,1,1]* and [*3,2,1]*.

This strategy, of looping only to *largest/2*, prunes the solution tree so that we don't return solutions such as *[2,4]* since that's the same as *[4,2]*. All of our solutions are in decreasing-or-equal order.

To turn this into a iterative solution, we need to identify the essential state of the iteration. In this case, it is the slice of integers that we are generating new solutions from (e.g. the state *[5,1]* results in new states *[4,1,1]* and [*3,2,1]*). On change from the recursive solution is that the iterative solution will return the list in a sorted order such as this:

    6
    5 1
    4 2
    4 1 1
    3 3
    3 2 1
    3 1 1 1
    2 2 2
    2 2 1 1
    2 1 1 1 1

To produce the sorted version, we leverage the *sort* package provided by Go. We need to define three methods on our to-be-sorted collection, *Len()*, *Swap()* and *Less()*.

    type items [][]int
    
    func (it items) Len() int { 
        return len(it) 
    }
    
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

Go's package is geared to sort a container from the smallest element to the largest, we want to sort in reverse order so we switch the expression; return *true* if item *i* is larger than *j* and vice versa.

    type intpartition struct {
        agenda items
        item   []int
        err    error
    }

*agenda* is the state of the iteration. *item* is the current solution that *Value()* will return.

    func IntPartition(n int) i.Forward {
        var p intpartition
        item := make([]int, 1)
        item[0] = n
        p.agenda = append(p.agenda, item)
        p.Next()
        return &p
    }

To start the iteration we simply construct a slice from the supplied starting value and call the *Next()* method.

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

An empty *agenda* means the iteration is over.

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

The *Next()* method starts py popping the first item from the *agenda* and assigning it to *item*. Then we break the *item* into two parts, *largest* and *rest* and proceded exacly like the recursive version, except we push new solutions on the *agenda* where we would call the function again. Lastly, we sort the agenda.

        itr := IntPartition(n)
        for ; !itr.AtEnd(); itr.Next() {
            fmt.Println(itr.Value())
        }

Looping through the solution space is trivial.

A version that only prints out sums that contain distinct components, e.g. no number repeats itself, is easy. Simply leverage the *i.Filter* function and write a filter for the serie of integers.

    func distinct(itr i.Iterator) bool {
        item, _ := itr.Value().([]int)
        for i := 1; i < len(item); i++ {
            if item[i-1] == item[i] {
                return false
            }
        }
        return true
    }

Since the list is in a descending order we simply need to check if the current component is equal to the previous component; if so we return false. If we can run through all components without this happening, we return true.

        itr = hoi.Filter(distinct, IntPartition(n))
        for ; !itr.AtEnd(); itr.Next() {
            fmt.Println(itr.Value())
        }

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c5/int_partition.go).