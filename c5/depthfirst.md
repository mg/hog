5.3: A Generic Search Iterator

We can encapsulate the core strategy of the previous example in a code structure that we can then reuse to create iterators such as the *IntPartition* iterator. The idea is to extract the code that revolves around the *agenda* and the iteration. To use this very generic iteration, we need to supply a function that will generate new items for the *agenda* from the current item.

    type ChildrenFunc func(interface{}) []interface{}
    
    type depthfirst struct {
        children ChildrenFunc
        agenda   []interface{}
        node     interface{}
        err      error
    }

The iterator consists of the *agenda*, the current value in the *node* item, and a function to generate new items.

    func Dfs(root interface{}, children ChildrenFunc) i.Forward {
        dfs := depthfirst{children: children}
        dfs.agenda = append(dfs.agenda, root)
        dfs.Next()
        return &dfs
    }

To create the iterator we set the agenda to the root item, call *Next()* to generate the first value and return.

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

*node* is the current value of the iteration and the iteration is over once the *agenda* is empty.

    func (dfs *depthfirst) Next() error {
        if dfs.AtEnd() {
            dfs.err = fmt.Errorf("Next: Beyond end")
            return dfs.err
        }
        dfs.node, dfs.agenda = dfs.agenda[0], dfs.agenda[1:len(dfs.agenda)]
        dfs.agenda = append(dfs.agenda, dfs.children(dfs.node)...)
        return nil
    }

The iteration should be familiar. Pop the head of the *agenda* and call the function we received to generate the next batch of items from the current item and append them to the *agenda*.

The *IntPartition* iterator is now a short piece of code that focuses on solving only the problem of partitioning an integer; all the iteration bookkeeping is gone. The only new thing here is, rather than adding items to the *agenda* for further iteration, we create a temporary slice and add our generated nodes to it. This slice is the list of children we return to the *Dfs* iterator.

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

Now its 15 lines of code rather than 50 lines of code. Yeah for iterator building blocks.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c5/depthfirst.go). The *Dfs* iterator is also available in the [iterator library](https://github.com/mg/i/blob/master/dfs.go).