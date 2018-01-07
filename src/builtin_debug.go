package main

import "fmt"

func invoke_debug() (v Variable) {
	v.init()
	v.add("printVarStack", storage{VAR_C_FUNCTION, CFunction{debug_printVarStack}})
	return v
}

func debug_printVarStack(data []storage, variable *Variable) *storage {
	fmt.Println(variable)
	return nil
}
