package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

type (
	Keyfunc  func(...interface{}) string
	Memfunc  func(...interface{}) interface{}
	Algofunc func(int) int
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

// Construct a fib calculator. If no function is passed, construct
// it to call it self recursively
func fib(mf Memfunc) Algofunc {

	// Fib calculator that calls either an memonized function or itself
	f := func(month int) int {
		if month < 2 {
			return 1
		}
		return toInt(mf(month-1)) + toInt(mf(month-2))
	}

	// If no memoize was supplied, create a dummy that simply calls
	// the fib calculator
	if mf == nil {
		mf = func(v ...interface{}) interface{} {
			return f(toInt(v[0]))
		}
	}

	return f
}

func main() {
	fibnum := 30
	if len(os.Args) > 1 {
		fibnum, _ = strconv.Atoi(os.Args[1])
	}

	// Regular fib calculation
	rfib := fib(nil)
	duration, res := timeThis(rfib, fibnum)
	fmt.Printf("Raw fib %d: %d at %s\n", fibnum, res, duration)

	// Construct memoize fib calculator
	mfib := func() Algofunc {
		// Forward decleration of fib calculator
		var calcfib Algofunc

		// Memoizer, uses yet to be defined but already declared fib calculator
		memfib := memoize(func(values ...interface{}) interface{} {
			return calcfib(toInt(values[0]))
		}, nil)

		// Construct the calculator, use the memonizer defined above
		calcfib = fib(memfib)

		// Return a function that will do the neccessary typecasting
		return func(month int) int {
			return toInt(memfib(month))
		}
	}()

	// Memoized fib calculation
	duration, res = timeThis(mfib, fibnum)
	fmt.Printf("Memoized fib %d: %d at %s\n", fibnum, res, duration)
}

// Utility functions
func toInt(v interface{}) int {
	if value, ok := v.(int); ok {
		return value
	}
	panic(fmt.Sprintf("Not an int value: %v", v))
}

func timeThis(f Algofunc, param int) (time.Duration, int) {
	start := time.Now()
	res := f(param)
	return time.Since(start), res
}
