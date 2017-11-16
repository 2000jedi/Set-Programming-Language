package main

import "container/list"

type array struct {
	data *list.List
}

func (s *array) get(n int) number {
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
	return p.Value.(number)
}

func (s *array) find(n number) bool {
	for e := s.data.Front(); e != nil; e = e.Next() {
		if e.Value == n {
			return true
		}
	}
	return false
}

func (s *array) new() {
	s.data = list.New()
}

func (s *array) remove(i int) {
	p := s.data.Front()
	for ; i > 0; p = p.Next() {
		if p == nil {
			panic("Index out of range")
		}
	}
	if p == nil {
		panic("Index out of range")
	}
	s.data.Remove(p)
}

func (s *array) append(n number) {
	s.data.PushBack(n)
}

func (s *array) toString() (ret string) {
	ret = "["
	var temp number
	for p := s.data.Front(); p != nil; p = p.Next() {
		temp = p.Value.(number)
		if p.Next() != nil {
			ret += temp.toString() + ", "
		} else {
			ret += temp.toString() + "]"
		}
	}
	return
}

func array_gen(data []storage, variable *Variable) *storage {
	var r array
	r.new()
	for _, v := range data {
		r.append(v.data.(number))
	}
	return &storage{VAR_ARRAY, r}
}
