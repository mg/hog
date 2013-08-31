3.5: The Memoize Module

As mentioned before Go does not have at this time (version 1.1) generics. Go does not have any kind of macro language that allows us to reason about code while compiling. And as far as I know Go has no way to change the function symbol table while compiling. There is a package called debug/gosym but it seems to be more about inspection than meddling; feel free to correct me.

The memoize function in HOP uses some Perl features that are hard to mimic in Go. For one, Perl is very liberal with function signatures, as in there are none. Every function is simply a name that accepts any number and type of arguments through the @_ symbol. In Go that means the variadic function f(...interface{}) and a whole lot of typecasting.

Second, and far harder, is Perl's ability to dynamically change the symbol table. The idea behind memoize is to speed up recursive functions by caching computations. Chaching the first computation is easy because we can simply call the memoized version rather than the raw version. But that doesn't help because the raw version then calls itself recursively rather than the memoized version, meaning the expensive computation is not cached at all.

But hope is here. Closures, anonymous functions and functions as first class objects ride in to save the day. Sort of.

First we start by defining a few function types. *Keyfunc* is a function to calculate the cache key for the memoizer. Memfunc is the memoizer. And *Algofunc* is the type of the algorithm we want to speed up. In this example we are using the fibonacci computation.

CODE

The memoizer accepts two functions. First a *Memfunc* which is the algorithm wrapped in a typecasting function; and a *Keyfunc*. It returns the actual memoizer that is a closure that contains a reference to the cache and the keyfunc.

If no *Keyfunc* is supplied we construct a naive key function that simply runs through a slice of interface values and constructs a string representation of its value with Sprint. 

Then the cache is created and a memoizer is constructed and returned. It simply constructs the cache key, checks the cache for already computed result, returns it if it exists, otherwise it calls the algorithm supplied, stores the result and returns it.

CODE

Next we create our fibonacci calculator. Ideally this code would be in no way aware of a possible memoization but that is sadly not possible since we need to be able to intercept recursive calls. We need to take a more co-operative route.

This functions accepts a possible memoizer or nil and returns an algorithm. The function that we return is the calculator that calls the supplied memoizer and returns typecasted results from it.

Next we check if a memoizer was supplied and if not we construct a small function that simply calls the fibonacci algorithm and assign that to the memoizing function, thus making the algorithm recursive.

CODE

In the main function we start by creating a simple fibonacci calculator (by passing nil into the constructor) and calculate a single fibonacci number, timing the execution. 

CODE

Then we proceed to finally create the memoized version of the calculator. The constructor starts by declaring a *Algofunc* variable without actually declaring the function. Then the memoize constructor is called, passing in a anonymous function that calls the previously declared function, typecasting interface parameter to an int. Then the calculator is created, passing in the memoizer for the recursion. Finally an anonymous function gets returned that simply typecasts from Algofunc to Memfunc.

CODE

What follows is the artists impression of this monstrosity:

![Chart of Memoize](memoize_thumb.jpg "Memoize flow")

I heard you liked closures in your closures in your closures in your closures. 

Amazingly, this frankensteinian construction works. Sample run for few choice calculations results in:

	Raw fib 1: 1 at 260ns
	Memoized fib 1: 1 at 57.066us
	Raw fib 10: 89 at 20.838us
	Memoized fib 10: 89 at 118.31us
	Raw fib 12: 233 at 69.558us
	Memoized fib 12: 233 at 142.454us
	Raw fib 13: 377 at 117.464us
	Memoized fib 13: 377 at 137.451us
	Raw fib 14: 610 at 180.078us
	Memoized fib 14: 610 at 145.62us
	Raw fib 20: 10946 at 3.305775ms
	Memoized fib 20: 10946 at 145.757us
	Raw fib 35: 14930352 at 5.566529702s
	Memoized fib 35: 14930352 at 252.713us
	Raw fib 40: 165580141 at 53.862903338s
	Memoized fib 40: 165580141 at 364.974us

The raw calculation is faster at lower number, with the memoized version becoming quicker around fib(14). At fib(40) the raw function takes around 54 seconds while the memoized version is still returning within the second (0.36s).

Possibly the most specialized "general" code I've ever written.

Get the source at [GitHub]().
