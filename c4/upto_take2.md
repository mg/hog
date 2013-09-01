4.2.1: A Trivial Iterator - Take two

The standard library in Go has some examples of iterators. One is the [bufio.Scanner](http://golang.org/pkg/bufio/#Scanner). The scanner has two primary methods for iteration, the *Scan()* method that moves the iterator forward and checks for *EOS*, and the *Text()* method that returns the value at the current location.

    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        text := scanner.Text()
    }

The Upto iterator can easily be written using this pattern. The previous function-of-three-hats is now broken into two.

    type upto struct {
        m, n int
    }
    
    func Upto(m, n int) *upto {
        return &upto{m: m - 1, n: n}
    }
    
    func (i *upto) Next() bool {
        i.m++
        return i.m < i.n
    }
    
    func (i *upto) Value() int {
        return i.m
    }

Rather than using a closure to manage state we now use a struct. We then define two methods that act on that structure, *Next()* and *Value()*. *Value()*, like *Text()* in *Scanner* is an idempotent function that returns the current value in the stream. *Next()*, like *Scan()*, both advances the iterator and checks for *EOS*.

    func main() {
        i := Upto(3, 5)
        for i.Next() {
            fmt.Println(i.Value())
        }
    }

Our iteration is simple enough, after constructing the function we can use the *while* form of the *for* loop to run through the data.

Even though this form requires a bit more boilerplate than the previous one, I believe it is a more useful way. Now we have something resembling a contract that other functions and algorithms can build on. But there are problems, *Next()* is still doing to much. An idempotent way to check for *EOS* is often neccessery. 

Another problem is that *Next()* is secretly serving a third purpose. Since it gets called once before we retrieve the first value to guard against an empty stream, it serves as an initialization function. This means the iterator must be in a weird pre-initialization state after the construction. It's easy enough in this case but could prove problematic when solving different problems.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/upto_take2.go).