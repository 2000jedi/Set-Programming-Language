#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Thu Jan 18 21:11:54 2017

@author: jedi
"""

import os
import lex
import syn
import execute

def main():
    # filename = sys.argv[1]
    filename = 'debug.sp'
    f = open(os.path.join(os.path.curdir, filename))
    inp = f.readlines()
    lexical = lex.parse(inp)
    #for i in lexical:
    #    print(i, '\n')
    nodes = syn.parse(lexical)
    #for i in nodes:
    #    print(i, '\n')
    execute.execute(nodes)
    
if __name__ == '__main__':
    main()
