4.4.4: append()

Go has an append function, it simply appends two slices and returns the resulting slice. An append iterator is similar, the major difference being that you can start to work on the joined list without actually going through the process of copying them together. If the lists are long, this is a good strategy.

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

*Append* simply runs through each iterator until it hits the end, then it switches to the next one. Once the last iterator is exhausted, *Append* quits.

The usage is very simple:

    for itr := Append(hoi.List(1, 2, 3), igen.Range(10, 20)); !itr.AtEnd(); itr.Next() {
        fmt.Println(itr.Value())
    }

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/append.go). It is also available in the [iterator library](https://github.com/mg/i/blob/master/hoi/append.go).
