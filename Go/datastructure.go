package main

import (
	"container/list"
	"fmt"
	"math"
	"strconv"
	"strings"
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

type number struct {
	numerator   int
	denominator int
}

var True, False storage

func (n *number) construct(v string) {
	raw := strings.Split(v, ".")
	switch len(raw) {
	case 1:
		temp, err := strconv.Atoi(raw[0])
		if err != nil {
			panic(err)
		}
		n.numerator = temp
		n.denominator = 1
	case 2:
		temp, err := strconv.Atoi(raw[0])
		if err != nil {
			panic(err)
		}
		n.numerator = temp
		temp, err = strconv.Atoi(raw[1])
		if err != nil {
			panic(err)
		}
		n.denominator = int(math.Pow10(int(math.Ceil(math.Log10(float64(temp + 1))))))
		n.numerator = n.numerator*n.denominator + temp
	default:
		panic("Invalid number: " + v + "\n")
	}
	n.reduce()
}

func (n *number) toString() string {
	if n.denominator == 1 {
		return fmt.Sprint(n.numerator)
	} else {
		return fmt.Sprintf("%d/%d", n.numerator, n.denominator)
	}
}

func (n *number) toInt() int {
	return n.numerator / n.denominator
}

func (n *number) print() {
	fmt.Print(n.toString())
}

type set struct {
	data *list.List
}

func (s *set) find(n number) bool {
	for e := s.data.Front(); e != nil; e = e.Next() {
		if e.Value == n {
			return true
		}
	}
	return false
}

func (s *set) new() {
	s.data = list.New()
}

func (s *set) remove(n number) bool {
	for e := s.data.Front(); e != nil; e = e.Next() {
		if e.Value == n {
			s.data.Remove(e)
			return true
		}
	}
	return false
}

func (s *set) append(n number) {
	flag := true
	for p := s.data.Front(); p != nil; p = p.Next() {
		if !lt(p.Value.(number), n) {
			flag = false
			if !equal(p.Value.(number), n) {
				s.data.InsertBefore(n, p)
				return
			} else {
				return
			}
		}
	}

	if flag {
		s.data.PushBack(n)
	}
}

func (s *set) toString() (ret string) {
	ret = "{"
	var temp number
	for p := s.data.Front(); p != nil; p = p.Next() {
		temp = p.Value.(number)
		if p.Next() != nil {
			ret += temp.toString() + ", "
		} else {
			ret += temp.toString() + "}"
		}
	}
	return
}

func (s *set) print() {
	fmt.Print(s.toString())
}

type function struct {
	argv  []storage
	exprs []storage
}

func (f *function) function(vals []storage, variable *Variable) (ret_val *storage) {
	for i, val := range f.argv {
		variable.add(val.data.(string), vals[i])
	}
	ret_val = evaluate(f.exprs, variable)
	for _, val := range f.argv {
		variable.del(val.data.(string))
	}
	return
}

type inherit struct {
	function func(data []storage, variable *Variable) *storage
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

func (q *Stack) Top() (n *storage) {
	x := q.Len() - 1
	n = (*q)[x]
	return
}

func (q *Stack) Len() int {
	return len(*q)
}

func init() {
	True = storage{var_fsm["number"], number{1, 1}}
	False = storage{var_fsm["number"], number{0, 0}}
	var_fsm = make(map[string]int)
	var_fsm["number"] = 0
	var_fsm["string"] = 1
	var_fsm["expr"] = 2
	var_fsm["fsm"] = 3
	var_fsm["inherit"] = 4
	var_fsm["set"] = 5
	var_fsm["array"] = 6
	var_fsm["namespace"] = 7

	lex_fsm = make(map[string]int)
	lex_fsm["number"] = 0
	lex_fsm["assign"] = 1
	lex_fsm["expr"] = 2
	lex_fsm["opr"] = 3
	lex_fsm["call"] = 4
	lex_fsm["bracket"] = 5
	lex_fsm["end_bracket"] = 6
	lex_fsm["set"] = 7
	lex_fsm["vec"] = 8
	lex_fsm["addr"] = 9
	lex_fsm["end_brace"] = 10
	lex_fsm["seperator"] = 11
	lex_fsm["str"] = 12
	lex_fsm["namespace"] = 13
	lex_fsm["func"] = 14
}
