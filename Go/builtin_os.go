// builtin_os
package main

func invoke_os() (v Variable) {
	v.add("open", storage{VAR_C_FUNCTION, CFunction{os_open}})
	return
}

func os_open(data []storage, v *Variable) *storage {
	return nil
}
