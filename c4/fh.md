4.3.3: Example - Filehandle Iterators

Go already provides a good package to read a file line by line. It is [bufio.Reader](http://golang.org/pkg/bufio/) so this iterator is simply a wrapper around that reader.

    type fh struct {
        in   *bufio.Reader
        line string
        err  error
    }
    
    func Fh(in io.Reader) i.Forward {
        f := fh{in: bufio.NewReader(in)}
        f.Next()
        return &f
    }
    
    func (f *fh) Error() error {
        if f.err == io.EOF {
            return nil
        }
        return f.err
    }
    
    func (f *fh) Value() interface{} {
        return f.line
    }
    
    func (f *fh) AtEnd() bool {
        return f.err == io.EOF
    }

We need to check if *err* is *io.EOF* or any other error since *io.EOF* is not an error condition but simply an *EOS* indicator. Everything else is pretty much self explanatory.

    func (f *fh) Next() error {
        if f.line, f.err = f.in.ReadString('\n'); f.err != nil && f.err != io.EOF {
            return f.err
        }
        // chomp
        f.line = strings.TrimSuffix(f.line, "\n")
        return nil
    }

The *Next()* function will read the next line from the file and chomp of the new line indicator.

    hoi.Each(
        Fh(os.Stdin),
        func(itr i.Iterator) bool {
            line, _ := itr.Value().(string)
            fmt.Println("Line: " + line)
            return true
        })

Usage is same as before.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/fh.go).    