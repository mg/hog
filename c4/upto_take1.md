4.2.1: A Trivial Iterator - Take one

Iterators are objects that allow us to iterate over a stream of data, be that from a file, a linked list, an array or any other data collection. The iterator encapsulates both location in the data and movement across it.

At very minimum, the iterator must support three operations: 

1. To move across the data. 
2. To return the data that the iterator currently points to.
3. Allow the user to check if he has reached the end of the data.
 
The Upto iterator represents iteration over the interval *[m,n)* which is the list of numbers from *m* to *n-1*. One of the benefits of using iterators in this example is that we don't have to generate the list of numbers before hand, we can simply evaluated it as needed.

A very simple way to create an iterator, and a very similar to the one that is used in HOP, is to use closures to capture the state and return a function that acts as an iterator. This function that is returned must be able to satisfy all three criterias in a single function call. 

CODE

The function accepts two numbers and returns a function that acts as an iterator. Upon invocation it increases the current position and returns the previous number and a boolean that indicates whether we have reached the end of the stream.

CODE

Looping is simple, but a bit clumsy. We start by constructing the iterator. Then in the initialization section of the *for* construct we retrieve the first value and the *EOS* (End Of Stream) indicator. The last section of the *for* constructs gets the next value along with updating the *EOS* indicator.

It's all very trivial and brief, there is very little code here that doesn't directly touch on either generating the interval or looping through it.

But I don't think it is a good fit for iterators in Go. For one, the invocation performs way to many functions in one call. There is e.g. no way to check for *EOS* without advancing it. Second, it is very limiting in the kind of iterators we can create. Random access is off the table. Thirdly, it is very hard to write algorithms for these kind of iterators. There is no way to declare what kind of iterator the algorithm needs. Fourth, I don't think this style of iterators scales very well across different problems.

Go is a static language; formalism is our friend and path to freedom.

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c4/upto_take1.go).