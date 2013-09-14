4.6.1: Using foreach to Loop Over More Than One Array

Zipping lists is to take some lists A, B, C .. and produce a new list containing the elements (a1, b1, c1..), (a2, b2, c2, ..) .. A key decision here is does the new list stop producing new items once the shortest list is exhausted, or does it continue to produce elements until the longest list is finished, and fill in *nil* values for the shorter lists.

    type stop int
    
    const (
        StopAtMin stop = iota
        StopAtMax
    )
    
    type eacharray struct {
        itrs   []i.Forward
        err    error
        stopAt stop
    }
    
    func EachArray(stopAt stop, itrs ...i.Forward) i.Forward {
        return &eacharray{itrs: itrs, stopAt: stopAt, err: nil}
    }

The constructor takes an enum to control the behaviour vis-a-vis the length of individual streams, and a list of iterators.

    func (i *eacharray) Value() interface{} {
        ret := make([]interface{}, len(i.itrs))
        for idx, itr := range i.itrs {
            if !itr.AtEnd() {
                ret[idx] = itr.Value()
            } else {
                ret[idx] = nil
            }
        }
        return ret
    }
    
    func (i *eacharray) Error() error {
        return i.err
    }
    
    func ( *eacharray) SetError(err error) {
        i.err = err
    }


The *Value()* function loops over all iterators, collecting a value from each of them. If the iterator is at end, we use a nil value for that iterator.

    func (i *eacharray) Next() error {
        for _, itr := range i.itrs {
            if !itr.AtEnd() {
                err := itr.Next()
                if err != nil {
                    i.err = err
                    break
                }
            }
        }
        return i.err
    }
    
    func (i *eacharray) AtEnd() bool {
        atEndCount := 0
        for _, itr := range i.itrs {
            if itr.AtEnd() {
                if i.stopAt == StopAtMin {
                    return true
                }
                atEndCount++
            }
        }
        return atEndCount == len(i.itrs)
    }

The *Next()* method loops over all iterators, calling *Next()* on each iterator that is not at end yet. The *AtEnd()* method loops over all the iterator and checks each one for *EOS*. If *stopAt* is *StopAtMin*, we return true on first *EOS* we encoundter, otherwise we return false until all iterators return *EOS*.

    itr := EachArray(
        StopAtMax,
        i.List(1, 2, 3, 4, 5, 6),
        i.List(6.4, 7.1, 8.2, 9.9),
        i.List("A", "B", "C", "D", "E"))
    for ; !itr.AtEnd(); itr.Next() {
        fmt.Println(itr.Value())
    }

Usage is simple and there is no need for all the lists to be of the same type.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/each_array.go). It is also available in the iterator library as [Zip](https://github.com/mg/i/blob/master/zip.go).