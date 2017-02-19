# -*- coding: utf-8 -*-
"""
Created on Sun Feb 19 09:50:52 2017

@author: 2000jedi
"""

from lib import operators as op
import sys

class Inherit:
    def __init__(self, func):
        self.func = func

def println(argc):
    print(*argc)
    
def printf(argc):
    print(end="", *argc)
        
def custom_for(argc):
    if len(argc) > 2:
        print("Augumentation Error: "+ repr(argc), file=sys.stderr)
    a = argc[0]
    b = argc[1]
    return_set = [b.func([i]) for i in a.data]
    return_set = [i for i in return_set if i is not None]
    if return_set == []:
        return None
    return op.Set(return_set)

def custom_range(arg):
    if len(arg) == 2:
        return op.Set(list(range(arg[0].toInt(), arg[1].toInt())))
    elif len(arg) == 3:
        return op.Set(list(range(arg[0].toInt(), arg[1].toInt(), arg[2].toInt())))
    else:
        print("Augumentation Error: "+ len(arg), file=sys.stderr)

def custom_import(a):
    pass