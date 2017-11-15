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

type inherit struct {
	function func(data []storage, variable *Variable) *storage
}
