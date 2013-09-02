4.2.1: A Trivial Iterator - Take three

Iterators are about movement across data, and what differentiates one class of iterators from another is the type of movement it allows. The *Upto* iterator is essentially a *Forward* iterator, it provides a way to move forward, to get the value at the current location, and a way to check if we have reached the end of the data.

    type upto struct {
        m, n int
    }
    
    func Upto(m, n int) *upto {
        return &upto{m: m, n: n}
    }
    
    func (i *upto) Next() {
        i.m++
    }
    
    func (i *upto) AtEnd() bool {
        return i.m >= i.n
    }
    
    func (i *upto) Value() int {
        return i.m
    }

As before, we use a struct to capture the state of the iterator. Three methods provide the functionality of the iterator: *Value()* gives us the value at the current location, *Next()* moves the iterator to the next location, and *AtEnd()* allows us to check if we've reached *EOS*. At the cost of a bit more of boilerplate, we now have an iterator with a well defined interface and no hidden behaviours.

    func main() {
        for i := Upto(3, 5); !i.AtEnd(); i.Next() {
            fmt.Println(i.Value())
        }
    }

This design fits very well in the *initialization*, *check*, *increment* design of the standard *for* loop.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/upto_take3.go).

## Formal design of iterators

Going forward this is the design that I will use, and I want to formalize it a bit more, especially with respect to types of iterators. These are all [external iterators](http://journal.stuffwithstuff.com/2013/01/13/iteration-inside-and-out/) classified according to the movements they provide and how they define the boundaries of the data sets they provide access too.

    Iterator interface {
        Value() interface{}
        Error() error
    }

The basic interface is the *Iterator*. It simply provides two idempotent ways to access the value and last error raised.

    Forward interface {
        Iterator
        Next() error
        AtEnd() bool
    }

The *Forward* is a extension of the *Iterator*. It provides previously seen *Next()* method and *AtEnd()* method. This is a iterator to provide access to any lazily evaluated, potentially infinite stream of data. The stream must have a way to move forward, but whether it has any end at all is unanswered.


    BiDirectional interface {
        Forward
        Prev() error
        AtStart() bool
    }

The *BiDirectional* is a *Forward* that allows us to move back in the stream with *Prev()* and check if we reached the start of the stream with *AtStart()*. *AtStart()* should be idempotent.

    BoundedAtStart interface {
        Forward
        First() error
    }

*BoundedAtStart* is a *Forward* iterator that has a well defined beginning that it can be reset to in an efficient manner.

    Bounded interface {
        BiDirectional
        First() error
        Last() error
    }

*Bounded* is a *BiDirectional* with the ability to jump to the end of the data set quickly. This is e.g. a file on disk where we can quickly seek to the end of the file.

    RandomAccess interface {
        Bounded
        Goto(int) error
        Len() int
    }

The last class of iterators is the *RandomAccess*. It is a *Bounded* iterator with two more functions. *Goto()* provides a quick way to jump into any point in the dataset while *Len()* gives us a quick count of distinct elements in the data set.

To understand the difference between *Bounded* and *RandomAccess* think of a text file on the disk. If we want to iterate over the individual bytes in the file, *RandomAccess* will provide us all the tools needed with acceptable performance. If we want to iterate over the individual lines in the file, we have to do with a *Bounded*, since there is no efficient way to either count the individual lines or position them in the file.

This is the design I will use moving forward and to that end I've set up a new project on [GitHub](https://github.com/mg/i). Right now it only contains the definitions of [iterator interfaces](https://github.com/mg/i/blob/master/iterator.go) and a *RandomAccess* version of the *Upto* iterator, now named [*Range*](https://github.com/mg/i/blob/master/iutil/range.go), but I will add more higher order algorithms and iterators to it as I progress through the chapter.