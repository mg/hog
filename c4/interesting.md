4.3: Examples

The first example is to use the [*Dirwalk*](http://higherordergo.blogspot.com/2013/09/422-dirwalk.html) iterator from the previous post and create some filters on top of it. To do that, we need to create a *Filter* iterator.

    type FilterFunc func(Iterator) bool
    
    type filter struct {
        f   FilterFunc
        itr Forward
    }
    
    func Filter(f FilterFunc, itr Forward) Forward {
        return &filter{f, itr}
    }
    
    func (i *filter) AtEnd() bool {
        for !i.itr.AtEnd() {
            if i.f(i.itr) {
                return false
            }
            i.itr.Next()
        }
        return true
    }
    
    func (i *filter) Next() error {
        return i.itr.Next()
    }
    
    func (i *filter) Value() interface{} {
        return i.itr.Value()
    }
    
    func (i *filter) Error() error {
        return i.itr.Error()
    }

The *Filter* iterator is a struct that wraps around a *Forward* iteartor. For most part it simply forwards any method calls to the enclosed iterator. The filtering function is provided by the *AtEnd()* method. If the wrapped iterator contains further elements, it simply loops until it finds an element that fulfills the requirements of the user supplied *FilterFunc*. At first, this seems to break the idempotency requirement. But if we call *AtEnd()* it will loop until it hits the end of stream or finds an element that is not filtered. If we then call *AtEnd()* again without calling *Next()*, it simply finds the same element again and returns at the same position.

I've put this *Filter* iterator into the [i package](https://github.com/mg/i/blob/master/filter.go).

Now the problem is simply to create a *FilterFunc* that performes the filtering we require.

    func hasInName(val string) i.FilterFunc {
        val = strings.ToUpper(val)
        return func(itr i.Iterator) bool {
            filename, _ := itr.Value().(string)
            return strings.Contains(strings.ToUpper(filename), val)
        }
    }
    
    func not(f i.FilterFunc) i.FilterFunc {
        return func(itr i.Iterator) bool {
            return !f(itr)
        }
    }

*hasInName* performes a case insensitive search, matching *val* to filenames that *Dirwalk* returns. *not* simply negates the result of any *FilterFunc*.

    i.Each(
        i.Filter(not(hasInName("example")), Dirwalk(os.Args[1])),
        func(itr i.Iterator) bool {
            fmt.Println(itr.Value())
            return true
        })

This constructs loops through all filenames under *os.Args[1]* and returns any that don't contain the string *example*.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/interesting.go).