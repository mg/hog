package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(binary(2))
	fmt.Println(binary(3))
	fmt.Println(binary(37))
	fmt.Println(binary(30001))
	fmt.Println(binary(1023))
}

func binary(n int) string {
	if n == 0 || n == 1 {
		return strconv.Itoa(n)
	}
	return binary(n/2) + strconv.Itoa(n%2)
}
