#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Fri Jan 20 07:00:46 2017

@author: jedi
"""

import lex, math

class Number:
    def __init__(self, lexi, denominator=None):
        if denominator is None:
            self.numerator = float(lexi.data)
            self.denominator = 1
            self.reduce()
        else:
            self.numerator = lexi
            self.denominator = denominator
            self.reduce()
            
    def reduce(self):
        if self.numerator == 0:
            if self.denominator != 0:
                self.denominator = 1
            return self
        while self.numerator % 1 != 0.0:
            self.denominator *= 10
            self.numerator *= 10
        while self.denominator % 1 != 0.0:
            self.denominator *= 10
            self.numerator *= 10
        gcd = math.gcd(int(self.denominator), int(self.numerator))
        self.denominator /= gcd
        self.numerator /= gcd
        
        if self.denominator < 0:
            self.numerator = -self.numerator
            self.denominator = -self.denominator
        return self
    
    def __add__(self, num):
        gcd = math.gcd(int(self.denominator), int(num.denominator))
        return Number(self.numerator * gcd / self.denominator + num.numerator * gcd / num.denominator, gcd).reduce()
    
    def __sub__(self, num):
        gcd = math.gcd(int(self.denominator), int(num.denominator))
        return Number(self.numerator * gcd / self.denominator - num.numerator * gcd / num.denominator, gcd).reduce()
    
    def __mul__(self, num):
        return Number(self.numerator * num.numerator, self.denominator * num.denominator).reduce()
        
    def __truediv__(self, num):
        return Number(self.numerator * num.denominator, self.denominator * num.numerator).reduce()
    
    def __float__(self):
        return self.numerator * 1.0 / self.denominator
    
    def __eq__(self, num):
        return self.denominator == num.denominator and self.numerator == num.numerator
    
    def __gt__(self, num):
        return (self - num).numerator > 0
    
    def __ge__(self, num):
        return (self - num).numerator >= 0
    
    def __lt__(self, num):
        return (self - num).numerator < 0

    def __le__(self, num):
        return (self - num).numerator <= 0
        
    def __repr__(self):
        return str(int(self.numerator)) + "/" + str(int(self.denominator)) if self.denominator != 1 else str(self.numerator)

class Set:
    def __init__(self, first_val):
        if type(first_val) != Number:
            first_val = Number(first_val)
        self.data = [first_val]
    
    def append(self, val):
        if type(val) != Number:
            val = Number(val)
        self.data.append(val)
        self.data.sort()

class CSet:
    pass

class Vector:
    def __init__(self, first_val):
        if type(first_val) != Number:
            first_val = Number(first_val)
        self.data = [first_val]
    
    def append(self, val):
        if type(val) != Number:
            val = Number(val)
        self.data.append(val)

class TVector:
    pass

class Function:
    def __init__(self, var, expr):
        self.var = var
        self.expr = expr

class String(str):
    pass

class Boolean(Number):
    def isTrue(self):
        return self.denominator != 0
        
true = Boolean(1, 1)
false = Boolean(0, 1)

class Inherit:
    def __init__(self, argc, func):
        self.func = func
        self.argc = argc
        
def varExchange(lexi):
    if type(lexi) == lex.lex_class:
        if lexi.fsm == lex.lex_fsm.NUMBER:
            return Number(lexi)
        elif lexi.fsm == lex.lex_fsm.EXPR:
            return lexi.data
        else:
            raise TypeError(lexi)
    else:
        return lexi

def operate(op, a, b):
    if op == '+':
        return a + b
    if op == '-':
        return a - b
    if op == '*':
        return a * b
    if op == '/':
        return a / b
    if op == '<':
        return true if a < b else false
    if op == '<=':
        return true if a <= b else false
    if op == '>':
        return true if a > b else false
    if op == '>=':
        return true if a >= b else false
    if op == '==':
        return true if a == b else false
    if op == '!=':
        return true if a != b else false
    if op == '\|\|':
        return true if (a == true) or (b == true) else false
    if op == '&&':
        return true if (a == true) and (b == true) else false
    if op == '~':
        return Function(a, b)
        
def c_for(a, b):
    for i in a.data:
        b.func(i)

def c_range(a, b, c=1):
    pass

def c_import(a):
    pass