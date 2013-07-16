package main

import (
	"fmt"
)

type (
	MoveFunc func(n int, start, end string)
)

func main() {
	hanoi(4, "A", "B", "C", checkmove(4, "A", move))
}

func hanoi(n int, pegStart, pegEnd, pegExtra string, mover MoveFunc) {
	if n == 1 {
		mover(1, pegStart, pegEnd)
	} else {
		hanoi(n-1, pegStart, pegExtra, pegEnd, mover)
		mover(n, pegStart, pegEnd)
		hanoi(n-1, pegExtra, pegEnd, pegStart, mover)
	}
}

func move(n int, start, end string) {
	fmt.Printf("Move disk %d from %q to %q.\n", n, start, end)
}

func checkmove(n int, peg string, move MoveFunc) MoveFunc {
	position := make([]string, n+1)
	for i := 1; i < len(position); i++ {
		position[i] = peg
	}
	return func(n int, start, end string) {
		if n < 1 || n > len(position)-1 {
			panic(fmt.Sprintf("Bad disk number %d, should be 1..%d", n, len(position)-1))
		}
		if position[n] != start {
			panic(fmt.Sprintf("Tried to move disk %d from %q, but it is on peg %q", n, start, position[n]))
		}
		for i := 1; i < n-1; i++ {
			if position[i] == start {
				panic(fmt.Sprintf("Can't move disk %n from %q because %n is on top of it", n, start, i))
			} else if position[i] == end {
				panic(fmt.Sprintf("Can't move disk %n to %q becasue %n is already there", n, end, i))
			}
		}
		move(n, start, end)
		position[n] = end
	}
}
