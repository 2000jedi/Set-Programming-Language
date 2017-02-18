# -*- coding: utf-8 -*-
"""
Created on Sat Feb 18 21:56:19 2017

@author: 2000jedi
"""

import lex
import syn
import execute

def main():
    # filename = sys.argv[1]
    while 1:
        line = input('>>>')
        lexical = lex.parse([line])
        nodes = syn.parse(lexical)
        execute.execute(nodes)
    
if __name__ == '__main__':
    main()
