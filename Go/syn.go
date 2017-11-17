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

func getlex(temp *storage) lexical {
	return *temp.data.(*lexical)
}

func segment(line Stack) (prog []storage) {
	var stack Stack
	i := 0
	for i < len(line) {
		switch getlex(line[i]).fsm {
		case LEX_EXPR, LEX_NUMBER, LEX_STR:
			prog = append(prog, *line[i])
		case LEX_FUNC:
			prog = append(prog, *line[i])
		case LEX_NAMESPACE:
			if getlex(line[i+1]).fsm != LEX_EXPR {
				panic(line[i+1].data.toString() + " is not an expression")
			}
			prog = append(prog, *line[i+1])
			prog = append(prog, *line[i])
			i++
		case LEX_CALL:
			if stack.Len() > 0 && getlex(stack.Top()).fsm == LEX_NAMESPACE {
				prog = append(prog, *line[i])
			}
			stack.Push(line[i])
			prog = append(prog, *stack.Top())
		case LEX_ADDR:
			switch getlex(line[i]).data {
			case "[":
				prog = append(prog, *line[i])
				stack.Push(line[i])
			case "]":
				for !(getlex(stack.Top()).fsm == LEX_ADDR && getlex(stack.Top()).data == "[") {
					prog = append(prog, *stack.Pop())
				}
				prog = append(prog, storage{VAR_FSM, &lexical{LEX_ADDR, "]"}})
				stack.Pop()
			}
		case LEX_BRACKET:
			stack.Push(line[i])
		case LEX_END_BRACKET:
			for !(getlex(stack.Top()).fsm == LEX_BRACKET || getlex(stack.Top()).fsm == LEX_CALL) {
				prog = append(prog, *stack.Pop())
			}
			if getlex(stack.Top()).fsm == LEX_CALL {
				prog = append(prog, storage{VAR_FSM, &lexical{LEX_CALL, ")"}})
			}
			stack.Pop()
		case LEX_BRACES:
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
		case LEX_SEPERATOR:
			for !(getlex(stack.Top()).fsm == LEX_BRACKET || getlex(stack.Top()).fsm == LEX_CALL) {
				prog = append(prog, *stack.Pop())
			}
			prog = append(prog, *line[i])
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
