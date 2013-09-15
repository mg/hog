package main

import (
	"bufio"
	"fmt"
	"github.com/mg/i"
	"github.com/mg/i/hoi"
	"io"
	"os"
	"strings"
)

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

func (f *fh) SetError(err error) {
	f.err = err
}

func (f *fh) Value() interface{} {
	return f.line
}

func (f *fh) AtEnd() bool {
	return f.err == io.EOF
}

func (f *fh) Next() error {
	if f.line, f.err = f.in.ReadString('\n'); f.err != nil && f.err != io.EOF {
		return f.err
	}
	// chomp
	f.line = strings.TrimSuffix(f.line, "\n")
	return nil
}

func main() {
	hoi.Each(
		Fh(os.Stdin),
		func(itr i.Iterator) bool {
			line, _ := itr.Value().(string)
			fmt.Println("Line: " + line)
			return true
		})
}
