3.7.4: Partitioning

The point of making a generic memoize is to be able to use it with different algorithms. The *partition* algorithm from [1.8: Partitioning](http://higherordergo.blogspot.com/2013/07/18-partitioning.html) has a different signature than the fibonacci algorithm yet we can use [*memoizer*](http://higherordergo.blogspot.com/2013/08/35-memoize-module_4.html) unchanged. First for the type definitions, the only difference being the signature of the *Algofunc*:

CODE

*memoize* is exactly the same as before:

CODE

The only change we need to make to *findshare* is to change it to a constructor that accepts a *Memfunc* and returns the actual algorithm. It is the same patterns as we used for *fib*:

CODE

In the *main* function we start with doing unmemoized calculations:

CODE

Then we proceed to do the memoized calculations on the same data, constructing the same typecasting code structure as before:

CODE

The speedup we gain for the memoization is about 96% for this example:

	Raw partition [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20] into 209: [2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20] at 558.018056ms
	Memoized partition [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20] into 209: [2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20] at 18.783834ms

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c3/partition.md).