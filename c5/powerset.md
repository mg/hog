5.4.1: Someone Else's Problem

A power set P of set S is a collection of sets constructed in such a way that each set Sn contains a different combination of the elements of S and P contains all such distinct sets Sn. 

    P(S={s1, s2, s3}) = [
        S1={s1, s2, s3},
        S2={s1, s2},
        S3={s1, s3},
        S4={s1},
        S5={s2, s3},
        S6={s2},
        S7={s3}
    ]

The Perl code in *HOP* is taken from the book *Mastering Algorithms with Perl* and the power set is constructed as a hash of hashes. It is full of weird perlisms and I even had to seek the help of MJD (thanks!) to understand it. Still I can't say I fully understand it and why it is implemented the way it is implemented. If anyone is interested, the original code is [here](https://gist.github.com/mg/6747532).

    type set map[string]string
    type powerset []set

For me the most natural way to construct this is to represent the powerset as a slice of maps. Perls magic hash foo is hard to grok and I'm not sure what the benefit is to do it like that. Feel free to tell me how my version is way worse Big Oh wise!

    func keysAndValues(s set) ([]string, []string) {
        var ks []string
        var vs []string
        for k := range s {
            ks = append(ks, k)
            vs = append(vs, s[k])
        }
        return ks, vs
    }

Perl has *keys* and *values* to retrieve the keys and values from hashes. We need something similar to retrieve them from maps. Doing it in a single function is beneficial since then we can do it in a single pass.

The signature of the *powerset* function is *powerset_recurse(s set, p powerset, keys, values []string, n, i int) powerset*. *s* is our set S that we wish to use as a seed for the powerset. All of the other parameters are the state that needs to be passed from each invocation to the next. This is something that the function maintains and the caller does not need to know about. In Perl we simply call the function without these parameters:

    powerset_recurse($set)

Go has no default parameters which forces us to provide *nil* values for all the state parameters: 

    powerset_recurse(s, nil, nil, nil, 0, 0)

This is clumsy and ugly. To avoid this, we can wrap the recursive function in another function that provides the nicer interface to the caller:

    func powerset_recurse(s set) powerset {
        var f func(s set, p powerset, keys, values []string, n, i int) powerset
        f = func(s set, p powerset, keys, values []string, n, i int) powerset {
    
            // function body
    
            return f(s, p, keys, values, n, i+1)
        }
        return f(s, nil, nil, nil, 0, 0)
    }

Now we can call the function and start the recursive computation in the same way as in Perl.

        if p == nil {
            keys, values = keysAndValues(s)
            n = len(keys)
            p = make(powerset, int(math.Pow(2, float64(n)))-1)
            i = 1
        }
        if i-1 == n {
            return p
        }

The first part of the body simply sets up the state variables for the recursion if the power set that we wish to create is *nil*. We retrieve a slice of keys and values from the source set and make the powerset to contain (2 to the power of n) - 1 values, n being the number of keys in the source set. *i* is our counter, it represent the number of recursions we've executed. Once *i-1* is equal to *n* the powerset is ready and we stop. Naturally, this initialization code could simply be outside the recursive function and inside the wrapping function that we created to improve the interface. But it is here in the Perl version so I'll do the same.

        c := int(math.Pow(2, float64(i-1)))
        for j := 0; j < c; j++ {
            ss := make(set, i)
            for k := 0; k < i; k++ {
                flag := 1 << uint(k)
                if (c+j)&flag == flag {
                    ss[keys[k]] = values[k]
                }
            }
            p[c-1+j] = ss
        }

This is where we generate each Sn in the powerset and where the Perl code contains much of its black magic. My strategy is to use bitmasks to generate the sets.

*c* is the number of sets to generate in this invocation, 1 for the first, 2 for the second, and then 4, 8, 16, etc., as needed. The invocation counter *i* controls the size of each set, 1 for the first time, 2 then, and 3, 4, 5, etc.

In the inner loop, we copy keys and values from the source set to set Sn. We generate a flag for each value such that value n is the flag ...0001..., that is, contains the value 1 at the nth position and 0 elsewhere. *c+j* is the mask for set Sn. A mask such as *1101* means the set contains keys and values 1, 2 and 4 while a mask such as *001011* means that the set contains keys and values 3, 5 and 6.

After this loop is done we assign the set Sn to the correct location in the powerset.

The complete function is as follows:

    func powerset_recurse(s set) powerset {
        var f func(s set, p powerset, keys, values []string, n, i int) powerset
        f = func(s set, p powerset, keys, values []string, n, i int) powerset {
            if p == nil {
                keys, values = keysAndValues(s)
                n = len(keys)
                p = make(powerset, int(math.Pow(2, float64(n)))-1)
                i = 1
            }
            if i-1 == n {
                return p
            }
    
            c := int(math.Pow(2, float64(i-1)))
            for j := 0; j < c; j++ {
                ss := make(set, i)
                for k := 0; k < i; k++ {
                    flag := 1 << uint(k)
                    if (c+j)&flag == flag {
                        ss[keys[k]] = values[k]
                    }
                }
                p[c-1+j] = ss
            }
            return f(s, p, keys, values, n, i+1)
        }
        return f(s, nil, nil, nil, 0, 0)
    }

To remove the recursion we simply have to add a for loop around the latter part of the body. And now we can remove the inner function body as well, making the whole thing a lot simpler:

    func powerset_loop(s set) powerset {
        keys, values := keysAndValues(s)
        n := len(keys)
        p := make(powerset, int(math.Pow(2, float64(n)))-1)
    
        for i := 1; i <= n; i++ {
            c := int(math.Pow(2, float64(i-1)))
            for j := 0; j < c; j++ {
                ss := make(set, i)
                for k := 0; k < i; k++ {
                    flag := 1 << uint(k)
                    if (c+j)&flag == flag {
                        ss[keys[k]] = values[k]
                    }
                }
                p[c-1+j] = ss
            }
        }
        return p
    }

The benchmarks point to something around a 40% improvement in execution speed over the recursive version:

    Powerset recursive at 14074
    Powerset looping at 8858
    Powerset recursive at 15000
    Powerset looping at 10000
    Powerset recursive at 14138
    Powerset looping at 11909

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c5/powerset.go).