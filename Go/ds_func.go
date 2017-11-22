package main

type function struct {
	argv  []storage
	exprs []lexical
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

func do_func(lambda storage, argc []storage, variable *Variable) (s *storage) {
	if lambda.vartype == VAR_FUNCTION {
		f := lambda.data.(function)
		stackTrace = append(stackTrace, "In function")
		s = f.function(argc, variable)
		stackTrace = stackTrace[:len(stackTrace)-1]
	} else if lambda.vartype == VAR_C_FUNCTION {
		f := lambda.data.(c_function)
		stackTrace = append(stackTrace, "In function:")
		s = f.function(argc, variable)
		stackTrace = stackTrace[:len(stackTrace)-1]
	} else {
		panic(ERR_NOT_FUNC)
	}
	return
}

const (
	ERR_NOT_FUNC   = "not a function"
	ERR_FUNC_PRINT = "cannot print function object"
)
