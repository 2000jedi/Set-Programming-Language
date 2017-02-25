package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func interactive() {
	fmt.Println("SPL 0.1 (build on Go1.8)")
	fmt.Println("Type exit() to exit")
	reader := bufio.NewReader(os.Stdin)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
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
		execute(syn_parse(lex_parse(file)))
	}
}

func debug(s string) {
	dat, err := ioutil.ReadFile(s)
	if err != nil {
		panic(err)
	}
	file := strings.Split(string(dat), "\n")
	execute(syn_parse(lex_parse(file)))
}

func help() {
	fmt.Println("Help of SPL:")
	fmt.Println("Specify a filename or use interactive shell")
	fmt.Println("Example:")
	fmt.Println("spl debug.sp")
}

func main() {
	if len(os.Args) == 1 {
		interactive()
	} else {
		if os.Args[1] == "--help" {
			help()
		} else {
			debug(os.Args[1])
		}
	}
}
