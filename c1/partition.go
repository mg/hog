package main

import (
	"fmt"
)

func findshare(target int, treasures []int) []int {
	if target == 0 {
		return make([]int, 0)
	}
	if target < 0 || len(treasures) == 0 {
		return nil
	}
	first, rest := treasures[0:1], treasures[1:]
	solution := findshare(target-first[0], rest)
	if solution != nil {
		return append(first, solution...)
	}
	return findshare(target, rest)
}

func main() {
	treasure := []int{5, 2, 4, 8, 1}
	total := 0
	for _, v := range treasure {
		total += v
	}
	share := findshare(total/2, treasure)
	fmt.Print(share)
}
