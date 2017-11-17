package main

import "fmt"

func invoke_builtin(v *Variable) {
	v.add("println", storage{VAR_C_FUNCTION, c_function{builtin_println}})
	v.add("print", storage{VAR_C_FUNCTION, c_function{builtin_printf}})
	v.add("if", storage{VAR_C_FUNCTION, c_function{builtin_if}})
	v.add("for", storage{VAR_C_FUNCTION, c_function{builtin_for}})
	v.add("range", storage{VAR_C_FUNCTION, c_function{builtin_range}})
	v.add("import", storage{VAR_C_FUNCTION, c_function{builtin_import}})
	v.add("array", storage{VAR_C_FUNCTION, c_function{array_gen}})
	v.add("set", storage{VAR_C_FUNCTION, c_function{set_gen}})
}

func builtin_printf(data []storage, variable *Variable) *storage {
	for k, v := range data {
		switch v.vartype {
		case VAR_NUMBER:
			num := v.data.(number)
			if k != len(data)-1 {
				fmt.Print(num.toString(), " ")
			} else {
				fmt.Print(num.toString())
			}
		case VAR_SET:
			num := v.data.(set)
			if k != len(data)-1 {
				fmt.Print(num.toString(), " ")
			} else {
				fmt.Print(num.toString())
			}
		case VAR_ARRAY:
			num := v.data.(array)
			if k != len(data)-1 {
				fmt.Print(num.toString(), " ")
			} else {
				fmt.Print(num.toString())
			}
		default:
			if k != len(data)-1 {
				fmt.Print(v.data, " ")
			} else {
				fmt.Print(v.data)
			}
		}
	}
	return nil
}

func builtin_println(data []storage, variable *Variable) *storage {
	builtin_printf(data, variable)
	fmt.Println()
	return nil
}

func builtin_if(data []storage, variable *Variable) *storage {
	if len(data) > 3 || len(data) < 2 {
		panic("Wrong number of arguments")
	}
	cond := data[0].data.(number) != False // Judge whether the condition is true
	if len(data) == 2 {
		branch_then := data[1]
		if cond {
			return do_func(branch_then, []storage{}, variable)
		}
	} else {
		branch_then := data[1]
		branch_else := data[2]
		if cond {
			return do_func(branch_then, []storage{}, variable)
		} else {
			return do_func(branch_else, []storage{}, variable)
		}
	}
	return nil
}

func builtin_for(data []storage, variable *Variable) *storage {
	if len(data) > 2 {
		panic("Wrong number of arguments")
	}
	if data[0].vartype != VAR_SET || (data[1].vartype != VAR_C_FUNCTION && data[1].vartype != VAR_FUNCTION) {
		panic("Wrong type of arguments")
	}
	a := data[0].data.(set).data
	b := data[1]

	return_set := set{}
	return_set.new()
	for p := a.Front(); p != nil; p = p.Next() {
		var stg []storage
		stg = append(stg, storage{VAR_NUMBER, p.Value.(number)})
		ret := do_func(b, stg, variable)
		if ret != nil {
			return_set.append(ret.data.(number))
		}
	}
	return &storage{VAR_SET, return_set}
}

func builtin_range(data []storage, variable *Variable) *storage {
	init := 0
	end := 0
	step := 1
	var temp number
	switch len(data) {
	case 2:
		temp = data[0].data.(number)
		init = temp.toInt()
		fallthrough
	case 1:
		temp = data[len(data)-1].data.(number)
		end = temp.toInt()
	case 3:
		temp = data[0].data.(number)
		init = temp.toInt()
		temp = data[1].data.(number)
		end = temp.toInt()
		temp = data[2].data.(number)
		step = temp.toInt()
	default:
		panic("Wrong number of arguments")
	}
	var ret set
	ret.new()
	for i := init; i < end; i += step {
		ret.append(number{i, 1})
	}
	return &storage{VAR_SET, ret}
}

func builtin_import(data []storage, variable *Variable) *storage {
	if len(data) > 1 {
		panic("Wrong number of arguments")
	}
	var_ := runfile(data[0].data.(string))
	return &storage{VAR_CLASS, var_}
}
