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
	priority["("] = 0
}

func segment(line Lexs) (prog []lexical) {
	var stack Lexs
	i := 0
	for i < len(line) {
		switch line[i].fsm {
		case LEX_EXPR, LEX_NUMBER, LEX_STR:
			prog = append(prog, line[i])
		case LEX_FUNC:
			prog = append(prog, line[i])
		case LEX_NAMESPACE:
			if line[i+1].fsm != LEX_EXPR {
				panic(line[i+1].data + " is not an expression")
			}
			prog = append(prog, line[i+1])
			prog = append(prog, line[i])
			i++
		case LEX_CALL:
			if stack.Len() > 0 && stack.Top().fsm == LEX_NAMESPACE {
				prog = append(prog, line[i])
			}
			stack.Push(line[i])
			prog = append(prog, stack.Top())
		case LEX_ADDR:
			switch line[i].data {
			case "[":
				prog = append(prog, line[i])
				stack.Push(line[i])
			case "]":
				for !(stack.Top().fsm == LEX_ADDR && stack.Top().data == "[") {
					prog = append(prog, stack.Pop())
				}
				prog = append(prog, lexical{LEX_ADDR, "]"})
				stack.Pop()
			}
		case LEX_BRACKET:
			stack.Push(line[i])
		case LEX_END_BRACKET:
			for !(stack.Top().fsm == LEX_BRACKET || stack.Top().fsm == LEX_CALL) {
				prog = append(prog, stack.Pop())
			}
			if stack.Top().fsm == LEX_CALL {
				prog = append(prog, lexical{LEX_CALL, ")"})
			}
			stack.Pop()
		case LEX_BRACES:
			switch line[i].data {
			case "{":
				prog = append(prog, line[i])
				stack.Push(line[i])
			case "}":
				for stack.Top().data != "{" {
					prog = append(prog, stack.Pop())
				}
				stack.Pop()
				prog = append(prog, line[i])
			default:
				panic("Wrong literal " + line[i].data + "\n")
			}
		case LEX_SEPERATOR:
			for !(stack.Top().fsm == LEX_BRACKET || stack.Top().fsm == LEX_CALL) {
				prog = append(prog, stack.Pop())
			}
			prog = append(prog, line[i])
		default:
			if stack.Len() == 0 || priority[stack.Top().data] < priority[line[i].data] {
				stack.Push(line[i])
			} else {
				for !(stack.Len() == 0) && priority[stack.Top().data] >= priority[line[i].data] {
					prog = append(prog, stack.Pop())
				}
				stack.Push(line[i])
			}
		}
		i++
	}
	for !(stack.Len() == 0) {
		prog = append(prog, stack.Pop())
	}
	return
}

func syn_parse(lines []Lexs) (prog [][]lexical) {
	for _, line := range lines {
		prog = append(prog, segment(line))
	}
	return
}
