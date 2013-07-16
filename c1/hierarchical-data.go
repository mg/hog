package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage %s NAME\n", os.Args[0])
		os.Exit(0)
	}
	fmt.Println(totalsize(os.Args[1]))
}

func totalsize(name string) int64 {
	f, err := os.Open(name)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		panic(err.Error())
	}
	size := stat.Size()
	if stat.IsDir() {
		files, err := f.Readdirnames(0)
		if err != nil {
			panic(err.Error())
		}
		name = name + string(os.PathSeparator)
		for _, file := range files {
			size += totalsize(name + file)
		}
	}
	return size
}
