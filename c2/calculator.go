package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type (
	ActionFunc  func(string, *Stack)
	ActionTable map[string]ActionFunc
)

// evaluates expression using the supplied dispatch table
func evaluate(expression []string, actions ActionTable, stack *Stack) interface{} {
	for _, t := range expression {
		var action ActionFunc
		if _, err := strconv.ParseFloat(t, 64); err == nil {
			action = actions["NUMBER"]
		} else {
			var ok bool
			if action, ok = actions[t]; !ok {
				action = actions["__DEFAULT__"]
			}
		}
		action(t, stack)
	}
	return stack.Pop()
}

func main() {

	// dispatch table for calculator
	calcActions := ActionTable{
		"+": func(token string, stack *Stack) {
			stack.Push(stack.PopFloat() + stack.PopFloat())
		},
		"-": func(token string, stack *Stack) {
			v := stack.PopFloat()
			stack.Push(stack.PopFloat() - v)
		},
		"*": func(token string, stack *Stack) {
			stack.Push(stack.PopFloat() * stack.PopFloat())
		},
		"/": func(token string, stack *Stack) {
			v := stack.PopFloat()
			stack.Push(stack.PopFloat() / v)
		},
		"sqrt": func(token string, stack *Stack) {
			stack.Push(math.Sqrt(stack.PopFloat()))
		},
		"NUMBER": func(token string, stack *Stack) {
			v, _ := strconv.ParseFloat(token, 64)
			stack.Push(v)
		},
		"__DEFAULT__": func(token string, stack *Stack) {
			panic(fmt.Sprintf("Unkown token %q", token))
		},
	}

	v := evaluate(os.Args[1:], calcActions, new(Stack))
	fmt.Printf("Result: %f\n", toFloat(v))

	// dispatch table for AST tree generator
	astActions := ActionTable{
		"NUMBER": func(token string, stack *Stack) {
			stack.Push(token)
		},
		"__DEFAULT__": func(token string, stack *Stack) {
			t := stack.Pop()
			if stack.Len() > 0 {
				stack.Push([]interface{}{token, stack.Pop(), t})
			} else {
				stack.Push([]interface{}{token, t})
			}

		},
	}

	v = evaluate(os.Args[1:], astActions, new(Stack))
	fmt.Printf("AST tree: %v\n", v)
	fmt.Printf("AST to string: %q\n", astToString(toInterfaces(v)))

}

// build an infix string from AST tree
func astToString(tree []interface{}) string {
	if len(tree) == 1 {
		return toString(tree[0])
	}
	if len(tree) == 2 {
		op, val := toString(tree[0]), toInterfaces(tree[1])
		s := astToString(val)
		return "( " + op + " " + s + " )"
	}
	op, l, r := toString(tree[0]), toInterfaces(tree[1]), toInterfaces(tree[2])
	s1 := astToString(l)
	s2 := astToString(r)
	return "( " + s1 + " " + op + " " + s2 + " )"
}

// Stack structure
// Based on https://gist.github.com/bemasher/1777766
type Stack struct {
	head *item
	size int
}

type item struct {
	value interface{}
	next  *item
}

func (s *Stack) Len() int {
	return s.size
}

func (s *Stack) Push(value interface{}) {
	s.head = &item{value, s.head}
	s.size++
}

func (s *Stack) Pop() (value interface{}) {
	if s.size > 0 {
		value, s.head = s.head.value, s.head.next
		s.size--
		return
	}
	return nil
}

func (s *Stack) PopFloat() float64 {
	return toFloat(s.Pop())
}

// General typecasting functions
func toFloat(value interface{}) float64 {
	if v, ok := value.(float64); ok {
		return v
	}
	return 0.0
}

func toString(value interface{}) string {
	if v, ok := value.(string); ok {
		return v
	}
	return ""
}

func toInterfaces(value interface{}) []interface{} {
	if v, ok := value.([]interface{}); ok {
		return v
	}
	return []interface{}{value}
}
