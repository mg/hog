5.4.2: Creating Tail Calls

If a recursive function doesn't have a tail call to eliminate, we can often create it. And then promptly eliminate it. A truly creative destruction.

The binary function accepts an integer and returns a string that is the binary representation of the integer. E.g. *binary(10)* returns the string *"1010"*.

    func binary1(n int) string {
        if n <= 1 {
            return strconv.Itoa(n)
        }
        k := int(n / 2)
        b := int(math.Mod(float64(n), 2))
        return binary1(k) + strconv.Itoa(b)
    }

Even though the calls is at the tail of the function, this is not a tail call since the last computation performed in the function is the adding of the two strings; the return of the recursive call and the conversion of *b*.

    func binary2(n int, retval string) string {
        if n < 1 {
            return retval
        }
        k := int(n / 2)
        b := int(math.Mod(float64(n), 2))
        retval = strconv.Itoa(b) + retval
        return binary2(k, retval)
    }

To eliminate the final computation we introduce a new state variable, *retval*, which will contain the state of the computation at all times and be passed along the recursive calls. Now the final computation is the recursive call which is now a tail call that we can eliminate. The cost is that the computation takes one extra recursive call. But that is acceptable since that one extra call will be turned into one extra iteration.

    func binary3(n int) string {
        var retval string
        for n >= 1 {
            b := int(math.Mod(float64(n), 2))
            retval = strconv.Itoa(b) + retval
            n = int(n / 2)
        }
        return retval
    }

The recursive function now becomes a simple function with a *for* loop, just as the [gcd]() example. *n>=1* guards our loop that allows us to build up the string in *retval* that we return once the loop exists.

    Binary1: binary(10) ="1010" at 13.448us
    Binary3: binary(10) ="1010" at 2.461us
    Binary1: binary(10) ="1010" at 8.896us
    Binary3: binary(10) ="1010" at 1.581us
    Binary1: binary(10) ="1010" at 12.516us
    Binary3: binary(10) ="1010" at 2.02us

Again, the loop version of the computation destroys the recursive version. And this is even though we are using a somewhat inefficient (and not recommended by the Go authors) method of string concatenation.

The strategy can be summarized as to identify the tail computation that prevents the tail-call elimination, store the computation in a variable that can be passed with the recursion, isolate the tail-call and finally eliminate it.

The source is at [GitHub](https://github.com/mg/hog/blob/master/c5/binary.go). Also, another example rewriting a factorial function in the same way is [available](https://github.com/mg/hog/blob/master/c5/factorial.go). The strategy is exactly the same so I won't go further into it.
