package main

type function struct {
	argv  []storage
	exprs []storage
}

func (f *function) function(vals []storage, variable *Variable) (ret_val *storage) {
	for i, val := range f.argv {
		variable.add(val.data.(Var), vals[i])
	}
	ret_val = evaluate(f.exprs, variable)
	for _, val := range f.argv {
		variable.del(val.data.(Var))
	}
	return
}

func (f function) toString() string {
	panic(ERR_FUNC_PRINT)
}

type c_function struct {
	function func(data []storage, variable *Variable) *storage
}

func (f c_function) toString() string {
	panic(ERR_FUNC_PRINT)
}

func do_func(lambda storage, argc []storage, variable *Variable) *storage {
	if lambda.vartype == VAR_FUNCTION {
		f := lambda.data.(function)
		return f.function(argc, variable)
	} else if lambda.vartype == VAR_C_FUNCTION {
		f := lambda.data.(c_function)
		return f.function(argc, variable)
	} else {
		panic(ERR_NOT_FUNC)
	}
}

const (
	ERR_NOT_FUNC   = "not a function"
	ERR_FUNC_PRINT = "cannot print function object"
)
