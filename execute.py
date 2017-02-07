#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Thu Jan 19 19:50:15 2017

@author: jedi
"""
import lex
import operators as op
import sys

class namespace:
    def __init__(self, name):
        self.name = name
        self.vars = {}

namespaces = {"": namespace("")}
namespaces[""].vars["print"] = op.Inherit(1, sys.stdout.write)
namespaces[""].vars["println"] = op.Inherit(1, print)
namespaces[""].vars["input"] = op.Inherit(0, input)
namespaces[""].vars["raw_input"] = op.Inherit(0, sys.stdin.readline)
namespaces[""].vars["for"] =  op.Inherit(2, op.c_for)
namespaces[""].vars["range"] = op.Inherit(2, op.c_range)
namespaces[""].vars["import"] = op.Inherit(1, op.c_import)

def evaluate(line):
    stack = []
    ret = None
    temp = None
    i = 0
    while i < len(line):
        if line[i].fsm in (lex.lex_fsm.NUMBER, lex.lex_fsm.EXPR, lex.lex_fsm.STR):
            stack.append(line[i])
        elif line[i].fsm == lex.lex_fsm.OPR:
            temp = op.varExchange(stack[-1])
            stack.pop()
            ret = op.varExchange(stack[-1])
            stack.pop()
            if type(temp) != op.Number:
                try:
                    temp = namespaces[""].vars[temp]
                except:
                    sys.stderr.write('Variable Undefined ' + temp + '\n')
            if type(ret) != op.Number:
                try:
                    ret = namespaces[""].vars[ret]
                except:
                    sys.stderr.write('Variable Undefined ' + temp + '\n')
            op.operate(line[i].data, ret, temp)
            stack.append(ret)
        elif line[i].fsm == lex.lex_fsm.SEPERATOR:
            stack.append(line[i])
        elif line[i].fsm == lex.lex_fsm.SET and line[i].data == '}':
            s = op.Set(stack[-1])
            stack.pop()
            temp = stack[-2:]
            while len(temp) > 1 and temp[1].fsm == lex.lex_fsm.SEPERATOR:
                stack = stack[:-2]
                s.append(temp[0])
                temp = stack[-2:]
            stack.append(s)
        elif line[i].fsm == lex.lex_fsm.VEC and line[i].data == '>':
            s = op.Vector(stack[-1])
            stack.pop()
            temp = stack[-2:]
            while len(temp) > 1 and temp[1].fsm == lex.lex_fsm.SEPERATOR:
                stack = stack[:-2]
                s.append(temp[0])
                temp = stack[-2:]
            stack.append(s)
        elif line[i].fsm == lex.lex_fsm.CALL:
            if line[i].data == '(':
                stack.append(line[i])
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
                    try:
                        if argc[j].fsm == lex.lex_fsm.NUMBER:
                            argc[j] = op.Number(argc[j])
                        elif argc[j].fsm == lex.lex_fsm.STR:
                            argc[j] = argc[j].data
                    except:
                        pass
                func = namespaces[""].vars[stack[-1].data]
                stack.pop()
                if type(func) == op.Inherit:
                    if func.argc == 1:
                        func.func(argc[0])
                    if func.argc == 2:
                        func.func(argc[1], argc[0])
                    if func.argc == 3:
                        func.func(argc[2], argc[1], argc[0])
                    
        elif line[i].fsm == lex.lex_fsm.ASSIGN:
            # print(temp, ret)
            temp = op.varExchange(stack[-1])
            stack.pop()
            ret = op.varExchange(stack[-1])
            stack.pop()
            namespaces[""].vars[ret] = temp
            stack.append(temp)
        i+=1

def execute(syn):
    #print(syn)
    for line in syn:
        evaluate(line)
