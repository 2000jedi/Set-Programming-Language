package main

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

type c_function struct {
	function func(data []storage, variable *Variable) *storage
}

func do_func(lambda storage, argc []storage, variable *Variable) *storage {
	if lambda.vartype == VAR_FUNCTION {
		f := lambda.data.(function)
		return f.function(argc, variable)
	} else if lambda.vartype == VAR_C_FUNCTION {
		f := lambda.data.(c_function)
		return f.function(argc, variable)
	} else {
		panic("not a function")
	}
}
