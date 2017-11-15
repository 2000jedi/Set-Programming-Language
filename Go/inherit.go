package main

import "fmt"

func invoke_inherit(v *Variable){
  v.add("println", storage{var_fsm["c_function"], inherit{inherit_println}})
	v.add("print", storage{var_fsm["c_function"], inherit{inherit_printf}})
	v.add("for", storage{var_fsm["c_function"], inherit{inherit_for}})
	v.add("range", storage{var_fsm["c_function"], inherit{inherit_range}})
	v.add("import", storage{var_fsm["c_function"], inherit{inherit_import}})
	v.add("array", storage{var_fsm["c_function"], inherit{inherit_gen_array}})
  v.add("set", storage{var_fsm["c_function"], inherit{inherit_gen_set}})
}

func inherit_printf(data []storage, variable *Variable) *storage {
	for _, val := range data {
		switch val.vartype {
		case var_fsm["number"]:
			num := val.data.(number)
			fmt.Print(num.toString())
		case var_fsm["set"]:
			num := val.data.(set)
			fmt.Print(num.toString())
		case var_fsm["array"]:
			num := val.data.(array)
			fmt.Print(num.toString())
		default:
			fmt.Print(val.data)
		}
	}
	return nil
}

func inherit_println(data []storage, variable *Variable) *storage {
	inherit_printf(data, variable)
	fmt.Println()
	return nil
}

func inherit_for(data []storage, variable *Variable) *storage {
	if len(data) > 2 {
		panic("Wrong number of arguments")
	}
	if data[0].vartype != var_fsm["set"] || (data[1].vartype != var_fsm["c_function"] && data[1].vartype != var_fsm["function"]) {
		panic("Wrong type of arguments")
	}
	a := data[0].data.(set).data
	b := data[1]

	return_set := set{}
	return_set.new()
	for p := a.Front(); p != nil; p = p.Next() {
		var stg []storage
		stg = append(stg, storage{var_fsm["number"], p.Value.(number)})
		ret := do_func(b, stg, variable)
		if ret != nil {
			return_set.append(ret.data.(number))
		}
	}
	return &storage{var_fsm["set"], return_set}
}

func inherit_range(data []storage, variable *Variable) *storage {
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
	return &storage{var_fsm["set"], ret}
}

func inherit_import(data []storage, variable *Variable) *storage {
	if len(data) > 1 {
		panic("Wrong number of arguments")
	}
	var_ := runfile(data[0].data.(string) + ".sp")
	return &storage{var_fsm["namespace"], var_}
}

func inherit_gen_array(data []storage, variable *Variable) *storage {
	var r array
	r.new()
	for _, v := range data {
		r.append(v.data.(number))
	}
	return &storage{var_fsm["array"], r}
}

func inherit_gen_set(data []storage, variable *Variable) *storage {
  var r set
  r.new()
  for _, v := range data {
    r.append(v.data.(number));
  }
  return &storage{var_fsm["set"], r}
}
