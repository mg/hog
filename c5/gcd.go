package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

func gcd(m, n int) int {
	if n == 0 {
		return m
	}
	m = int(math.Mod(float64(m), float64(n)))
	return gcd(n, m)
}

func gcd_notail(m, n int) int {
	for n != 0 {
		m, n = n, int(math.Mod(float64(m), float64(n)))
	}
	return m
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage %s M N\n", os.Args[0])
		os.Exit(0)
	}
	m, _ := strconv.Atoi(os.Args[1])
	n, _ := strconv.Atoi(os.Args[2])
	d, r := timeThis(gcd, m, n)
	fmt.Printf("GCD Recursive of %d and %d: %d at %v\n", m, n, r, d)
	d, r = timeThis(gcd_notail, m, n)
	fmt.Printf("GCD NoTail of %d and %d: %d at %v\n", m, n, r, d)
}

func timeThis(f func(int, int) int, p1, p2 int) (time.Duration, int) {
	start := time.Now()
	res := f(p1, p2)
	return time.Since(start), res
}
