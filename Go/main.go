package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	fmt.Println()
	dat, err := ioutil.ReadFile("./debug.sp")
	if err != nil {
		panic(err)
	}
	file := strings.Split(string(dat), "\n")
	execute(syn_parse(lex_parse(file)))
}
