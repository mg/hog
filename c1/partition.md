1.8: Partitioning

The last problem in chapter 1 is intended to warn against the follies of reckless recursion. It works well for a small problem space but quickly blows up once we increase it. The problem is to divide a treasure (e.g. consisting of items valued at $5, $2, $4, $8 and $1) into shares. 

*findshare* accepts a target and a list of valuables and returns a solution. If the target is 0 returns a list with zero elements and if target is less than zero or the treasure has no elements it returns a nil value. It then splits the treasure list into a first element and the tail and attempts to find a solution with that first element. If that fails, it tries again but now without the first element. Repeating this for every element in the list, we've exhausted the problem space and have found a solution if it exists at all. 

CODE

The main function simply creates a treasure list, calculates it total value, asks for a two way split, and prints the result.

CODE

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c1/partition.go).