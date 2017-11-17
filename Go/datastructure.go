package main

import (
	"fmt"
)

type storage struct {
	vartype int
	data    interface {
		toString() string
	}
}

func (s storage) toString() string {
	return fmt.Sprint(s.vartype, " ", s.data)
}

type lexical struct {
	fsm  int
	data string
}

func (l lexical) toString() string {
	return fmt.Sprintf("(%d,%s)\n", l.fsm, l.data)
}

func (l lexical) print() {
	fmt.Printf("(%d,%s)\n", l.fsm, l.data)
}

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

func (q *Stack) Top() *storage {
	return (*q)[q.Len()-1]
}

func (q *Stack) DeVarPop(variable *Variable) *storage {
	x := q.Len() - 1
	n := *(*q)[x]
	*q = (*q)[:x]
	for n.vartype == VAR_VAR {
		n = variable.get(n.data.(Var))
	}
	return &n
}

func (q *Stack) DeVarTop(variable *Variable) *storage {
	n := *(*q)[q.Len()-1]
	for n.vartype == VAR_VAR {
		n = variable.get(n.data.(Var))
	}
	return &n
}

func (q *Stack) Len() int {
	return len(*q)
}

type Lexs []lexical

func (q Lexs) Len() int {
	return len(q)
}

func (q *Lexs) Push(n lexical) {
	*q = append(*q, n)
}

func (q *Lexs) Pop() (n lexical) {
	x := q.Len() - 1
	n = (*q)[x]
	*q = (*q)[:x]
	return
}

func (q Lexs) Top() lexical {
	return q[q.Len()-1]
}

const (
	VAR_NUMBER int = iota
	VAR_STRING
	VAR_VAR
	VAR_FSM
	VAR_FUNCTION
	VAR_C_FUNCTION
	VAR_SET
	VAR_ARRAY
	VAR_CLASS
)

const (
	LEX_NUMBER int = iota
	LEX_ASSIGN
	LEX_EXPR
	LEX_OPR
	LEX_CALL
	LEX_BRACKET
	LEX_END_BRACKET
	LEX_BRACES
	LEX_VEC
	LEX_ADDR
	LEX_SEPERATOR
	LEX_STR
	LEX_NAMESPACE
	LEX_FUNC
)
