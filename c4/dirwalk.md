4.2.2: dir_walk()

The recursive [dirwalk](http://higherordergo.blogspot.com/2013/07/15-applications-and-variations-of.html) can be rewritten as an *Forward* iterator. 

    type dirwalk struct {
        cur   string
        queue []string
        err   error
    }
    
    func Dirwalk(filename string) i.Forward {
        // remove trailing /
        strings.TrimSuffix(filename, "/")
        
        // construct and initialize
        var dw dirwalk
        dw.queue = []string{filename}
        dw.Next()
        return &dw
    }

The *dirwalk* structure contains a *queue* of filenames to be processed. The *cur* variable holds the filename that will be returned when *Value()* is called. *Dirwalk* simply constructs this structure and returns it, encapsulated as a *Forward* iterator.

    func (dw *dirwalk) Value() interface{} {
        return dw.cur
    }
    
    func (dw *dirwalk) Error() error {
        return dw.err
    }
    
    func (dw *dirwalk) AtEnd() bool {
        return len(dw.queue) == 0
    }

The length of the *queue* represents the state of the iteration; once it hits zero we are done.

    func (dw *dirwalk) Next() error {
        // pop head from queue
        dw.cur, dw.queue = dw.queue[0], append(dw.queue[1:])
    
        // open file
        var file *os.File
        if file, dw.err = os.Open(dw.cur); dw.err != nil {
            return dw.err
        }
        defer file.Close()
    
        // stat file
        var stat os.FileInfo
        if stat, dw.err = file.Stat(); dw.err != nil {
            return dw.err
        }
    
        if stat.IsDir() {
            // read files in directory
            var files []string
            if files, dw.err = file.Readdirnames(0); dw.err != nil {
                return dw.err
            }
    
            // add files in directory to queue
            for _, subfile := range files {
                dw.queue = append(dw.queue, dw.cur+string(os.PathSeparator)+subfile)
            }
        }
        return dw.err
    }

*Next* is the engine of the iterator. It pops the queue for the current filename and then processes it. If it is a simple file we are done and return. If it is a directory we read through it and put the filenames in it on the queue.

    for itr := Dirwalk(os.Args[1]); !itr.AtEnd(); itr.Next() {
        fmt.Println(itr.Value())
    }

Using the iterator is very simple, simply construct it and loop through it. 

This provides an opportunity to write the first higher order function that depends on the *Iterator* interface. The loop can be written as an *Each* function as follows:

    type EachFunc func(i Iterator) bool
    
    func Each(i Forward, e EachFunc) {
        for ; !i.AtEnd(); i.Next() {
            if !e(i) {
                break
            }
        }
    }

It accepts a *Forward* iterator and a function of type EachFunc, loops through the data stream and calls the function on each *Iterator*. If the function returns false, the loop breaks.

Now we can use the *Dirwalk()* iterator as follows:

    hoi.Each(Dirwalk(os.Args[1]), func(itr i.Iterator) bool {
        fmt.Println(itr.Value())
        return true
    })

Yes, this is very trivial and of dubious value, but it is a start.

You can get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/dirwalk.go) and you can get the source for the Each function in the [i package](https://github.com/mg/i/blob/master/hoi/each.go).