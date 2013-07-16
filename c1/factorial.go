package main

import (
	"fmt"
)

func main() {
	fmt.Println(factorial(0))
	fmt.Println(factorial(1))
	fmt.Println(factorial(2))
	fmt.Println(factorial(3))
	fmt.Println(factorial(4))
	fmt.Println(factorial(5))
	fmt.Println(factorial(6))
	fmt.Println(factorial(7))
	fmt.Println(factorial(8))
	fmt.Println(factorial(9))
	fmt.Println(factorial(10))
}

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return factorial(n-1) * n
}
