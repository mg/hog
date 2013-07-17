package main

import (
	"fmt"
	"os"
)

type (
	ComputeType interface{}
	FileFunc    func(f *os.File) ComputeType
	DirFunc     func(d *os.File, results []ComputeType) ComputeType
)

func empty() ComputeType {
	var e ComputeType
	return e
}

func dirwalk(name string, f FileFunc, d DirFunc) ComputeType {
	file, err := os.Open(name)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err.Error())
	}
	if stat.IsDir() {
		files, err := file.Readdirnames(0)
		if err != nil {
			panic(err.Error())
		}
		name = name + string(os.PathSeparator)
		results := make([]ComputeType, 0)
		for _, subfile := range files {
			results = append(results, dirwalk(name+subfile, f, d))
		}
		if d == nil {
			return empty()
		}
		return d(file, results)
	}
	if f == nil {
		return empty()
	}
	return f(file)
}

func int64value(val interface{}) int64 {
	if v, ok := val.(int64); ok {
		return v
	}
	panic(fmt.Sprintf("Unexpected type: %v", val))
}

func filefunc(f *os.File) ComputeType {
	stat, err := f.Stat()
	if err != nil {
		panic(err.Error())
	}
	return ComputeType(stat.Size())
}

func dirfunc(d *os.File, results []ComputeType) ComputeType {
	var total int64
	for _, v := range results {
		total += int64value(v)
	}
	return total
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage %s NAME\n", os.Args[0])
		os.Exit(0)
	}
	fmt.Println(int64value(dirwalk(os.Args[1], filefunc, dirfunc)))
}
