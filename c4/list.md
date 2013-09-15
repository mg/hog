4.4.3: list_iterator()

In this section MJD introduces the *imap* and *igrep* iterators. We've already used them; *imap* was introduced in [4.3.3: Example - A Flat-File Database - part 2]() (available in the iterator library as [i.Map](https://github.com/mg/i/blob/master/map.go)) while *igrep* was introduced in [4.3: Examples]() (available in the iterator library as [i.Filter](https://github.com/mg/i/blob/master/filter.go)). I've won't spend any more time on them here.

A *list_iterator()* transforms a *Perl* array into an iterator. Go's native container is the *slice* so we will build an adaptor for the same purpose here.

One problem we run into is that a slice of some type, say *[]int*, can not be typecasted to a *[]interface{}* type, we would need to create a i*nterface{}* of equal length as the *[]int* and then copy the values over. This is not efficient and not something I want to do.

I have two solutions to this problem. The first one is to write an adapter speficic for the type you want to use.

    type intslice struct {
        slice []int
        pos   int
    }
    
    func IntSlice(vals []int) i.Forward {
        return &intslice{slice: vals}
    }
    
    func (i *intslice) AtEnd() bool {
        return i.pos >= len(i.slice)
    }
    
    func (i *intslice) Next() error {
        i.pos++
        return nil
    }
    
    func (i *intslice) Value() interface{} {
        return i.slice[i.pos]
    }
    
    func (i *intslice) Error() error {
        return nil
    }
    
    func (i *intslice) SetError(err error) {
    }


This is a typical *i.Forward* iterator that simply wraps around the *int* slice you provide.

    itr := IntSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
    itr = hoi.Filter(
        func(itr i.Iterator) bool {
            n, _ := itr.Value().(int)
            return math.Mod(float64(n), 3) == 0
        }, itr)
    
    for !itr.AtEnd() {
        fmt.Println(itr.Value())
        itr.Next()
    }

I've created a script that generates iterators such as this for all the base types in Go. They are all available in the [iterator library](https://github.com/mg/i/icon).

A second strategy is not to wrap a slice but to use Go's support for *variable arguments*. Yes, if you are forced to wrap a slice this will not help you but sometimes all we need is to wrap a list of values that we've written out in the source file.

    type list struct {
        slice []interface{}
        pos   int
    }
    
    func List(vals ...interface{}) i.Forward {
        return &list{slice: vals}
    }
    
    func (i *list) AtEnd() bool {
        return i.pos >= len(i.slice)
    }
    
    func (i *list) Next() error {
        i.pos++
        return nil
    }
    
    func (i *list) Value() interface{} {
        return i.slice[i.pos]
    }
    
    func (i *list) Error() error {
        return nil
    }
    
    func (i *list) SetError(err error) {
    }


This will work for any type, even a mix of different types.

    itr = List(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    itr = hoi.Filter(func(itr i.Iterator) bool {
        n, _ := itr.Value().(int)
        return math.Mod(float64(n), 2) == 0
    }, itr)
    for !itr.AtEnd() {
        fmt.Println(itr.Value())
        itr.Next()
    }

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/list.go). The *variable argument* version is also available in the [iterator library](https://github.com/mg/i/blob/master/icon/list.go).