package main

import "fmt"

type Variable struct {
	stack map[string]*Stack
}

func (v *Variable) init() {
	v.stack = make(map[string]*Stack)
	invoke_inherit(v)
}

func (v *Variable) add(name string, val storage) {
	if _, ok := v.stack[name]; ok {
		v.stack[name].Push(&val)
	} else {
		v.stack[name] = &Stack{&val}
	}
}

func (v *Variable) get(name string) (val storage) {
	if _, ok := v.stack[name]; ok {
		val = *v.stack[name].Top()
	} else {
		panic("Variable Undefined: " + name)
	}
	return
}

func (v *Variable) del(name string) {
	if _, ok := v.stack[name]; ok {
		v.stack[name].Pop()
	} else {
		panic("Variable Undefined: " + name)
	}
}

func (v *Variable) set(name string, val storage) {
	if _, ok := v.stack[name]; ok {
		if v.stack[name].Len() > 0 {
			v.stack[name].Pop()
		}
		v.stack[name].Push(&val)
	} else {
		v.stack[name] = &Stack{&val}
	}
}

func operation(op string, num1, num2 storage) *storage {
	if num1.vartype == VAR_NUMBER && num2.vartype == VAR_NUMBER {
		var temp number
		switch op {
		case "+":
			temp = add(num1.data.(number), num2.data.(number))
			return &storage{VAR_NUMBER, temp}
		case "-":
			temp = sub(num1.data.(number), num2.data.(number))
			return &storage{VAR_NUMBER, temp}
		case "*":
			temp = mul(num1.data.(number), num2.data.(number))
			return &storage{VAR_NUMBER, temp}
		case "/":
			temp = div(num1.data.(number), num2.data.(number))
			return &storage{VAR_NUMBER, temp}
		case "==":
			if equal(num1.data.(number), num2.data.(number)) {
				return &True
			} else {
				return &False
			}
		case ">":
			if gt(num1.data.(number), num2.data.(number)) {
				return &True
			} else {
				return &False
			}
		case ">=":
			if ge(num1.data.(number), num2.data.(number)) {
				return &True
			} else {
				return &False
			}
		case "<":
			if lt(num1.data.(number), num2.data.(number)) {
				return &True
			} else {
				return &False
			}
		case "<=":
			if le(num1.data.(number), num2.data.(number)) {
				return &True
			} else {
				return &False
			}
		case "!=":
			if equal(num1.data.(number), num2.data.(number)) {
				return &False
			} else {
				return &True
			}
		case "+=":
			num1.data = add(num1.data.(number), num2.data.(number))
			return nil
		case "-=":
			num1.data = sub(num1.data.(number), num2.data.(number))
			return nil
		case "*=":
			num1.data = mul(num1.data.(number), num2.data.(number))
			return nil
		case "/=":
			num1.data = div(num1.data.(number), num2.data.(number))
			return nil
		default:
			panic("Unknown operator: " + op)
		}
	}
	if num1.vartype == VAR_SET && num2.vartype == VAR_NUMBER {
		var temp set
		switch op {
		case "+":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(p.Value)
			}
			temp.data = num1.data.(set).data
			temp.append(num2.data.(number))
			return &storage{VAR_SET, temp}
		case "-":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(p.Value)
			}
			temp.remove(num2.data.(number))
			return &storage{VAR_SET, temp}
		case "*":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(p.Value)
			}
			num := num2.data.(number)
			for p := temp.data.Front(); p != nil; p = p.Next() {
				p.Value = mul(p.Value.(number), num)
			}
			return &storage{VAR_SET, temp}
		case "/":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(p.Value)
			}
			num := num2.data.(number)
			for p := temp.data.Front(); p != nil; p = p.Next() {
				p.Value = div(p.Value.(number), num)
			}
			return &storage{VAR_SET, temp}
		default:
			panic("Unknown operator: " + op)
		}
	}
	if debug_flag {
		fmt.Println(num1.vartype, num1.data, num2.vartype, num2.data)
	}
	panic("Method not allowed")
}

func evaluate(line []storage, variable *Variable) *storage {
	if debug_flag {
		for _, v := range line {
			fmt.Println(v.vartype, v.data)
		}
	}
	var stack Stack
	i := 0
	if !debug_flag {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
				fmt.Printf("\033[0m\n")
			}
		}()
	}
	for i < len(line) {
		switch getlex(&line[i]).fsm {
		case LEX_NUMBER:
			var n number
			n.construct(getlex(&line[i]).data)
			stack.Push(&storage{VAR_NUMBER, n})
		case LEX_STR:
			stack.Push(&storage{VAR_STRING, getlex(&line[i]).data})
		case LEX_EXPR:
			stack.Push(&storage{VAR_EXPR, getlex(&line[i]).data})
		case LEX_OPR:
			temp := *stack.Pop()
			ret := *stack.Pop()
			if temp.vartype == VAR_EXPR {
				temp = variable.get(temp.data.(string))
			}
			if ret.vartype == VAR_EXPR {
				ret = variable.get(ret.data.(string))
			}
			ret_pointer := operation(getlex(&line[i]).data, ret, temp)
			if ret_pointer != nil {
				stack.Push(ret_pointer)
			}
		case LEX_NAMESPACE:
			temp := *stack.Pop()
			ret := *stack.Pop()
			namespace := variable.get(ret.data.(string)).data.(Variable)
			val := namespace.get(temp.data.(string))
			stack.Push(&val)
		case LEX_SEPERATOR:
			stack.Push(&line[i])
		case LEX_BRACES:
			if getlex(&line[i]).data == "{" {
				is_function := 0
				var segment_data []storage
				i++
				another_para := 0
				for getlex(&line[i]).data != "}" || another_para != 0 {
					if getlex(&line[i]).data == "{" {
						another_para++
					}
					if getlex(&line[i]).data == "}" {
						another_para--
					}
					segment_data = append(segment_data, line[i])
					if getlex(&line[i]).fsm == LEX_FUNC {
						is_function = len(segment_data)
					}
					i++
				}
				if is_function == 0 {
					panic("This method to create set is depreciated. Use set() function instead")
				}
				var argv []storage
				for _, varname := range segment_data[:is_function-1 : 2] {
					argv = append(argv, storage{VAR_EXPR, varname.data.(*lexical).data})
				}

				exprs := segment_data[is_function:]
				stack.Push(&storage{VAR_FUNCTION, function{argv, exprs}})
			}
		case LEX_CALL:
			var argc []storage
			if getlex(&line[i]).data == "(" {
				stack.Push(&line[i])
			} else {
				if stack.Top().vartype == VAR_FSM && getlex(stack.Top()).data == "(" {
					stack.Pop()
				} else {
					argc = append(argc, *stack.Pop())
					for stack.Len() > 1 && getlex(stack.Top()).fsm != LEX_CALL {
						stack.Pop()
						argc = append(argc, *stack.Pop())
					}
					stack.Pop()
					for j, _ := range argc {
						if argc[j].vartype == VAR_EXPR {
							argc[j] = variable.get(argc[j].data.(string))
						}
					}
				}
				lambda := variable.get(stack.Pop().data.(string))
				argc_ := make([]storage, len(argc), len(argc))
				for index, data := range argc {
					argc_[len(argc)-index-1] = data
				}
				ret := do_func(lambda, argc_, variable)
				if ret != nil {
					stack.Push(ret)
				}
			}
		case LEX_ASSIGN:
			temp := stack.Pop()
			if temp.vartype == VAR_EXPR {
				temp_ := variable.get(temp.data.(string))
				temp = &temp_
			}
			ret := stack.Pop()
			variable.set(ret.data.(string), *temp)
		case LEX_ADDR:
			if getlex(&line[i]).data == "[" {
				var addr []storage
				i++
				for getlex(&line[i]).data != "]" {
					addr = append(addr, line[i])
					i++
				}
				ret := evaluate(addr, variable).data.(number)
				loc := ret.toInt()
				arr := variable.get(stack.Pop().data.(string)).data.(array)
				stack.Push(&storage{VAR_NUMBER, arr.get(loc)})
			}
		}
		i++
	}
	if stack.Len() != 0 {
		s := *stack.Top()
		for s.vartype == VAR_EXPR {
			s = variable.get(s.data.(string))
		}
		return &s
	}
	return nil
}

func execute(lines [][]storage, variable *Variable) {
	for _, line := range lines {
		evaluate(line, variable)
	}
}
