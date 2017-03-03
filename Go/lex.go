package main

func push(lexs *Stack, fsm string, data string) {
	lexs.Push(&storage{var_fsm["fsm"], &lexical{lex_fsm[fsm], data}})
}

func pop(lexs Stack) lexical {
	return *lexs.Pop().data.(*lexical)
}

func top(lexs Stack) lexical {
	return *lexs.Top().data.(*lexical)
}

func process_escapes(s string) string {
	if s == "\\n" {
		return "\n"
	} else {
		return "\\"
	}
}

func is_expr_letters(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '_'
}

func is_expr_letter_latter(c byte) bool {
	return '0' <= c && c <= '9'
}

func is_number(c byte) bool {
	return ('0' <= c && c <= '9') || c == '.'
}

func not_opr(c int) bool {
	return c != var_fsm["number"] && c != var_fsm["set"] && c != var_fsm["addr"] && c != var_fsm["expr"]
}

func lex_parse(lines []string) (lex_lines []Stack) {
	for _, line := range lines {
		var lexs Stack
		line = line + " "
		i := 0
		for i < len(line) {
			switch line[i] {
			case '\n':
			case ' ':
			case '#':
				i = len(line)
			case '"':
				temp := ""
				i++
				for line[i] != '"' {
					if line[i] == '\\' {
						temp += process_escapes(string(line[i]) + string(line[i+1]))
						i += 2
					} else {
						temp += string(line[i])
						i++
					}
				}
				push(&lexs, "str", temp)
			case '=':
				if line[i+1] == '=' {
					push(&lexs, "opr", "==")
					i++
				} else {
					push(&lexs, "assign", "=")
				}
			case '+', '*', '/', '<', '>':
				if line[i+1] == '=' {
					push(&lexs, "opr", string(line[i])+"=")
					i++
				} else {
					push(&lexs, "opr", string(line[i]))
				}
			case '-':
				if is_number(line[i+1]) && not_opr(lexs.Top().data.(*lexical).fsm) {
					temp := string(line[i])
					i++
					for is_number(line[i]) {
						temp += string(line[i])
						i++
					}
					i--
					push(&lexs, "number", temp)
				} else {
					if line[i+1] == '=' {
						push(&lexs, "opr", string(line[i])+"=")
						i++
					} else {
						push(&lexs, "opr", string(line[i]))
					}
				}
			case '!':
				if line[i+1] == '=' {
					push(&lexs, "opr", "!=")
					i++
				} else {
					panic("Illegal literal '!'\n")
				}
			case '|':
				if line[i+1] == '|' {
					push(&lexs, "opr", "||")
					i++
				} else {
					push(&lexs, "opr", "|")
				}
			case '&':
				if line[i+1] == '&' {
					push(&lexs, "opr", "&&")
					i++
				} else {
					panic("Illegal literal '&'\n")
				}
			case '~':
				push(&lexs, "func", "~")
			case ':':
				push(&lexs, "namespace", ":")
			case ',':
				push(&lexs, "seperator", ",")
			case '(':
				if top(lexs).fsm == lex_fsm["expr"] {
					push(&lexs, "call", "(")
				} else {
					push(&lexs, "bracket", "(")
				}
			case ')':
				push(&lexs, "end_bracket", ")")
			case '{', '}':
				push(&lexs, "set", string(line[i]))
			case '[', ']':
				push(&lexs, "addr", string(line[i]))
			default:
				if is_expr_letters(line[i]) {
					temp := string(line[i])
					i++
					for is_expr_letters(line[i]) {
						temp += string(line[i])
						i++
					}
					i--
					push(&lexs, "expr", temp)
				}
				if is_number(line[i]) {
					temp := string(line[i])
					i++
					for is_number(line[i]) {
						temp += string(line[i])
						i++
					}
					i--
					push(&lexs, "number", temp)
				}
			}
			i++
		}
		if lexs.Len() != 0 {
			lex_lines = append(lex_lines, lexs)
		}
	}
	return
}
