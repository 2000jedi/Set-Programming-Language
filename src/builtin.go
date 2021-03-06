package main

import (
	"fmt"
	"os"
)

func invoke_builtin(v *Variable) {
	v.add("println", storage{VAR_C_FUNCTION, CFunction{builtin_println}})
	v.add("print", storage{VAR_C_FUNCTION, CFunction{builtin_printf}})
	v.add("if", storage{VAR_C_FUNCTION, CFunction{builtin_if}})
	v.add("for", storage{VAR_C_FUNCTION, CFunction{builtin_for}})
	v.add("range", storage{VAR_C_FUNCTION, CFunction{builtin_range}})
	v.add("import", storage{VAR_C_FUNCTION, CFunction{builtin_import}})
	v.add("array", storage{VAR_C_FUNCTION, CFunction{array_gen}})
	v.add("set", storage{VAR_C_FUNCTION, CFunction{set_gen}})
	v.add("exit", storage{VAR_C_FUNCTION, CFunction{builtin_exit}})
}

func builtin_printf(data []storage, variable *Variable) *storage {
	for k, v := range data {
		if k != len(data)-1 {
			fmt.Print(v.data.toString() + " ")
		} else {
			fmt.Print(v.data.toString())
		}
	}
	newlineFlag = false
	return nil
}

func builtin_println(data []storage, variable *Variable) *storage {
	builtin_printf(data, variable)
	fmt.Println()
	newlineFlag = true
	return nil
}

func builtin_if(data []storage, variable *Variable) *storage {
	if len(data) > 3 || len(data) < 2 {
		panic(ERR_ARG_NUM)
	}
	cond := data[0].data.(number) != False
	if len(data) == 2 {
		branch_then := data[1]
		if cond {
			if ans, err := do_func(branch_then, []storage{}, variable); err != nil {
				panic(err)
			} else {
				return ans
			}
		}
	} else {
		branch_then := data[1]
		branch_else := data[2]
		if cond {
			if ans, err := do_func(branch_then, []storage{}, variable); err != nil {
				panic(err)
			} else {
				return ans
			}
		} else {
			if ans, err := do_func(branch_else, []storage{}, variable); err != nil {
				panic(err)
			} else {
				return ans
			}
		}
	}
	return nil
}

func builtin_for(data []storage, variable *Variable) *storage {
	if len(data) > 2 {
		panic(ERR_ARG_NUM)
	}
	if data[0].vartype != VAR_SET || (data[1].vartype != VAR_C_FUNCTION && data[1].vartype != VAR_FUNCTION) {
		panic(ERR_ARG_TYPE)
	}
	a := data[0].data.(set).data
	b := data[1]

	return_set := set{}
	return_set.new()
	for p := a.Front(); p != nil; p = p.Next() {
		var stg []storage
		stg = append(stg, storage{VAR_NUMBER, p.Value.(number)})
		if ans, err := do_func(b, stg, variable); err != nil {
			panic(err)
		} else {
			if ans != nil {
				return_set.append(ans.data.(number))
			}
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
		panic(ERR_ARG_NUM)
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
		panic(ERR_ARG_NUM)
	}
	v := (data[0].data.toString())
	if len(v) >= 3 && v[len(v)-3:] == ".sp" {
		return &storage{VAR_CLASS, runfile(v)}
	} else {
		return &storage{VAR_CLASS, builtin_imports[v]}
	}
}

func builtin_exit(data []storage, variable *Variable) *storage {
	os.Exit(0)
	return nil
}

func init() {
	builtin_imports = map[string]Variable{
		"debug": invoke_debug(),
		"os":    invoke_os(),
	}
}

const (
	ERR_ARG_NUM  = "wrong number of arguments"
	ERR_ARG_TYPE = "wrong type of arguments"
)
