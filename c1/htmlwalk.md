1.7: Html

This section took more time than I expected simply because I was hit by the double whammy of not knowing Perl well enough and not knowing Go well enough at the same time.

What I failed to appreciate in the previous section ([1.5: Applications and Variations of Directory Walking](http://higherordergo.blogspot.com/2013/07/15-applications-and-variations-of.html)) was that Perls Push method would, if it received an array to push onto a array, flatten out the second array so I would end up with a single array rather than an array of arrays. This would blow up the second use (promoting elements, see below) since I would end up with a tree of slices rather than a single slice of strings. 

My second problem was that once I figured out the Push problem, I failed to figure out the proper syntax for appending slices to slices in Go. In the end [Effective Go](http://golang.org/doc/effective_go.html) lead me to the correct solution, to append "..." to the second slice argument.

We start by declaring some types and functions for the generic walking function. These are analogues to the previous *dirwalk* function.

CODE

*htmlwalk* is our html walker. It will accept a *html.Node* and call the *TextFunc* on text elements and *ElementFunc* on other elements. It will recursively call itself and check if the result is a single *ResultType* element or a slice of *ResultType* elements and append them to the final result list.

CODE

We apply this function to two use cases. The first is to simply walk through the html file and return all the text stripped of tags.

Simple and very permissive typecasting function. If it fails it simply returns an empty string.

CODE

This is called on text elments and it simply returns the text from the html node. 

CODE

The element function. Run through the result list and concatenate all the strings.  

CODE

The second use case for the html walker is to construct code that will allow us to dig into an html document and print out only the text contained in the tag we specify. In this example we ask for the h1 tag.

First we define our tag structure (which *htmlwalker* will pass around as *ResultType*). Strings that we want to be promoted will be tagged with *Keep*, other strings will be tagged with *Maybe*.

CODE

Extract the value from *ResultType*. It can be both either a *\*tag* structure or a slice of *ResultType*. 

CODE

Any text node is a *Maybe*.

CODE

This functions constructs the *ElementFunc* function. If the tagname is found we concatenate all the strings in the result list and promote the result to *Keep*. If not, we simply return the result list that *htmlwalk* will the flatten into the final result list.

CODE

The *htmlwalk* function will return a list of *ResultType* elements. We run through it, find all elements that are tagged with *Keep* and concatenate them into one string that we then return. 

CODE

The final main function simply tests both of these use cases. You can call it either on a local html file or supply an url using the *-u* flag. 

CODE

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c1/htmlwalk.go).