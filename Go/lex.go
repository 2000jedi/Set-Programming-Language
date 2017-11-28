package main

import "fmt"

func processEscapes(s string) string {
	if s == "\\n" {
		return "\n"
	} else {
		return "\\"
	}
}

func isExprLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '_'
}

func isExprLetterLatter(c byte) bool {
	return ('0' <= c && c <= '9') || isExprLetter(c)
}

func isNum(c byte) bool {
	return ('0' <= c && c <= '9') || c == '.'
}

func notOpr(c int) bool {
	return c != VAR_NUMBER && c != VAR_SET && c != VAR_VAR
	//return c != VAR_NUMBER && c != VAR_SET && c != VAR_ADDR && c != VAR_EXPR
}

func lexParse(lines []string) (lexLines []Lexs) {
	if !*debugFlag {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
	}
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
						temp += processEscapes(string(line[i]) + string(line[i+1]))
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
						temp += processEscapes(string(line[i]) + string(line[i+1]))
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
				if isNum(line[i+1]) && notOpr(lexs.Top().fsm) {
					temp := string(line[i])
					i++
					for isNum(line[i]) {
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
				if isExprLetter(line[i]) {
					temp := string(line[i])
					i++
					for isExprLetter(line[i]) {
						temp += string(line[i])
						i++
					}
					i--
					lexs.Push(lexical{LEX_EXPR, temp})
				} else if line[i] == '.' && isExprLetter(line[i+1]) {
					lexs.Push(lexical{LEX_NAMESPACE, "."})
				} else if isNum(line[i]) {
					temp := string(line[i])
					i++
					for isNum(line[i]) {
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
			lexLines = append(lexLines, lexs)
		}
	}
	return
}
