5.1: The Partition Problem Revisited

The recurisve partition algorithm form [1.8](http://higherordergo.blogspot.com/2013/07/18-partitioning.html) solved for a single solution when trying to search a list of treasures for a share of certain value. The iterator version will return single solution on each iteration, successive iterations will return new solutions until all possible solutions have been found when we reach AtEnd().

    type entry struct {
        target      int
        pool, share []int
    }
    
    type partition struct {
        todo  *list.List
        err   error
        share []int
    }

The *entry* item contains the current target, the *pool* of items to consider, and the *share* that make up a possible solution. The *todo* is a list of entries to inspect.

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

To construct the iterator we create the first *entry* item for the *todo* list from the *target* and *treasures* parameters. If *target* is zero we put a empty list for the *pool* so the iterator will return right after the first iteration. Then we execute the first iteration.

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

An empty *todo* list indicates that we are at the end of the iteration. The *share* is the current value.

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

We start by popping the first entry off the *todo* list and splitting the pool into a *first* item and the *rest*. If the *rest* contains items there are more possible solutions to consider so we put it on the *todo* list for future inspection.

If the *target* is equal to the *first* element of the *pool* we've found a solution and return. Otherwise if the *target* is larger than the *first* element and the *rest* contains more elements there could be a solution or two still in the *rest* list, so we put it on the *todo* list to inspect.

If we reach this point we did not find a solution so we pop the next item of the *todo* list and start again, or return if it is empty.

    itr := Partition(target, treasures)
    for ; !itr.AtEnd(); itr.Next() {
        fmt.Println(itr.Value())
    }

Construct the iterator and loop through it to find all solutions.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c5/partition.go).