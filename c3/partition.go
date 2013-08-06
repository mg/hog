package main

import (
	"bytes"
	"fmt"
	"time"
)

type (
	Keyfunc  func(...interface{}) string
	Memfunc  func(...interface{}) interface{}
	Algofunc func(int, []int) []int
)

func memoize(f Memfunc, kf Keyfunc) Memfunc {
	if kf == nil {
		kf = func(v ...interface{}) string {
			// Very naive keygen func, simply use string representation
			// of every param
			var buffer bytes.Buffer
			for i := 0; i < len(v); i++ {
				buffer.WriteString(fmt.Sprint(v[0]))
				buffer.WriteString(",")
			}
			return buffer.String()
		}
	}

	cache := make(map[string]interface{})
	return func(v ...interface{}) interface{} {
		cachekey := kf(v)
		if r, ok := cache[cachekey]; ok {
			// Hit, return previous result
			return r
		}
		// Miss, call calculator
		r := f(v...)
		cache[cachekey] = r
		return r
	}
}

func findshare(mf Memfunc) Algofunc {

	f := func(target int, treasures []int) []int {
		if target == 0 {
			return make([]int, 0)
		}
		if target < 0 || len(treasures) == 0 {
			return nil
		}
		first, rest := treasures[0:1], treasures[1:]
		solution := toIntSlice(mf(target-first[0], rest))
		if solution != nil {
			return append(first, solution...)
		}
		return toIntSlice(mf(target, rest))
	}

	if mf == nil {
		mf = func(v ...interface{}) interface{} {
			return f(toInt(v[0]), toIntSlice(v[1]))
		}
	}
	return f
}

func main() {
	treasure := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	target := 209

	rawfindshare := findshare(nil)
	duration, res := timeThis(rawfindshare, target, treasure)
	fmt.Printf("Raw partition %v into %d: %d at %s\n", treasure, target, res, duration)

	// Construct memoize partitioner
	mfindshare := func() Algofunc {
		// Forward decleration of calculator
		var calcshare Algofunc

		// Memoizer
		memfindshare := memoize(func(values ...interface{}) interface{} {
			return calcshare(toInt(values[0]), toIntSlice(values[1]))
		}, nil)

		// Construct the calculator, use the memonizer defined above
		calcshare = findshare(memfindshare)

		// Return a function that will do the neccessary typecasting
		return func(target int, treasures []int) []int {
			return toIntSlice(memfindshare(target, treasures))
		}
	}()

	duration, res = timeThis(mfindshare, target, treasure)
	fmt.Printf("Memoized partition %v into %d: %d at %s\n", treasure, target, res, duration)

}

// Utility functions
func toInt(v interface{}) int {
	if value, ok := v.(int); ok {
		return value
	}
	panic(fmt.Sprintf("Not an int value: %v", v))
}

func toIntSlice(v interface{}) []int {
	if value, ok := v.([]int); ok {
		return value
	}
	panic(fmt.Sprintf("Not an int slice value: %v", v))
}

func timeThis(f Algofunc, param1 int, param2 []int) (time.Duration, []int) {
	start := time.Now()
	res := f(param1, param2)
	return time.Since(start), res
}
