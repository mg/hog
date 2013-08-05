The Towers of Hanoi is an old puzzle and quite the popular example for recursion in many languages. The legend is that the world will end once you move 64 disks.

This is type of the callback that performs the Hanoi move. An *interface{}* could be used here but since no state is being maintained a function type is a better choice.

CODE

The Hanoi function. Moves n pegs from start to end using 3 pegs.

CODE

A basic moving function. Simply prints out the actual move (e.g. Move disk 2 from "A" to "C").

CODE

Kicks of the process for 4 disks.

CODE

This function constructs a *MoveFunc* function that wraps around an actual move function. The idea here is to check if certain invariants hold before each move.

CODE

Again, kicking of the process for  4 disks, now using the wrapper to construct the invariant function around the mover.

CODE

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c1/hanoi.go).