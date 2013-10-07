package main

import (
	"container/list"
	"fmt"
	"time"
)

func fib1(n int) int {
	if n < 2 {
		return n
	}
	return fib1(n-2) + fib1(n-1)
}

func fib2(n int) int {
	for {
		if n < 2 {
			return n
		}
		s1 := fib2(n - 2)
		s2 := fib1(n - 1)
		return s1 + s2
	}
}

type state3 struct {
	BRANCH, s1, s2, n int
}

func fib3(n int) int {
	var s1, s2, retval int
	BRANCH := 0
	var STACK list.List
	for {
		if n < 2 {
			retval = n
		} else {
			if BRANCH == 0 {
				STACK.PushBack(&state3{BRANCH, s1, s2, n})
				n -= 2
				BRANCH = 0
				continue
			} else if BRANCH == 1 {
				s1 = retval
				STACK.PushBack(&state3{BRANCH, s1, s2, n})
				n -= 1
				BRANCH = 0
				continue
			} else if BRANCH == 2 {
				s2 = retval
				retval = s1 + s2
			}
		}
		if STACK.Len() == 0 {
			return retval
		}
		s := STACK.Back().Value.(*state3)
		STACK.Remove(STACK.Back())
		BRANCH, s1, s2, n = s.BRANCH, s.s1, s.s2, s.n
		BRANCH++
	}
}

type state4 struct {
	BRANCH, s1, n int
}

func fib4(n int) int {
	var s1, retval int
	BRANCH := 0
	var STACK list.List
	for {
		if n < 2 {
			retval = n
		} else {
			if BRANCH == 0 {
				for n >= 2 {
					STACK.PushBack(&state4{BRANCH, 0, n})
					n -= 2
				}
				retval = n
			} else if BRANCH == 1 {
				s1 = retval
				STACK.PushBack(&state4{BRANCH, retval, n})
				n -= 1
				BRANCH = 0
				continue
			} else if BRANCH == 2 {
				retval += s1
			}
		}
		if STACK.Len() == 0 {
			return retval
		}
		s := STACK.Back().Value.(*state4)
		STACK.Remove(STACK.Back())
		BRANCH, s1, n = s.BRANCH, s.s1, s.n
		BRANCH++
	}
}

func fib5(n int) int {
	var s1, retval int
	BRANCH := 0
	var STACK list.List
	for {
		if n < 2 {
			retval = n
		} else {
			if BRANCH == 0 {
				for n >= 2 {
					STACK.PushBack(&state4{1, 0, n})
					n -= 1
				}
				retval = n
			} else if BRANCH == 1 {
				s1 = retval
				STACK.PushBack(&state4{2, retval, n})
				n -= 2
				BRANCH = 0
				continue
			} else if BRANCH == 2 {
				retval += s1
			}
		}
		if STACK.Len() == 0 {
			return retval
		}
		s := STACK.Back().Value.(*state4)
		STACK.Remove(STACK.Back())
		BRANCH, s1, n = s.BRANCH, s.s1, s.n
	}
}

func main() {
	n := 20
	d, r := timeThis(fib1, n)
	fmt.Printf("Fib1: fib(%d) = %d at %v\n", n, r, d)
	fmt.Printf("Fib2: fib(%d) = %d\n", n, fib2(n))
	fmt.Printf("Fib3: fib(%d) = %d\n", n, fib3(n))
	fmt.Printf("Fib4: fib(%d) = %d\n", n, fib4(n))
	d, r = timeThis(fib5, n)
	fmt.Printf("Fib5: fib(%d) = %d at %v\n", n, r, d)
}

func timeThis(f func(int) int, p int) (time.Duration, int) {
	start := time.Now()
	res := f(p)
	return time.Since(start), res
}
