package main

func push(lexs *Stack, fsm int, data string) {
	lexs.Push(&storage{VAR_FSM, &lexical{fsm, data}})
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
	return ('0' <= c && c <= '9') || is_expr_letters(c)
}

func is_number(c byte) bool {
	return ('0' <= c && c <= '9') || c == '.'
}

func not_opr(c int) bool {
	return c != VAR_NUMBER && c != VAR_SET && c != VAR_EXPR
	//return c != VAR_NUMBER && c != VAR_SET && c != VAR_ADDR && c != VAR_EXPR
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
				push(&lexs, LEX_STR, temp)
			case '\'':
				temp := ""
				i++
				for line[i] != '\'' {
					if line[i] == '\\' {
						temp += process_escapes(string(line[i]) + string(line[i+1]))
						i += 2
					} else {
						temp += string(line[i])
						i++
					}
				}
				push(&lexs, LEX_STR, temp)
			case '=':
				if line[i+1] == '=' {
					push(&lexs, LEX_OPR, "==")
					i++
				} else {
					push(&lexs, LEX_ASSIGN, "=")
				}
			case '+', '*', '/', '<', '>':
				if line[i+1] == '=' {
					push(&lexs, LEX_OPR, string(line[i])+"=")
					i++
				} else {
					push(&lexs, LEX_OPR, string(line[i]))
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
					push(&lexs, LEX_NUMBER, temp)
				} else {
					if line[i+1] == '=' {
						push(&lexs, LEX_OPR, string(line[i])+"=")
						i++
					} else {
						push(&lexs, LEX_OPR, string(line[i]))
					}
				}
			case '!':
				if line[i+1] == '=' {
					push(&lexs, LEX_OPR, "!=")
					i++
				} else {
					panic("Illegal literal '!'\n")
				}
			case '|':
				if line[i+1] == '|' {
					push(&lexs, LEX_OPR, "||")
					i++
				} else {
					push(&lexs, LEX_OPR, "|")
				}
			case '&':
				if line[i+1] == '&' {
					push(&lexs, LEX_OPR, "&&")
					i++
				} else {
					panic("Illegal literal '&'\n")
				}
			case ':':
				push(&lexs, LEX_FUNC, ":")
			case ',':
				push(&lexs, LEX_SEPERATOR, ",")
			case '(':
				if top(lexs).fsm == LEX_EXPR {
					push(&lexs, LEX_CALL, "(")
				} else {
					push(&lexs, LEX_BRACKET, "(")
				}
			case ')':
				push(&lexs, LEX_END_BRACKET, ")")
			case '{', '}':
				push(&lexs, LEX_BRACES, string(line[i]))
			case '[', ']':
				push(&lexs, LEX_ADDR, string(line[i]))
			default:
				if is_expr_letters(line[i]) {
					temp := string(line[i])
					i++
					for is_expr_letters(line[i]) {
						temp += string(line[i])
						i++
					}
					i--
					push(&lexs, LEX_EXPR, temp)
				} else if line[i] == '.' && is_expr_letters(line[i+1]) {
					push(&lexs, LEX_NAMESPACE, ".")
				} else if is_number(line[i]) {
					temp := string(line[i])
					i++
					for is_number(line[i]) {
						temp += string(line[i])
						i++
					}
					i--
					push(&lexs, LEX_NUMBER, temp)
				} else {
					panic(line[i])
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
