package main

import (
	"fmt"
	"github.com/mg/i"
	"os"
	"strings"
)

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

func (dw *dirwalk) Value() interface{} {
	return dw.cur
}

func (dw *dirwalk) Error() error {
	return dw.err
}

func (dw *dirwalk) SetError(err error) {
	dw.err = err
}

func (dw *dirwalk) AtEnd() bool {
	return len(dw.queue) == 0
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage %s NAME\n", os.Args[0])
		os.Exit(0)
	}

	i.Each(Dirwalk(os.Args[1]), func(itr i.Iterator) bool {
		fmt.Println(itr.Value())
		return true
	})
}
