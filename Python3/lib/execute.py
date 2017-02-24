#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Thu Jan 19 19:50:15 2017

@author: jedi
"""

from lib import lex
from lib import operators as op
from lib import inherit as ih
import sys


class var:
    def __init__(self, val):
        self.stack = [val]

    def top(self):
        return self.stack[-1]

    def pop(self):
        self.stack.pop()

    def push(self, val):
        self.stack.append(val)


def add_var(name, val):
    if name in vars.keys():
        vars[name].push(val)
    else:
        vars[name] = var(val)


def get_var(name):
    return vars[name].top()


def set_var(name, val):
    if name in vars.keys():
        if len(vars[name].stack) != 0:
            vars[name].stack[-1] = val
        else:
            vars[name].stack.append(val)
    else:
        vars[name] = var(val)


def del_var(name):
    vars[name].pop()


vars = {}
add_var("print", ih.Inherit(lambda x: ih.printf(x)))
add_var("println", ih.Inherit(ih.println))
add_var("input", ih.Inherit(input))
add_var("raw_input", ih.Inherit(sys.stdin.readline))
add_var("for", ih.Inherit(ih.custom_for))
add_var("range", ih.Inherit(ih.custom_range))
add_var("import", ih.Inherit(ih.custom_import))
add_var("true", op.Number(1, 1))
add_var("false", op.Number(0, 0))


def evaluate(line):
    stack = []
    ret = None
    temp = None
    i = 0
    while i < len(line):
        if line[i].fsm == lex.lex_fsm.NUMBER:
            stack.append(op.Number(line[i]))
        elif line[i].fsm == lex.lex_fsm.STR:
            stack.append(line[i].data)
        elif line[i].fsm == lex.lex_fsm.EXPR:
            stack.append(line[i])
        elif line[i].fsm == lex.lex_fsm.OPR:
            temp = op.varExchange(stack[-1])
            stack.pop()
            ret = op.varExchange(stack[-1])
            stack.pop()
            if type(temp) != op.Number:
                try:
                    temp = get_var(temp)
                except:
                    sys.stderr.write('Variable Undefined ' + temp + '\n')
            if type(ret) != op.Number:
                try:
                    ret = get_var(ret)
                except:
                    sys.stderr.write('Variable Undefined ' + temp + '\n')
            ret = op.operate(line[i].data, ret, temp)
            stack.append(ret)
        elif line[i].fsm == lex.lex_fsm.SEPERATOR:
            stack.append(line[i])
        elif line[i].fsm == lex.lex_fsm.SET and line[i].data == '{':
            is_function = 0
            segment = []
            i += 1
            while line[i].data != '}':
                segment.append(line[i])
                if line[i].fsm == lex.lex_fsm.FUNC:
                    is_function = len(segment)
                i += 1
            if is_function:
                var = segment[:is_function - 1:2]
                var = [p.data for p in var]
                exprs = segment[is_function:]
                stack.append(op.Function(var, exprs))
            else:
                segment.append(lex.lex_class(lex.lex_fsm.SEPERATOR, ','))
                ret = op.Set()
                while len(segment) != 0:
                    j = 0
                    while segment[j].data != ',':
                        j += 1
                    split = segment[:j]
                    segment = segment[j + 1:]
                    ret.append(evaluate(split))
                stack.append(ret)
        elif line[i].fsm == lex.lex_fsm.CALL:
            if line[i].data == '(':
                stack.append(line[i])
            else:
                if not type(stack[-1]) in (op.Number, op.Set, str) and stack[-1].data == '(':
                    argc = []
                    stack.pop()
                else:
                    argc = [stack[-1]]
                    stack.pop()
                    temp = stack[-2:]
                    while len(temp) > 1 and temp[1].fsm != lex.lex_fsm.CALL:
                        stack = stack[:-2]
                        argc.append(temp[0])
                        temp = stack[-2:]
                    stack.pop()
                    for j in range(len(argc)):
                        if not type(argc[j]) in (op.Number, op.Set, str):
                            if argc[j].fsm == lex.lex_fsm.EXPR:
                                argc[j] = get_var(argc[j].data)
                func = get_var(stack[-1].data)
                stack.pop()
                stack.append(func.func(argc[::-1]))
        elif line[i].fsm == lex.lex_fsm.ASSIGN:
            temp = op.varExchange(stack[-1])
            stack.pop()
            ret = op.varExchange(stack[-1])
            stack.pop()
            set_var(ret, temp)
            # stack.append(temp)
        i += 1
    if len(stack) != 0:
        return stack[-1]


def execute(syn):
    for line in syn:
        evaluate(line)