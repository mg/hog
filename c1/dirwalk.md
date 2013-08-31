1.5: Applications and Variations of Directory Walking

The *dirwalk* is an generic function (yes, Go lacks generics, we will get to that) that walks a directory tree and executes a supplied file function or a supplied directory function for each entry. 

To start there are some types that need to be declared. *ComputeType* is our generic type and *FileFunc* and *DirFunc* are the two user supplied functions.

CODE

Next comes a utility function that is only needed so that the user can call dirwalk without either a *FileFunc* or a *DirFunc* if those are not needed. It simply returns a empty value.

CODE

*dirwalk* opens a file entry and checks if it is a regular file or a directory. If it is a regular file it calls the *FileFunc* on it and returns the result. If it is a directory it loops through all entries of that directory calling itself recursively on each entry, collecting the results in a slice. It then calls the *DirFunc* on the result collection and returns the result of the *DirFunc*.

CODE

Now the lack of proper generics starts to complicate things. So far *dirwalk* has only used the *interface{}* type but the point of all this is to compute actual values. Typecasting is needed. To reduce the code the following utility function is used. It tries to typecast and panics on failure. Another way to write this function is to accept a default value and return that on failure.

CODE

The *FileFunc* used here simply returns the size of the file, which is an int64 value.

CODE

*DirFunc* goes through the results slice, extracting the *int64* value for each entry and returning the total.

CODE

The *main* function simply calls *dirwalk* with the size calculation functions and then extracts the *int64* value.

CODE

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c1/dirwalk.go).