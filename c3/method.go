package main

import (
	"bytes"
	"fmt"
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

type Object struct {
}

func (o *Object) Say(p int) int {
	time.Sleep(1 * time.Second)
	return p
}

func main() {
	var obj Object
	duration, res := timeThis(obj.Say, 1)
	fmt.Printf("Object says: %d, takes %s to do it.\n", res, duration)

	// Construct memoize partitioner
	msay := func(f Algofunc) Algofunc {
		memsay := memoize(func(values ...interface{}) interface{} {
			return f(toInt(values[0]))
		}, nil)

		// Return a function that will do the neccessary typecasting
		return func(val int) int {
			return toInt(memsay(val))
		}
	}(obj.Say)

	duration, res = timeThis(msay, 1)
	fmt.Printf("Memoized object says: %d, takes %s to do it.\n", res, duration)
	duration, res = timeThis(msay, 1)
	fmt.Printf("Memoized object says: %d, takes %s to do it.\n", res, duration)
	duration, res = timeThis(msay, 2)
	fmt.Printf("Memoized object says: %d, takes %s to do it.\n", res, duration)
	duration, res = timeThis(msay, 2)
	fmt.Printf("Memoized object says: %d, takes %s to do it.\n", res, duration)
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
