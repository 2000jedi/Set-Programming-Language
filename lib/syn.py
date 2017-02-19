#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Thu Jan 19 07:17:32 2017

@author: jedi
"""

# import sys
from lib import lex

priority = {
    '=': 0,
    '{': 0,
    '&&': 1,
    '\|\|': 1,
    '<':2,
    '>':2,
    '<=': 2,
    '>=': 2,
    '!=': 2,
    '==': 2,
    '+': 3,
    '-': 3,
    '*': 4,
    '/': 4,
    '(': 0,
    '[': 0
}

def segment(line):
    stack = []
    prog = []
    i = 0
    while i < len(line):
        if line[i].fsm in (lex.lex_fsm.EXPR, lex.lex_fsm.NUMBER, lex.lex_fsm.STR):
            prog.append(line[i])
        elif line[i].fsm == lex.lex_fsm.FUNC:
            prog.append(line[i])
        elif line[i].fsm in (lex.lex_fsm.BRACKET, lex.lex_fsm.CALL, lex.lex_fsm.ADDR):
            stack.append(line[i])
            if stack[-1].fsm == lex.lex_fsm.CALL:
                prog.append(stack[-1])
        elif line[i].fsm == lex.lex_fsm.END_BRACKET:
            while not stack[-1].fsm in (lex.lex_fsm.BRACKET, lex.lex_fsm.CALL):
                prog.append(stack[-1])
                stack.pop()
            if stack[-1].fsm == lex.lex_fsm.CALL:
                prog.append(lex.lex_class(lex.lex_fsm.CALL, ")"))
            stack.pop()
        elif line[i].fsm == lex.lex_fsm.END_BRACE:
            while stack[-1].fsm != lex.lex_fsm.ADDR:
                prog.append(stack[-1])
                stack.pop()
            prog.append(stack[-1])
            stack.pop()
        elif line[i].fsm in (lex.lex_fsm.VEC, lex.lex_fsm.SET):
            if line[i].data == '{':
                prog.append(line[i])
                stack.append(line[i])
            elif line[i].data == '}':
                while stack[-1].data != '{':
                    prog.append(stack[-1])
                    stack.pop()
                stack.pop()
                prog.append(line[i])
        elif line[i].fsm == lex.lex_fsm.SEPERATOR:
            prog.append(line[i])
            temp_line = []
            i += 1
            para = 0
            while (not line[i].fsm in (lex.lex_fsm.SEPERATOR, lex.lex_fsm.VEC, lex.lex_fsm.SET)) and (line[i].fsm != lex.lex_fsm.END_BRACKET or para != 0):
                temp_line.append(line[i])
                if line[i].data == '(':
                    para += 1
                if line[i].data == ')':
                    para -= 1
                i += 1
            i -= 1
            prog += segment(temp_line)
        else:
            if len(stack) == 0 or priority[stack[-1].data] < priority[line[i].data]:
                stack.append(line[i])
            else:
                while not len(stack) == 0 and priority[stack[-1].data] >= priority[line[i].data]:
                    prog.append(stack[-1])
                    stack.pop()
                stack.append(line[i])
        i += 1
    while not len(stack) == 0:
        prog.append(stack[-1])
        stack.pop()
    return prog

def parse(lexs):
    prog = []
    for line in lexs:
        prog.append(segment(line))
    return prog
