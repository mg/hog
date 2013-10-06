package main

import (
	"fmt"
)

func factorial1(n int) int {
	if n == 0 {
		return 1
	}
	return factorial1(n-1) * n
}

func factorial2(n, product int) int {
	if n == 0 {
		return product
	}
	return factorial2(n-1, n*product)
}

func factorial3(n int) int {
	product := 1
	for n > 0 {
		product *= n
		n--
	}
	return product
}

func main() {
	n := 10
	fmt.Println(factorial1(n))
	fmt.Println(factorial2(n, 1))
	fmt.Println(factorial3(n))
}
