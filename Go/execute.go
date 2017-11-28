package main

import (
	"fmt"
)

var stackTrace []string

func operation(op string, num1, num2 storage) *storage {
	if num1.vartype == VAR_NUMBER && num2.vartype == VAR_NUMBER {
		switch op {
		case "+":
			return &storage{VAR_NUMBER, add(num1.data.(number), num2.data.(number))}
		case "-":
			return &storage{VAR_NUMBER, sub(num1.data.(number), num2.data.(number))}
		case "*":
			return &storage{VAR_NUMBER, mul(num1.data.(number), num2.data.(number))}
		case "/":
			return &storage{VAR_NUMBER, div(num1.data.(number), num2.data.(number))}
		case "==":
			return &storage{VAR_NUMBER, equal(num1.data.(number), num2.data.(number))}
		case ">":
			return &storage{VAR_NUMBER, gt(num1.data.(number), num2.data.(number))}
		case ">=":
			return &storage{VAR_NUMBER, ge(num1.data.(number), num2.data.(number))}
		case "<":
			return &storage{VAR_NUMBER, lt(num1.data.(number), num2.data.(number))}
		case "<=":
			return &storage{VAR_NUMBER, le(num1.data.(number), num2.data.(number))}
		case "!=":
			if equal(num1.data.(number), num2.data.(number)) == True {
				return &storage{VAR_NUMBER, False}
			} else {
				return &storage{VAR_NUMBER, True}
			}
		case "+=":
			num1.data = add(num1.data.(number), num2.data.(number))
			return &num1
		case "-=":
			num1.data = sub(num1.data.(number), num2.data.(number))
			return &num1
		case "*=":
			num1.data = mul(num1.data.(number), num2.data.(number))
			return &num1
		case "/=":
			num1.data = div(num1.data.(number), num2.data.(number))
			return &num1
		default:
			panic("unknown operator for NUMBER and NUMBER: " + op)
		}
	}
	if num1.vartype == VAR_SET && num2.vartype == VAR_NUMBER {
		var temp set
		switch op {
		case "+":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(add(p.Value.(number), num2.data.(number)))
			}
			return &storage{VAR_SET, temp}
		case "+=":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(add(p.Value.(number), num2.data.(number)))
			}
			*num1.data.(set).data = *temp.data
			return &storage{VAR_SET, num1}
		case "-":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(sub(p.Value.(number), num2.data.(number)))
			}
			return &storage{VAR_SET, temp}
		case "-=":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(sub(p.Value.(number), num2.data.(number)))
			}
			*num1.data.(set).data = *temp.data
			return &storage{VAR_SET, num1}
		case "*":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(mul(p.Value.(number), num2.data.(number)))
			}
			return &storage{VAR_SET, temp}
		case "*=":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(mul(p.Value.(number), num2.data.(number)))
			}
			*num1.data.(set).data = *temp.data
			return &storage{VAR_SET, num1}
		case "/":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(div(p.Value.(number), num2.data.(number)))
			}
			return &storage{VAR_SET, temp}
		case "/=":
			temp.new()
			for p := num1.data.(set).data.Front(); p != nil; p = p.Next() {
				temp.data.PushBack(div(p.Value.(number), num2.data.(number)))
			}
			*num1.data.(set).data = *temp.data
			return &storage{VAR_SET, num1}
		default:
			panic("unknown operator for SET and NUMBER: " + op)
		}
	}
	panic("method not allowed: " + fmt.Sprint(num1.vartype, num1.data, num2.vartype, num2.data))
}

func evaluate(line []lexical, variable *Variable) *storage {
	if !*debug_flag {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("\nTraceback: ")
				for _, v := range stackTrace {
					fmt.Printf("    %s\n", v)
				}
				fmt.Println(r)
				fmt.Printf("\033[0m\n")
			}
		}()
	}

	var stack Stack
	i := 0
	for i < len(line) {
		switch line[i].fsm {
		case LEX_NUMBER:
			var n number
			n.construct(line[i].data)
			stack.Push(&storage{VAR_NUMBER, n})
		case LEX_STR:
			var str ds_string
			str.fromString(line[i].data)
			stack.Push(&storage{VAR_STRING, str})
		case LEX_EXPR:
			stack.Push(&storage{VAR_VAR, Var(line[i].data)})
		case LEX_OPR:
			temp, err := stack.DeVarPop(variable)
			if err != nil {
				panic(err)
			}
			ret, err := stack.DeVarPop(variable)
			if err != nil {
				panic(err)
			}
			ret_pointer := operation(line[i].data, *ret, *temp)
			if ret_pointer != nil {
				stack.Push(ret_pointer)
			}
		case LEX_NAMESPACE:
			temp, err := stack.Pop()
			if err != nil {
				panic(err)
			}
			ret, err := stack.DeVarPop(variable)
			if err != nil {
				panic(err)
			}
			class := ret.data.(Variable)
			val, err := class.get(temp.data.(Var))
			if err != nil {
				panic(err)
			}
			stack.Push(val)
		case LEX_SEPERATOR:
			stack.Push(&storage{VAR_FSM, Var(line[i].data)}) // TODO: improvement
		case LEX_BRACES:
			if line[i].data == "{" {
				is_function := 0
				var segment_data []lexical
				i++
				another_para := 0
				for line[i].data != "}" || another_para != 0 {
					if line[i].data == "{" {
						another_para++
					}
					if line[i].data == "}" {
						another_para--
					}
					segment_data = append(segment_data, line[i])
					if line[i].fsm == LEX_FUNC && another_para == 0 {
						is_function = len(segment_data)
					}
					i++
				}
				if is_function == 0 {
					is_function = 1
					segment_data = append([]lexical{{LEX_FUNC, ""}}, segment_data...)
				}
				var argv []storage
				for _, varname := range segment_data[:is_function-1 : 2] {
					argv = append(argv, storage{VAR_VAR, Var(varname.data)})
				}

				exprs := segment_data[is_function:]
				stack.Push(&storage{VAR_FUNCTION, function{argv, exprs}})
			}
		case LEX_CALL:
			var argc []storage
			if line[i].data == "(" {
				stack.Push(&storage{VAR_FSM, Var(line[i].data)}) // TODO: improvement
			} else {
				top, err := stack.Top()
				if err != nil {
					panic(err)
				}
				if top.vartype == VAR_FSM && top.data.toString() == "(" {
					stack.Pop()
				} else {
					temp, err := stack.DeVarPop(variable)
					if err != nil {
						panic(err)
					}
					argc = append(argc, *temp)
					top, err = stack.Top()
					if err != nil {
						panic(err)
					}
					for stack.Len() > 1 && top.data.toString() != "(" { // TODO: improvement
						stack.Pop()
						temp, err = stack.DeVarPop(variable)
						if err != nil {
							panic(err)
						}
						argc = append(argc, *temp)
						top, err = stack.Top()
						if err != nil {
							panic(err)
						}
					}
					stack.Pop()
				}
				lambda, err := stack.DeVarPop(variable)
				if err != nil {
					panic(err)
				}
				argc_ := make([]storage, len(argc), len(argc))
				for index, data := range argc {
					argc_[len(argc)-index-1] = data
				}
				ret := do_func(*lambda, argc_, variable)
				if ret != nil {
					stack.Push(ret)
				}
			}
		case LEX_ASSIGN:
			temp, err := stack.DeVarPop(variable)
			if err != nil {
				panic(err)
			}
			ret, err := stack.Pop()
			if err != nil {
				panic(err)
			}
			variable.set(ret.data.(Var), *temp)
			stack.Push(temp)
		case LEX_ADDR:
			if line[i].data == "[" {
				var addr []lexical
				i++
				for line[i].data != "]" {
					addr = append(addr, line[i])
					i++
				}
				ret := evaluate(addr, variable).data.(number)
				loc := ret.toInt()
				val, err := stack.DeVarPop(variable)
				if err != nil {
					panic(err)
				}
				arr := val.data.(array)
				stack.Push(arr.get(loc))
			}
		}
		i++
	}
	if stack.Len() != 0 {
		val, err := stack.Top()
		if err != nil {
			panic(err)
		}
		return val
	}
	return nil
}

func execute(lines [][]lexical, variable *Variable) {
	for k, line := range lines {
		stackTrace = []string{fmt.Sprintf("In line %d: ", k)}
		evaluate(line, variable)
	}
}
