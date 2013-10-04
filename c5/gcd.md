5.4.1: Tail-Call Elimination

One of the reasons recursive solutions often look so nice is the automatic stack management the computer performs for us. Every time we call a function the parent functions parameters get saved on the stack and then popped back when the function returns and the parent regains control. Furthermore, any return value of the function gets pushed on the stack and the tail end of the function and then popped back in the parent function. As nice as such automatic stack management is, it is not free. Pushing and popping takes time.

A tail call occurs when a function ends by calling another function (such as itself) and then immediately returning the value of the last function call without any additional computation. A recursive function R that calls itself three times R1, R2, and R3 at the end will end up pushing and popping the same return value three times as the call stack unwinds.

Tail-Call elimination is the process of eliminating this unnecessary pushing-and-popping, simply returning the value of R3 straight to the caller of R1 without going through R1 and R2.

Many compilers and runtimes provide for this optimization automaticly. As I write this, Go does not. But we can usually refactor a recursive tail-call function into a simple for loop without any recursion, therefore removing the tail call.

    func gcd(m, n int) int {
        if n == 0 {
            return m
        }
        m = int(math.Mod(float64(m), float64(n)))
       return gcd(n, m)
    }

The *gcd* computes the *Greatest Common Divisor* between *m* and *n*, using Euclid's algorithm. At each call we make *m* be equal to *n* and *n* be equal to the modulus of *m* and *n*. Once *n* hits 0 we are done and return *m* as the final result. Notice that no additional computation is performed on this value of *m*, it is simply passed back up the entire call stack as the result of the original computation.

    func gcd_notail(m, n int) int {
        for n != 0 {
            m, n = n, int(math.Mod(float64(m), float64(n)))
        }
        return m
    }

To eliminate the tail call we rewrite the function as a *for* loop. *n* is the condition of the loop as it was the condition of the call stack, once it hits zero we are done. And at each iteration *m* becomes *n* and *n* becomes the modulus of *m* and *n*.

    GCD Recursive of 48 and 20: 4 at 3.298us
    GCD NoTail of 48 and 20: 4 at 431ns
    GCD Recursive of 48 and 20: 4 at 3.55us
    GCD NoTail of 48 and 20: 4 at 624ns
    GCD Recursive of 48 and 20: 4 at 4.669us
    GCD NoTail of 48 and 20: 4 at 1.05us

As we can see, the non recursive version of GCD is always faster than the recursive one, by a factor of 4 to 8 in these three samples.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c5/gcd.go).