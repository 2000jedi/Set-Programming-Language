package main

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
	return c != VAR_NUMBER && c != VAR_SET && c != VAR_VAR
	//return c != VAR_NUMBER && c != VAR_SET && c != VAR_ADDR && c != VAR_EXPR
}

func lex_parse(lines []string) (lex_lines []Lexs) {
	for _, line := range lines {
		var lexs Lexs
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
				lexs.Push(lexical{LEX_STR, temp})
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
				lexs.Push(lexical{LEX_STR, temp})
			case '=':
				if line[i+1] == '=' {
					lexs.Push(lexical{LEX_OPR, "=="})
					i++
				} else {
					lexs.Push(lexical{LEX_ASSIGN, "="})
				}
			case '+', '*', '/', '<', '>':
				if line[i+1] == '=' {
					lexs.Push(lexical{LEX_OPR, string(line[i]) + "="})
					i++
				} else {
					lexs.Push(lexical{LEX_OPR, string(line[i])})
				}
			case '-':
				if is_number(line[i+1]) && not_opr(lexs.Top().fsm) {
					temp := string(line[i])
					i++
					for is_number(line[i]) {
						temp += string(line[i])
						i++
					}
					i--
					lexs.Push(lexical{LEX_NUMBER, temp})
				} else {
					if line[i+1] == '=' {
						lexs.Push(lexical{LEX_OPR, string(line[i]) + "="})
						i++
					} else {
						lexs.Push(lexical{LEX_OPR, string(line[i])})
					}
				}
			case '!':
				if line[i+1] == '=' {
					lexs.Push(lexical{LEX_OPR, "!="})
					i++
				} else {
					panic("Illegal literal '!'\n")
				}
			case '|':
				if line[i+1] == '|' {
					lexs.Push(lexical{LEX_OPR, "||"})
					i++
				} else {
					lexs.Push(lexical{LEX_OPR, "|"})
				}
			case '&':
				if line[i+1] == '&' {
					lexs.Push(lexical{LEX_OPR, "&&"})
					i++
				} else {
					panic("Illegal literal '&'\n")
				}
			case ':':
				lexs.Push(lexical{LEX_FUNC, ":"})
			case ',':
				lexs.Push(lexical{LEX_SEPERATOR, ","})
			case '(':
				if lexs.Top().fsm == LEX_EXPR {
					lexs.Push(lexical{LEX_CALL, "("})
				} else {
					lexs.Push(lexical{LEX_BRACKET, "("})
				}
			case ')':
				lexs.Push(lexical{LEX_END_BRACKET, ")"})
			case '{', '}':
				lexs.Push(lexical{LEX_BRACES, string(line[i])})
			case '[', ']':
				lexs.Push(lexical{LEX_ADDR, string(line[i])})
			default:
				if is_expr_letters(line[i]) {
					temp := string(line[i])
					i++
					for is_expr_letters(line[i]) {
						temp += string(line[i])
						i++
					}
					i--
					lexs.Push(lexical{LEX_EXPR, temp})
				} else if line[i] == '.' && is_expr_letters(line[i+1]) {
					lexs.Push(lexical{LEX_NAMESPACE, "."})
				} else if is_number(line[i]) {
					temp := string(line[i])
					i++
					for is_number(line[i]) {
						temp += string(line[i])
						i++
					}
					i--
					lexs.Push(lexical{LEX_NUMBER, temp})
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
