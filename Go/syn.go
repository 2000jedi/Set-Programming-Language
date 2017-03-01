package main

var priority map[string]int

func init() {
	priority = make(map[string]int)
	priority["="] = 0
	priority["{"] = 0
	priority["&&"] = 1
	priority["||"] = 1
	priority["<"] = 2
	priority[">"] = 2
	priority["<="] = 2
	priority[">="] = 2
	priority["!="] = 2
	priority["=="] = 2
	priority["+"] = 3
	priority["-"] = 3
	priority["*"] = 4
	priority["/"] = 4
	priority[":"] = 5
	priority["("] = 0
	priority["["] = 0
}

func getlex(temp *storage) lexical {
	return *temp.data.(*lexical)
}

func segment(line Stack) (prog []storage) {
	var stack Stack
	i := 0
	for i < len(line) {
		switch getlex(line[i]).fsm {
		case lex_fsm["expr"], lex_fsm["number"], lex_fsm["str"]:
			prog = append(prog, *line[i])
		case lex_fsm["func"]:
			prog = append(prog, *line[i])
		case lex_fsm["bracket"], lex_fsm["call"], lex_fsm["addr"]:
			stack.Push(line[i])
			if getlex(stack.Top()).fsm == lex_fsm["call"] {
				prog = append(prog, *stack.Top())
			}
		case lex_fsm["end_bracket"]:
			for !(getlex(stack.Top()).fsm == lex_fsm["bracket"] || getlex(stack.Top()).fsm == lex_fsm["call"]) {
				prog = append(prog, *stack.Pop())
			}
			if getlex(stack.Top()).fsm == lex_fsm["call"] {
				prog = append(prog, storage{var_fsm["fsm"], &lexical{lex_fsm["call"], ")"}})
			}
			stack.Pop()
		case lex_fsm["end_brace"]:
			for getlex(stack.Top()).fsm != lex_fsm["addr"] {
				prog = append(prog, *stack.Pop())
			}
			prog = append(prog, *stack.Pop())
		case lex_fsm["set"]:
			switch getlex(line[i]).data {
			case "{":
				prog = append(prog, *line[i])
				stack.Push(line[i])
			case "}":
				for getlex(stack.Top()).data != "{" {
					prog = append(prog, *stack.Pop())
				}
				stack.Pop()
				prog = append(prog, *line[i])
			default:
				panic("Wrong literal " + getlex(line[i]).data + "\n")
			}
		case lex_fsm["seperator"]:
			prog = append(prog, *line[i])
			temp_line := Stack{}
			i++
			para := 0
			for !(getlex(line[i]).fsm == lex_fsm["seperator"] || getlex(line[i]).fsm == lex_fsm["vec"] || getlex(line[i]).fsm == lex_fsm["set"]) && (getlex(line[i]).fsm != lex_fsm["end_bracket"] || para != 0) {
				temp_line.Push(line[i])
				if getlex(line[i]).data == "(" {
					para++
				}
				if getlex(line[i]).data == ")" {
					para--
				}
				i++
			}
			i--
			for _, t := range segment(temp_line) {
				prog = append(prog, t)
			}
		default:
			if stack.Len() == 0 || priority[getlex(stack.Top()).data] < priority[getlex(line[i]).data] {
				stack.Push(line[i])
			} else {
				for !(stack.Len() == 0) && priority[getlex(stack.Top()).data] >= priority[getlex(line[i]).data] {
					prog = append(prog, *stack.Pop())
				}
				stack.Push(line[i])
			}
		}
		i++
	}
	for !(stack.Len() == 0) {
		prog = append(prog, *stack.Pop())
	}
	return
}

func syn_parse(lines []Stack) (prog [][]storage) {
	for _, line := range lines {
		prog = append(prog, segment(line))
	}
	return
}
