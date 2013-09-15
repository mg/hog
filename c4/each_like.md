4.6.2: An Iterator with an each-Like Interface

The concept behind *each-Like* is to run a *Map* operation transforming a value and returning it along with the original to the caller.

This is a *Map* iterator in almost every way.

    type eachlike struct {
        itr i.Forward
        f   hoi.MapFunc
    }
    
    func EachLike(f i.MapFunc, itr i.Forward) *eachlike {
        return &eachlike{f: f, itr: itr}
    }

The constructor takes a *f.MapFunc* function and a *i.Forward* iterator to work on.

    func (i *eachlike) Value() interface{} {
        ret := make([]interface{}, 2)
        ret[0] = i.itr.Value()
        ret[1] = i.f(i.itr)
        return ret
    }
    
    func (i *eachlike) ValuePair() (interface{}, interface{}) {
        res := i.Value().([]interface{})
        return res[0], res[1]
    }

The *Value()* method returns an *interface{}* value that contains both the original value and the transformed value. The *ValuePair()* method unpacks the *interface{}* value into a *interface{}* slice.

    func (i *eachlike) Error() error {
        return i.itr.Error()
    }
    
    func (i *eachlike) Next() error {
        return i.itr.Next()
    }
    
    func (i *eachlike) AtEnd() bool {
        return i.itr.AtEnd()
    }

These functions simply forward to the wrapped iterator.

    func multiply(n int) i.MapFunc {
        return func(itr i.Iterator) interface{} {
            return n * itr.Value().(int)
        }
    }

This is a simple closure that creates a *i.MapFunc* that multiplies two integers together.

    itr := EachLike(multiply(10), icon.List(1, 2, 3, 4, 5, 6, 7, 8, 9))
    for ; !itr.AtEnd(); itr.Next() {
        val, res := itr.ValuePair()
        fmt.Printf("Value: %v, result: %v\n", val, res)
    }

To use it we construct the iterator from the *i.MapFunc* and the *i.Forward* iterator. We then loop over it and print out every value, result pair.

The rational behind this is to exploit a feature in *Perl* where the caller can indicate to a function whether it expects a single value or a list value returned. The *EachLike* would return the *transformed value* in the single context, and the (*transformed value, original value)* in the list context. This is not possible in Go. And the core functionality of returning the transformed value with the origianl can easily be achived with the regular *i.Map* iterator and a *i.MapFunc* function that simply returns a slice with both values.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/each_like.go).