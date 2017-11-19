package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

var debug_flag bool

func interactive() {
	fmt.Println("SPL 17.11.17 (build on Go 1.9)")
	fmt.Println("Type exit() to exit")

	if !debug_flag {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
	}

	l, err := readline.NewEx(&readline.Config{
		Prompt:            ">>> ",
		HistoryFile:       "/tmp/spl.tmp",
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit()",
		HistorySearchFold: true,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	var variable Variable
	variable.init()

	for {
		text, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(text) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
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
	debug_flag = false
	if len(os.Args) == 1 {
		interactive()
	} else {
		if os.Args[1] == "--help" {
			help()
		} else {
			if os.Args[1] == "--debug" {
				debug_flag = true
				interactive()
			} else {
				runfile(os.Args[1])
			}
		}
	}
}
