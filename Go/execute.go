package main

import "fmt"

type Variable struct {
	stack map[string]*Stack
}

func (v *Variable) add(name string, val storage) {
	if _, ok := v.stack[name]; ok {
		v.stack[name].Push(&val)
	} else {
		v.stack[name] = &Stack{&val}
		//v.stack[name].Push(&val)
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

var variable Variable

func operation(op string, num1, num2 storage) *storage {
	if num1.vartype == var_fsm["number"] && num2.vartype == var_fsm["number"] {
		var temp number
		switch op {
		case "+":
			temp = add(num1.data.(number), num2.data.(number))
			return &storage{var_fsm["number"], temp}
		case "-":
			temp = sub(num1.data.(number), num2.data.(number))
			return &storage{var_fsm["number"], temp}
		case "*":
			temp = mul(num1.data.(number), num2.data.(number))
			return &storage{var_fsm["number"], temp}
		case "/":
			temp = div(num1.data.(number), num2.data.(number))
			return &storage{var_fsm["number"], temp}
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
	} else {
		panic("Method not allowed")
	}
}

func evaluate(line []storage) *storage {
	var stack Stack
	i := 0
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\n\033[1;31mError in char: %d ->", i)
			fmt.Println(line[i])
			fmt.Println("Debug Detail:")
			for _, val := range line {
				fmt.Print(val.data)
			}
			fmt.Println()
			fmt.Println(r)
			fmt.Printf("\033[0m\n")
		}
	}()
	for i < len(line) {
		switch getlex(&line[i]).fsm {
		case lex_fsm["number"]:
			var n number
			n.construct(getlex(&line[i]).data)
			stack.Push(&storage{var_fsm["number"], n})
		case lex_fsm["str"]:
			stack.Push(&storage{var_fsm["string"], getlex(&line[i]).data})
		case lex_fsm["expr"]:
			stack.Push(&storage{var_fsm["expr"], getlex(&line[i]).data})
		case lex_fsm["opr"]:
			temp := *stack.Pop()
			ret := *stack.Pop()
			if temp.vartype == var_fsm["expr"] {
				temp = variable.get(temp.data.(string))
			}
			if ret.vartype == var_fsm["expr"] {
				ret = variable.get(ret.data.(string))
			}
			ret_pointer := operation(getlex(&line[i]).data, ret, temp)
			if ret_pointer != nil {
				stack.Push(ret_pointer)
			}
		case lex_fsm["seperator"]:
			stack.Push(&line[i])
		case lex_fsm["set"]:
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
					if getlex(&line[i]).fsm == lex_fsm["func"] {
						is_function = len(segment_data)
					}
					i++
				}
				if is_function != 0 {
					var argv []storage
					for _, varname := range segment_data[:is_function-1 : 2] {
						argv = append(argv, storage{var_fsm["expr"], varname.data.(*lexical).data})
					}

					exprs := segment_data[is_function:]
					stack.Push(&storage{var_fsm["function"], function{argv, exprs}})
				} else {
					segment_data = append(segment_data, storage{var_fsm["fsm"], &lexical{lex_fsm["seperator"], ","}})
					var ret set
					for len(segment_data) != 0 {
						j := 0
						for getlex(&segment_data[j]).data != "," {
							j++
						}
						split := segment_data[:j]
						split_segment := segment_data[j+1:]
						split_seperated := make([]storage, len(split), len(split))
						copy(split_seperated, split)
						segment_data = make([]storage, len(split_segment), len(split_segment))
						copy(segment_data, split_segment)
						ret.append(evaluate(split).data.(number))
					}
					stack.Push(&storage{var_fsm["set"], ret})
				}
			}
		case lex_fsm["call"]:
			var argc []storage
			if getlex(&line[i]).data == "(" {
				stack.Push(&line[i])
			} else {
				if stack.Top().vartype == var_fsm["fsm"] && getlex(stack.Top()).data == "(" {
					stack.Pop()
				} else {
					argc = append(argc, *stack.Pop())
					for stack.Len() > 1 && getlex(stack.Top()).fsm != lex_fsm["call"] {
						stack.Pop()
						argc = append(argc, *stack.Pop())
					}
					stack.Pop()
					for j, _ := range argc {
						if argc[j].vartype == var_fsm["expr"] {
							argc[j] = variable.get(argc[j].data.(string))
						}
					}
				}
				lambda := variable.get(stack.Pop().data.(string))
				argc_ := make([]storage, len(argc), len(argc))
				for index, data := range argc {
					argc_[len(argc)-index-1] = data
				}
				ret := do_func(lambda, argc_)
				if ret != nil {
					stack.Push(ret)
				}
			}
		case lex_fsm["assign"]:
			temp := stack.Pop()
			ret := stack.Pop()
			variable.set(ret.data.(string), *temp)
		}
		i++
	}
	if stack.Len() != 0 {
		return stack.Top()
	}
	return nil
}

func execute(lines [][]storage) {
	for _, line := range lines {
		evaluate(line)
	}
}

func init() {
	variable.stack = make(map[string]*Stack)
}
