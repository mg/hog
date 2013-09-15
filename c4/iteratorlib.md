Concerning the i library

During the last few days, the spinof library that is the [i library](https://github.com/mg/i) has seen many changes. I've refactored it into subpackages (*i*, *hoi*, *icon*, *iadapt*, *ityped*, *itk*), added tests and examples, and a first draft of a documentation is up on [godoc.org](http://godoc.org/github.com/mg/i). Over the next days I will add to the documentation and examples.

The [*hoi*](https://github.com/mg/i/tree/master/hoi) package is the *Higher Order Iterator* package, as such it contains mainy useful higher order functions to use with iterators. As of today it contains the following iterators:

* [append](https://github.com/mg/i/tree/master/hoi/append.go)
* [count](https://github.com/mg/i/tree/master/hoi/count.go)
* [cycle](https://github.com/mg/i/tree/master/hoi/cycle.go)
* [filter](https://github.com/mg/i/tree/master/hoi/filter.go)
* [foldl](https://github.com/mg/i/tree/master/hoi/fold.go)
* [foldr](https://github.com/mg/i/tree/master/hoi/fold.go)
* [map](https://github.com/mg/i/tree/master/hoi/map.go)
* [repeat](https://github.com/mg/i/tree/master/hoi/repeat.go)
* [reverse](https://github.com/mg/i/tree/master/hoi/reverse.go)
* [sample](https://github.com/mg/i/tree/master/hoi/sample.go)
* [shuffle](https://github.com/mg/i/tree/master/hoi/sample.go)
* [slice](https://github.com/mg/i/tree/master/hoi/slice.go)
* [zip](https://github.com/mg/i/tree/master/hoi/zip.go)
* [zip longest](https://github.com/mg/i/tree/master/hoi/zip.go)

Documentation for those functions is up on [GoDoc](http://godoc.org/github.com/mg/i/hoi).

As a result of the refactoring, I've had to go back over some of the examples in chapter 4 and change the code in the *hog* repository. I've also gone through the past blogposts and updated the links to files in the *i* repository. 