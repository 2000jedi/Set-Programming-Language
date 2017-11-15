package main

import (
	"fmt"
)

type storage struct {
	vartype int
	data    interface{}
}

type lexical struct {
	fsm  int
	data string
}

func (l lexical) print() {
	fmt.Printf("(%d,%s)\n", l.fsm, l.data)
}

var var_fsm map[string]int
var lex_fsm map[string]int

type Stack []*storage

func (q *Stack) Push(n *storage) {
	*q = append(*q, n)
}

func (q *Stack) Pop() (n *storage) {
	x := q.Len() - 1
	n = (*q)[x]
	*q = (*q)[:x]
	return
}

func (q *Stack) Top() (n *storage) {
	x := q.Len() - 1
	n = (*q)[x]
	return
}

func (q *Stack) Len() int {
	return len(*q)
}

func init() {
	var_fsm = make(map[string]int)
	var_fsm["number"] = 0
	var_fsm["string"] = 1
	var_fsm["expr"] = 2
	var_fsm["fsm"] = 3
	var_fsm["c_function"] = 4
	var_fsm["set"] = 5
	var_fsm["array"] = 6
	var_fsm["namespace"] = 7
	var_fsm["function"] = 8

	lex_fsm = make(map[string]int)
	lex_fsm["number"] = 0
	lex_fsm["assign"] = 1
	lex_fsm["expr"] = 2
	lex_fsm["opr"] = 3
	lex_fsm["call"] = 4
	lex_fsm["bracket"] = 5
	lex_fsm["end_bracket"] = 6
	lex_fsm["braces"] = 7
	lex_fsm["vec"] = 8
	lex_fsm["addr"] = 9
	lex_fsm["seperator"] = 10
	lex_fsm["str"] = 11
	lex_fsm["namespace"] = 12
	lex_fsm["func"] = 13
}
