package main

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func binary1(n int) string {
	if n <= 1 {
		return strconv.Itoa(n)
	}
	k := int(n / 2)
	b := int(math.Mod(float64(n), 2))
	return binary1(k) + strconv.Itoa(b)
}

func binary2(n int, retval string) string {
	if n < 1 {
		return retval
	}
	k := int(n / 2)
	b := int(math.Mod(float64(n), 2))
	retval = strconv.Itoa(b) + retval
	return binary2(k, retval)
}

func binary3(n int) string {
	var retval string
	for n >= 1 {
		b := int(math.Mod(float64(n), 2))
		retval = strconv.Itoa(b) + retval
		n = int(n / 2)
	}
	return retval
}

func main() {
	n := 10
	d, r := timeThis(binary1, n)
	fmt.Printf("Binary1: binary(%d) = %q at %v\n", n, r, d)
	fmt.Printf("Binary3: binary(%d) = %q\n", n, binary2(n, ""))
	d, r = timeThis(binary3, n)
	fmt.Printf("Binary3: binary(%d) = %q at %v\n", n, r, d)
}

func timeThis(f func(int) string, p int) (time.Duration, string) {
	start := time.Now()
	res := f(p)
	return time.Since(start), res
}
