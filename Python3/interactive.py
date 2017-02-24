# -*- coding: utf-8 -*-
"""
Created on Sat Feb 18 21:56:19 2017

@author: 2000jedi
"""

from lib import lex, syn, execute
import traceback

def main():
    # filename = sys.argv[1]
    while 1:
        line = input('>>> ')
        try:
            if line == 'exit()':
                break
            lexical = lex.parse([line])
            nodes = syn.parse(lexical)
            ret = execute.evaluate(nodes[0])
            if ret is not None:
                if type(ret) == lex.lex_class:
                    ret = execute.get_var(ret.data)
                print(ret)
        except Exception as e:
            traceback.print_exc()
    
if __name__ == '__main__':
    main()
