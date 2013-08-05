Chapter 2 is all about breaking long if-else chains into maps of values and functions. This will allow us to separate the branching logic from the algorithm, enabling us to easily modify and extend the branching logic and even completely replace it with a different one.

In this first example we create a process that will read through a configuration file that has the form:

DIRECTIVE PARAMETERS

The algorithm runs through the file and executes the function for the *DIRECTIVE* passing to it the parameters defined in the file.

To start we define our types, the table that holds the branching logic and the functions that will represent the individual branches. In this case the function accepts a slice of strings and a reference to the branching logic.

CODE

The *onReadConfig* function expects a filename for its argument. It opens the file, reads each line of text, breaks it into tokens, looks up the function to execute using the dispatch table and the first token and then executes that function passing the rest of the line as parameter. It is the core algorithm of this program, but interestingly, it is fully reentrant and has the same signature as other functions in the dispatch table. It can therefore execute itself through the dispatch table in a round-about recursive way.

CODE

*onDefine* is a meta function that defines a directive in terms of another existing directive. What it allows us to do is to define a name that executes an directive with default parameters. For example, if directive *CD* is defined and it expects a directory as a parameter, we can write *DEFINE HOME CD /home/* in the configuration file, therefore creating a new directive *HOME* that simply executes directive *CD* with the parameters */home/*. In HOP, MJD uses *DEFINE* to define a directive and actual Perl code in the configuration file that is then dynamically evaluated. With our statically compiled code we will have to do with a less powerful version.

CODE

The main function creates the dispatch table with four directives. *CONFIG* and *DEFINE* point to our previous functions while *PRINT* and *CD* are two very self-explanatory functions. Examples of configuration files can been seen [here](https://github.com/mg/hog/blob/master/c2/1.conf) and [here](https://github.com/mg/hog/blob/master/c2/2.conf).

CODE

Get the source at [GitHub](https://github.com/mg/hog/blob/master/c2/configuration-file-handling.go).