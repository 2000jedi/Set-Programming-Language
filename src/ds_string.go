// ds_string
package main

import (
	"container/list"
	"fmt"
)

type ds_string array

func (s ds_string) toString() (ret string) {
	ret = ""
	var temp storage
	var n number
	for p := s.data.Front(); p != nil; p = p.Next() {
		temp = p.Value.(storage)
		n = temp.data.(number)
		ret += fmt.Sprintf("%c", byte(n.toInt()))
	}
	return
}

func (s ds_string) get(n int) (ret *storage) {
	p := s.data.Front()
	for ; n > 0; p = p.Next() {
		if p == nil {
			panic("Index out of range")
		}
		n--
	}
	if p == nil {
		panic("Index out of range")
	}
	*ret = p.Value.(storage)
	return ret
}

func (s *ds_string) new() {
	s.data = list.New()
}

func (s *ds_string) append(n storage) {
	s.data.PushBack(n)
}

func (s *ds_string) fromString(str string) {
	s.new()
	for _, v := range str {
		s.append(storage{VAR_NUMBER, number{int(v), 1}})
	}
}
