#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Thu Jan 19 07:12:34 2017

@author: jedi
"""

import sys

class lex_fsm:
    NUMBER = 0
    ASSIGN = 1 # =
    EXPR = 2
    OPR = 3
    CALL = 4 # ()
    BRACKET = 5 # ()
    END_BRACKET = 4.5
    SET = 6 # {}
    VEC = 7 # <>
    ADDR = 8 # []
    END_BRACE = 7.5
    SEPERATOR = 9 # ,
    STR = 10
    NAMESPACE = 11
    
    expr_letters = [chr(i) for i in range(65, 90)] + [chr(i) for i in range(97, 122)]
    numbers = [str(i) for i in range(10)] + ['.']

class lex_class:
    def __init__(self, fsm, data):
        self.fsm = fsm
        self.data = data
        
    def __repr__(self):
        return str((self.fsm, self.data))

def process_escapes(raw):
    return raw.encode().decode('unicode_escape')

def parse(raw):
    lexs = []

    for line in raw:
        i = 0
        lexs.append([])
        while i < len(line):
            if line[i] in ['\n', ' ']:
                pass
            elif line[i] == '#':
                i = len(line)
            elif line[i] == '"':
                temp = ''
                i+=1
                while line[i]!='"':
                    if line[i] == '\\':
                        temp += process_escapes(line[i] + line[i+1])
                        i+=2
                    else:
                        temp += line[i]
                        i+=1
                i+=1
                lexs[-1].append(lex_class(lex_fsm.STR, temp))
            elif line[i] == '=':
                if line[i+1] == '=':
                    lexs[-1].append(lex_class(lex_fsm.OPR, '=='))
                    i+=1
                else:
                    lexs[-1].append(lex_class(lex_fsm.ASSIGN, '='))
            elif line[i] == '+':
                if line[i+1] == '=':
                    lexs[-1].append(lex_class(lex_fsm.OPR, '+='))
                    i+=1
                else:
                    lexs[-1].append(lex_class(lex_fsm.OPR, '+'))
            elif line[i] == '-':
                if line[i+1] == '=':
                    lexs[-1].append(lex_class(lex_fsm.OPR, '-='))
                    i+=1
                else:
                    lexs[-1].append(lex_class(lex_fsm.OPR, '-'))
            elif line[i] == '*':
                if line[i+1] == '=':
                    lexs[-1].append(lex_class(lex_fsm.OPR, '*='))
                    i+=1
                else:
                    lexs[-1].append(lex_class(lex_fsm.OPR, '*'))
            elif line[i] == '/':
                if line[i+1] == '=':
                    lexs[-1].append(lex_class(lex_fsm.OPR, '/='))
                    i+=1
                else:
                    lexs[-1].append(lex_class(lex_fsm.OPR, '/'))
            elif line[i] == '<':
                if line[i+1] == '=':
                    lexs[-1].append(lex_class(lex_fsm.OPR, '<='))
                    i+=1
                else:
                    lexs[-1].append(lex_class(lex_fsm.OPR, '<'))
            elif line[i] == '>':
                if line[i+1] == '=':
                    lexs[-1].append(lex_class(lex_fsm.OPR, '>='))
                    i+=1
                else:
                    lexs[-1].append(lex_class(lex_fsm.OPR, '>'))
            elif line[i] == '!':
                if line[i+1] == '=':
                    lexs[-1].append(lex_class(lex_fsm.OPR, '!='))
                    i+=1
                else:
                    sys.stderr.write('Cannot identify operator: !\n')
            elif line[i] == '\|':
                if line[i+1] == '\|':
                    lexs[-1].append(lex_class(lex_fsm.OPR, '\|\|'))
                    i+=1
                else:
                    lexs[-1].append(lex_class(lex_fsm.OPR, '\|'))
            elif line[i] == '&':
                if line[i+1] == '&':
                    lexs[-1].append(lex_class(lex_fsm.OPR, '&&'))
                    i+=1
                else:
                    sys.stderr.write('Cannot identify operator: &\n')
            elif line[i] == '~':
                lexs[-1].append(lex_class(lex_fsm.OPR, '~'))
            elif line[i] == ':':
                lexs[-1].append(lex_class(lex_fsm.NAMESPACE, ':'))
            elif line[i] == ',':
                lexs[-1].append(lex_class(lex_fsm.SEPERATOR, ','))
            elif line[i] == '(':
                if lexs[-1][-1].fsm == lex_fsm.EXPR:
                    lexs[-1].append(lex_class(lex_fsm.CALL, '('))
                else:
                    lexs[-1].append(lex_class(lex_fsm.BRACKET, '('))
            elif line[i] == ')':
                lexs[-1].append(lex_class(lex_fsm.END_BRACKET, ')'))
            elif line[i] in ('{','}'):
                lexs[-1].append(lex_class(lex_fsm.SET, line[i]))
            elif line[i] == '[':
                lexs[-1].append(lex_class(lex_fsm.ADDR, '['))
            elif line[i] == ']':
                lexs[-1].append(lex_class(lex_fsm.END_BRACE, ']'))
            elif line[i] in ('<', '>'): # conflict!
                lexs[-1].append(lex_class(lex_fsm.VEC, line[i]))
            elif line[i] in lex_fsm.expr_letters:
                temp = line[i]
                i+=1
                while line[i] in lex_fsm.expr_letters:
                    temp += line[i]
                    i+=1
                i-=1
                lexs[-1].append(lex_class(lex_fsm.EXPR, temp))
            elif line[i] in lex_fsm.numbers:
                temp = line[i]
                i+=1
                while line[i] in lex_fsm.numbers:
                    temp += line[i]
                    i+=1
                i-=1
                lexs[-1].append(lex_class(lex_fsm.NUMBER, temp))
            i+=1
    return lexs
