package main

import (
	"errors"
)

type Function struct {
	argv  []storage
	exprs []lexical
}

func (f *Function) function(vals []storage, variable *Variable) (ret_val *storage) {
	for i, val := range f.argv {
		variable.add(val.data.(Var), vals[i])
	}
	ret_val = evaluate(f.exprs, variable)
	for _, val := range f.argv {
		if err := variable.del(val.data.(Var)); err != nil {
			panic(err)
		}
	}
	return
}

func (f Function) toString() string {
	panic(ERR_FUNC_PRINT)
}

type CFunction struct {
	function func(data []storage, variable *Variable) *storage
}

func (f CFunction) toString() string {
	panic(ERR_FUNC_PRINT)
}

func do_func(lambda storage, argc []storage, variable *Variable) (s *storage, err error) {
	if lambda.vartype == VAR_FUNCTION {
		f := lambda.data.(Function)
		stackTrace = append(stackTrace, "In function:")
		s = f.function(argc, variable)
		stackTrace = stackTrace[:len(stackTrace)-1]
	} else if lambda.vartype == VAR_C_FUNCTION {
		f := lambda.data.(CFunction)
		stackTrace = append(stackTrace, "In function:")
		s = f.function(argc, variable)
		stackTrace = stackTrace[:len(stackTrace)-1]
	} else {
		s, err = nil, errors.New(ERR_NOT_FUNC)
		return
	}
	err = nil
	return
}

const (
	ERR_NOT_FUNC   = "not a function"
	ERR_FUNC_PRINT = "cannot print function object"
)
