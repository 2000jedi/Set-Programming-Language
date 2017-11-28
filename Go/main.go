package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/chzyer/readline"
)

var debugFlag *bool
var helpFlag *bool
var addr *string

func interactive() {
	fmt.Println("SPL 17.11.17 (build on Go 1.9)")
	fmt.Println("Type exit() to exit")

	l, err := readline.NewEx(&readline.Config{
		Prompt:            ">>> ",
		HistoryFile:       "/tmp/spl.tmp",
		InterruptPrompt:   "^C",
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
		execute(synParse(lexParse(file)), &variable)
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
	execute(synParse(lexParse(file)), &variable)
	return variable
}

func main() {
	if *addr == "" {
		interactive()
	} else {
		runfile(*addr)
	}

}

func init() {
	// define & parse flags
	debugFlag = flag.Bool("debug", false, "Enable Debug Mode")
	addr = flag.String("d", "", "Execute a file")
	flag.Parse()
}
