5.4.3: Explicit Stacks

So far the strategy has been to identify what the state is that the call stack is managing for us and build a loop around that state to eliminate the function calls. But sometimes the state is complex enough that we need to bring in the heavy machinery.

    func fib1(n int) int {
        if n < 2 {
            return n
        }
        return fib1(n-2) + fib1(n-1)
    }

The Fibonacci function is a tricky thing when it comes to eliminating the tail call. First we have to account for the final computation in the function, the adding of the two values. But we also have to account for the fact that the function calls it self recursively twice every time, with two different values. So this is not a linear recursive calculation but an execution tree. There is some branching logic we have to account for if we wish to eliminate the recursion.

A sample call of fib(6) looks like this:

    fib(6)
        fib(4)
            fib(2)
                fib(0)
                fib(1)
            fib(3)
                fib(1)
                fib(2)
                    fib(1)
                    fib(0)
        fib(5)
            fib(3)
                fib(1)
                fib(2)
                    fib(0)
                    fib(1)
            fib(1)

Flattening this to something that can be run in a *for* loop will require us to do branch management along with stack management.

    func fib2(n int) int {
        for {
            if n < 2 {
                return n
            }
            s1 := fib2(n - 2)
            s2 := fib1(n - 1)
            return s1 + s2
        }
    }

First we make some cosmic changes to frame the problem a bit. The *for* loop is what will replace the recursion. And breaking up the calculation seems to point to a state consisting of three variables, *n*, *s1* and *s2*. We will use a stack to track these states along with the branch control.

    type state3 struct {
        BRANCH, s1, s2, n int
    }
    
    func fib3(n int) int {
        var s1, s2, retval int
        BRANCH := 0
        var STACK list.List
        for {
            if n < 2 {
                retval = n
            } else {
                if BRANCH == 0 {
                    STACK.PushBack(&state3{BRANCH, s1, s2, n})
                    n -= 2
                    BRANCH = 0
                    continue
                } else if BRANCH == 1 {
                    s1 = retval
                    STACK.PushBack(&state3{BRANCH, s1, s2, n})
                    n -= 1
                    BRANCH = 0
                    continue
                } else if BRANCH == 2 {
                    s2 = retval
                    retval = s1 + s2
                }
            }
            if STACK.Len() == 0 {
                return retval
            }
            s := STACK.Back().Value.(*state3)
            STACK.Remove(STACK.Back())
            BRANCH, s1, s2, n = s.BRANCH, s.s1, s.s2, s.n
            BRANCH++
        }
    }

The key part here is the value of *BRANCH*. 0 corresponds to *fib(n-2)*, 1 to f*ib(n-1)* and 2 to *s1+s2*. If we compare to the *fib(6)* tree illustrated above, the first three iterations *BRANCH* will be 0 and *n* will go *6, 4, 2* as the loop continues from the *BRANCH == 0* block. 

On the next iteration n will be 0 and we drop to the block below the branching logic. There the stack will be popped bringing back the values for *n == 2* (restoring *n* to 2) and the *BRANCH* will be increased to 1. 

So on the next iteration we hit the *BRANCH == 1* block, *n* becomes 1. As we continue we hit n < 2 again and pop the state from the BRANCH == 1 block bringing BRANCH to 2. As we loop back we hit the BRANCH == 2 bock and perform the addition. The addition is stored in *retval*. Now we pop n == 4 from the stack. This again brings us to the BRANCH == 1 block and brings n down to 3.

This sequence corresponds to *fib(6) -> fib(4) -> fib(2) -> fib(0) -> fib(1) -> fib(3)*. The tree branching logic has been flattened in such a way that the left child is always travelled first and then we go down the right branch once we pop back. This is called *pre-order tree traversal*. The rest of the tree will we traversed in the order of *fib(1) -> fib(2) -> fib(1) -> fib(0) -> fib(5) -> fib(3) -> fib(1) -> fib(2) -> fib(0) -> fib(1) -> fib(1)*, completing the same calculation as the recursive version.

    type state4 struct {
        BRANCH, s1, n int
    }
    
    func fib4(n int) int {
        var s1, retval int
        BRANCH := 0
        var STACK list.List
        for {
            if n < 2 {
                retval = n
            } else {
                if BRANCH == 0 {
                    for n >= 2 {
                        STACK.PushBack(&state4{BRANCH, 0, n})
                        n -= 2
                    }
                    retval = n
                } else if BRANCH == 1 {
                    s1 = retval
                    STACK.PushBack(&state4{BRANCH, retval, n})
                    n -= 1
                    BRANCH = 0
                    continue
                } else if BRANCH == 2 {
                    retval += s1
                }
            }
            if STACK.Len() == 0 {
                return retval
            }
            s := STACK.Back().Value.(*state4)
            STACK.Remove(STACK.Back())
            BRANCH, s1, n = s.BRANCH, s.s1, s.n
            BRANCH++
        }
    }

The hard part is over and all that is left is to do some code clean up and optimizations. The BRANCH == 0 can be optimized to include its own for loop to calculate the entire branch from n to 0 without involving the outer loop. So once we hit it with e.g. n == 6, n == 4 and n == 2 get pushed on the stack straight away before we continue with the main loop.

Furthermore, *s2* is not really needed. Its value is always either 0 (when *BRANCH == 0*) or the same as *retval* (when *BRANCH == 1*).

    func fib5(n int) int {
        var s1, retval int
        BRANCH := 0
        var STACK list.List
        for {
            if n < 2 {
                retval = n
            } else {
                if BRANCH == 0 {
                    for n >= 2 {
                        STACK.PushBack(&state4{1, 0, n})
                        n -= 1
                    }
                    retval = n
                } else if BRANCH == 1 {
                    s1 = retval
                    STACK.PushBack(&state4{2, retval, n})
                    n -= 2
                    BRANCH = 0
                    continue
                } else if BRANCH == 2 {
                    retval += s1
                }
            }
            if STACK.Len() == 0 {
                return retval
            }
            s := STACK.Back().Value.(*state4)
            STACK.Remove(STACK.Back())
            BRANCH, s1, n = s.BRANCH, s.s1, s.n
        }
    }

As a final optimization we swap the logic for *BRANCH == 0* and *BRANCH == 1*. It doesn't really matter whether we descend down the execution tree left child first or right child first, the final result is the same. After all, *a+b* is equal to *b+a*. But the inner loop for the left path is executed *n/2* times while the inner loop for the right path is executed *n* times. If the optimization is beneficial for *n/2*, it should be even better for *n*.

Furthermore, rather than pushing a branch value on the stack and then increment it by one when we get it back, we can simply push the incremented value on the stack and get rid of the trailing *BRANCH++* statement.

After all this hard work, the results are a bit disappointing. The optimized version is around 100 times slower than the recursive version.

    Fib1: fib(20) = 6765 at 83.407us
    Fib5: fib(20) = 6765 at 6.380718ms
    Fib1: fib(20) = 6765 at 83.304us
    Fib6: fib(20) = 6765 at 5.68734ms
    Fib1: fib(20) = 6765 at 83.497us
    Fib6: fib(20) = 6765 at 5.682758ms

So apparently Go is doing a good job here, at least better than we can accomplish.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c5/fib.go).