package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var debug_flag bool

func interactive() {
	fmt.Println("SPL 17.11.17 (build on Go 1.9)")
	fmt.Println("Type exit() to exit")
	reader := bufio.NewReader(os.Stdin)

	if !debug_flag {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
	}
	var variable Variable
	variable.init()
	for {
		fmt.Print(">>> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if text == "exit()\n" {
			panic("User Exit")
		}
		file := strings.Split(text, "\n")
		execute(syn_parse(lex_parse(file)), &variable)
	}
}

func runfile(s string) Variable {
	dat, err := ioutil.ReadFile(s)
	if err != nil {
		panic(err)
	}
	file := strings.Split(string(dat), "\n")
	var variable Variable
	variable.init()
	execute(syn_parse(lex_parse(file)), &variable)
	return variable
}

func help() {
	fmt.Println("Help of SPL:")
	fmt.Println("Specify a filename or use interactive shell")
	fmt.Println("Example:")
	fmt.Println("spl debug.sp")
}

func main() {
	debug_flag = true
	if len(os.Args) == 1 {
		interactive()
	} else {
		if os.Args[1] == "--help" {
			help()
		} else {
			if os.Args[1] == "--debug" {
				debug_flag = true
			}
			runfile(os.Args[1])
		}
	}
}
