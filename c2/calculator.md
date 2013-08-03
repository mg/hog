As in [2.1](http://higherordergo.blogspot.com/2013/07/21-configuration-file-handling.html) we are aiming to separate the branching logic from the algorithm, making it pluggable and the algorithm reusable. And in this example, we actually do code two separate branching tables using the same algorithm for two different effects.

This is a simple *Reverse Polish Notation* calculator, accepting expressions in the form of *"2 3 + 1 -"* and evaluating a result depending on the branching logic supplied.

We start with a very similar type declaration as in the previous post, a action function and a dispatch table. The *Stack* is used to record the state of the evaluation.

CODE

The *evaluate* function accepts the expression, the branching logic and a stack for the state. It loops through the expression and evaluates which action to take depending on the token. If the token is a number the *NUMBER* action is selected. Otherwise an action is selected depending on the token or a default one if that fails. The action is then executed.

The function ends with popping the top value of the stack and returning it to the caller.

CODE

The main function starts with declaring a branching logic that supports calculating an expression in RPN format using the +, -, \*, / and sqrt operators. Then the *evaluate* algorithm is called and the result is printed.

CODE

Next we create a branching logic that supports building an *Abstract Syntax Tree* from a RPN expression. The *\__DEFAULT\__* function is selected for all tokens except numbers, building a tree of slices on the stack. The result is then printed twice, once in the raw format of the Go data structure and once in the form of an infix string, built by the *astToString* function.

CODE

At the end of the file, not shown here, are some auxiliary functions. A stack data structure is defined along with three typecasting functions.

Get the full code at [GitHub](https://github.com/mg/hog/blob/master/c2/calculator.go).