package main

import "container/list"

type set struct {
	data *list.List
}

func (s set) find(n number) bool {
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
		if lt(p.Value.(number), n) == False {
			flag = false
			if equal(p.Value.(number), n) == False {
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

func (s set) toString() (ret string) {
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

func set_gen(data []storage, variable *Variable) *storage {
	var r set
	r.new()
	for _, v := range data {
		r.append(v.data.(number))
	}
	return &storage{VAR_SET, r}
}
