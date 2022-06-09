#
# lexerSintetico.py
# un programa que colorea la sintaxis de Python dentro de un html más css
# Marcelo Eduardo Guillen Castillo A00831137

import sys
path_entrada = sys.argv[1]
path_salida = sys.argv[2]
def lexer(path_entrada):
    # usaremos <span></span>
    archivo = open(path_entrada, 'r')
    size = len(archivo.readlines())
    archivo.close()
    
    archivo = open(path_entrada, 'r')
    argument = ''
    
    token = ['identificador','clave','operadores','delimitadores','literales', 'comentario', 'otro'] # 7 elementos
    elegido = token[-1]
    tipo_literal = ['string','numero','m-string','ninguno']
    literal_elegido = tipo_literal[-1]
    
    palabras_clave = ['and', 'del', 'for', 'is' ,'raise',
    'assert' ,'elif' ,'from' ,'lambda' ,'return',
    'break' ,'else' ,'global' ,'not' ,'try',
    'class' ,'except' ,'if' ,'or' ,'while',
    'continue' ,'exec' ,'import' ,'pass' ,'with',
    'def' ,'finally' ,'in' ,'print' ,'yield']
    
    alphabet_list = ['A','B','C','D','E','F','G','H','I','J','K','M','N','L','O','P','Q','R','S','T','U','V','W','X','Y','Z',
                     'a','b','c','d','e','f','g','h','i','j','k','m','n','l','o','p','q','r','s','t','u','v','w','x','y','z']
    
    lista_operadores_simples = ['&', '|', '^', '~', '+', '*', '/', '%', '<', '>', '!', '=','-']
    lista_operadores_compuestos = ['+=','-=','**','//','<>','<=','>=','==','!=','<<','>>']
    
    contLineas = 0
    script_temporal = ''
    
    unPunto = True
    unCientifico = True
    unNegativo = True
    
    for line in archivo:
        contLineas+=1
        if(len(line) >= 2):
            if(line[-1] == '\n'):
                line = line[:-1] + ' '
            else:
                line += ' '
        else:
            line += ' '
            
        line = line.replace('\t', '    ')
        
        contChar = 0
        for i in line:
            contChar+=1
            # logica python
            
            # comentario
            if(elegido == 'comentario'):
                if(contChar == len(line)): # termina
                    elegido = token[-1]
                    argument += '<span style="color: #00FF7F">'+script_temporal+'</span>'
                    script_temporal = ''
                else:
                    if(i == ' '):
                        script_temporal += '&nbsp;'
                    else:
                        script_temporal += i
            elif(i == '#' and elegido == 'otro'):
                elegido = token[5]
                script_temporal += i
            
            # delimitadores
            if(elegido == 'delimitadores'):
                elegido = token[-1]
                argument += '<span style="color: #FFA500">' + script_temporal + '</span>'
                script_temporal = ''
                if(i == '#' and elegido == 'otro'):
                    elegido = token[5]
                    script_temporal += i
                elif((i == '(' or i == '{' or i == '[' or i == ')' or i == '}' or i == ']') and elegido == 'otro'):
                    elegido = token[3]
                    script_temporal += i
            elif((i == '(' or i == '{' or i == '[' or i == ')' or i == '}' or i == ']') and elegido == 'otro'):
                elegido = token[3]
                script_temporal += i
            
            # literales
            if(elegido == 'literales'):
                if(literal_elegido == 'string'): # seccion strings
                    if(len(script_temporal) >= 2 and script_temporal[0] == "'" and script_temporal[1] == "'"
                       and i == "'"): # m-string ''
                        elegido = token[4]
                        literal_elegido = tipo_literal[2]
                        script_temporal += i
                    elif(len(script_temporal) >= 2 and script_temporal[0] == '"' and script_temporal[1] == '"'
                       and i == '"'): # m-string ""
                        elegido = token[4]
                        literal_elegido = tipo_literal[2]
                        script_temporal += i
                    elif(len(script_temporal) >= 2 and script_temporal[0] == "'" and script_temporal[-1] == "'"):
                        elegido = token[-1]
                        literal_elegido = tipo_literal[-1]
                        argument += '<span style="color: magenta; background-color: #2E8B57">'+script_temporal+'</span>'
                        script_temporal = ''
                        if((i == '(' or i == '{' or i == '[' or i == ')' or i == '}' or i == ']') and elegido == 'otro'):
                            elegido = token[3]
                            script_temporal += i
                        elif(i == '#' and elegido == 'otro'):
                            elegido = token[5]
                            script_temporal += i
                        elif(i.isnumeric() and elegido == 'otro'): # es un numero
                            elegido = token[4]
                            literal_elegido = tipo_literal[1]
                            script_temporal += i
                    elif(len(script_temporal) >= 2 and script_temporal[0] == '"' and script_temporal[-1] == '"'):
                        elegido = token[-1]
                        literal_elegido = tipo_literal[-1]
                        argument += '<span style="color: magenta; background-color: #2E8B57">'+script_temporal+'</span>'
                        script_temporal = ''
                        if((i == '(' or i == '{' or i == '[' or i == ')' or i == '}' or i == ']') and elegido == 'otro'):
                            elegido = token[3]
                            script_temporal += i
                        elif(i == '#' and elegido == 'otro'):
                            elegido = token[5]
                            script_temporal += i
                        elif(i.isnumeric() and elegido == 'otro'): # es un numero
                            elegido = token[4]
                            literal_elegido = tipo_literal[1]
                            script_temporal += i
                    elif(contChar == len(line)):
                        elegido = token[-1]
                        literal_elegido = tipo_literal[-1]
                        argument += '<span style="color: magenta; background-color: #2E8B57">'+script_temporal+'</span>'
                        script_temporal = ''
                    else:
                        if(i == ' '):
                            script_temporal += '&nbsp;'
                        else:
                            script_temporal += i
                elif(literal_elegido == 'numero'): # seccion numeros
                    if(i == '.' and unPunto): # un punto decimal
                        unPunto = False
                        unNegativo = False
                        script_temporal += i
                    elif(i == '-' and not(unPunto) and unNegativo and (script_temporal[-1] == 'e' or script_temporal[-1] == 'E')): # un negativo
                        unNegativo = False
                        script_temporal += i
                    elif((i == 'e' or i == 'E') and unCientifico): # una notacion cientifica
                        unPunto = False
                        unCientifico = False
                        unNegativo = True
                        script_temporal += i
                    elif(not(i.isnumeric())): # la iteracion ya no forma parte de un numero
                        elegido = token[-1]
                        literal_elegido = tipo_literal[-1]
                        argument += '<span style="color: #9370DB">'+script_temporal+'</span>'
                        script_temporal = ''
                        unPunto = True
                        unNegativo = True
                        unCientifico = True
                        if((i == '(' or i == '{' or i == '[' or i == ')' or i == '}' or i == ']') and elegido == 'otro'):
                            elegido = token[3]
                            script_temporal += i
                        elif((i == '"' or i == "'") and elegido == 'otro'): # es un string
                            elegido = token[4]
                            literal_elegido = tipo_literal[0]
                            script_temporal += i
                    else:
                        script_temporal += i
                elif(literal_elegido == 'm-string'): # super string
                    if(len(script_temporal) >= 6 and (script_temporal[:2] == '""' and script_temporal[2] == '"')
                       and (script_temporal[-3:-1] == '""' and script_temporal[-1] == '"')): # termina con triple "
                        elegido = token[-1]
                        literal_elegido = tipo_literal[-1]
                        argument += '<span style="color: magenta; background-color: #2E8B57">'+script_temporal+'</span>'
                        script_temporal = ''
                        if((i == '"' or i == "'") and elegido == 'otro'): # es un string
                            elegido = token[4]
                            literal_elegido = tipo_literal[0]
                            script_temporal += i
                        elif(i.isnumeric() and elegido == 'otro'): # es un numero
                            elegido = token[4]
                            literal_elegido = tipo_literal[1]
                            script_temporal += i
                    elif(len(script_temporal) >= 6 and (script_temporal[:2] == "''" and script_temporal[2] == "'")
                         and (script_temporal[-3:-1] == "''" and script_temporal[-1] == "'")): # termina con triple "
                        elegido = token[-1]
                        literal_elegido = tipo_literal[-1]
                        argument += '<span style="color: magenta; background-color: #2E8B57">'+script_temporal+'</span>'
                        script_temporal = ''
                        if((i == '"' or i == "'") and elegido == 'otro'): # es un string
                            elegido = token[4]
                            literal_elegido = tipo_literal[0]
                            script_temporal += i
                        elif(i.isnumeric() and elegido == 'otro'): # es un numero
                            elegido = token[4]
                            literal_elegido = tipo_literal[1]
                            script_temporal += i
                    elif(contChar == len(line) and contLineas == size): # se acaba el documento
                        elegido = token[-1]
                        literal_elegido = tipo_literal[-1]
                        argument += '<span style="color: magenta; background-color: #2E8B57">'+script_temporal+'</span>'
                        script_temporal = ''
                    else:
                        if(i == ' '): # espacios
                            script_temporal += '&nbsp;'
                        else:
                            script_temporal += i
                        if(contChar == len(line)): # si ya viene la siguiente linea
                            script_temporal += '<br>'
            elif(i.isnumeric() and elegido == 'otro'): # es un numero
                elegido = token[4]
                literal_elegido = tipo_literal[1]
                script_temporal += i
            elif((i == '"' or i == "'") and elegido == 'otro'): # es un string
                elegido = token[4]
                literal_elegido = tipo_literal[0]
                script_temporal += i
            
            # identificadores
            if(elegido == 'identificador'):
                if(not((i in alphabet_list) or i.isnumeric() or i == '_')): # acaba de leer identificadores
                    if(script_temporal in palabras_clave): # palabras clave
                        script_temporal = '<span style="color: crimson">'+script_temporal+'</span>'
                    argument += script_temporal + '</span>'
                    elegido = token[-1]
                    script_temporal = ''
                    if(i.isnumeric() and elegido == 'otro'): # es un numero
                        elegido = token[4]
                        literal_elegido = tipo_literal[1]
                        script_temporal += i
                    elif((i == '"' or i == "'") and elegido == 'otro'): # es un string
                        elegido = token[4]
                        literal_elegido = tipo_literal[0]
                        script_temporal += i
                    elif(i == '#'):
                        elegido = token[5]
                        script_temporal += i
                    elif((i == '(' or i == '{' or i == '[' or i == ')' or i == '}' or i == ']') and elegido == 'otro'):
                        elegido = token[3]
                        script_temporal += i
                else:
                    script_temporal += i
            elif((i in alphabet_list) and elegido == 'otro'): # primer valor para variables
                elegido = token[0]
                argument += '<span style="color: blue;">'
                script_temporal += i
                
            # operadores
            if(elegido == 'operadores'):
                if(script_temporal not in lista_operadores_compuestos and (i not in lista_operadores_simples)
                   and script_temporal not in lista_operadores_simples): # la combinacion no existe
                    argument += '<span>' + script_temporal + '</span>'
                    elegido = token[-1]
                    script_temporal = ''
                    if(i.isalpha() and elegido == 'otro'): # primer valor para variables
                        elegido = token[0]
                        argument += '<span style="color: blue;">'
                        script_temporal += i
                    elif(i.isnumeric() and elegido == 'otro'): # es un numero
                        elegido = token[4]
                        literal_elegido = tipo_literal[1]
                        script_temporal += i
                    elif((i == '"' or i == "'") and elegido == 'otro'): # es un string
                        elegido = token[4]
                        literal_elegido = tipo_literal[0]
                        script_temporal += i
                    elif(i == '#'):
                        elegido = token[5]
                        script_temporal += i
                    elif((i == '(' or i == '{' or i == '[' or i == ')' or i == '}' or i == ']') and elegido == 'otro'):
                        elegido = token[3]
                        script_temporal += i
                elif(script_temporal in lista_operadores_compuestos): # acaba de leer operadores compuestos
                    argument += '<span style="color: green; font-weight: bold;">' + script_temporal + '</span>'
                    elegido = token[-1]
                    script_temporal = ''
                    if(i.isalpha() and elegido == 'otro'): # primer valor para variables
                        elegido = token[0]
                        argument += '<span style="color: blue;">'
                        script_temporal += i
                    elif(i.isnumeric() and elegido == 'otro'): # es un numero
                        elegido = token[4]
                        literal_elegido = tipo_literal[1]
                        script_temporal += i
                    elif((i == '"' or i == "'") and elegido == 'otro'): # es un string
                        elegido = token[4]
                        literal_elegido = tipo_literal[0]
                        script_temporal += i
                    elif(i == '#'):
                        elegido = token[5]
                        script_temporal += i
                    elif((i == '(' or i == '{' or i == '[' or i == ')' or i == '}' or i == ']') and elegido == 'otro'):
                        elegido = token[3]
                        script_temporal += i
                elif(script_temporal in lista_operadores_simples and (i not in lista_operadores_simples)): # acaba de leer operadores simples
                    argument += '<span style="color: green; font-weight: bold;">' + script_temporal + '</span>'
                    elegido = token[-1]
                    script_temporal = ''
                    if(i.isalpha() and elegido == 'otro'): # primer valor para variables
                        elegido = token[0]
                        argument += '<span style="color: blue;">'
                        script_temporal += i
                    elif(i.isnumeric() and elegido == 'otro'): # es un numero
                        elegido = token[4]
                        literal_elegido = tipo_literal[1]
                        script_temporal += i
                    elif((i == '"' or i == "'") and elegido == 'otro'): # es un string
                        elegido = token[4]
                        literal_elegido = tipo_literal[0]
                        script_temporal += i
                    elif(i == '#'):
                        elegido = token[5]
                        script_temporal += i
                    elif((i == '(' or i == '{' or i == '[' or i == ')' or i == '}' or i == ']') and elegido == 'otro'):
                        elegido = token[3]
                        script_temporal += i
                elif(i in lista_operadores_simples): # la iteracion sigue siendo parte de los operadores
                    script_temporal += i
            elif(i in lista_operadores_simples and elegido == 'otro'):
                elegido = token[2]
                script_temporal += i
                
            # otros
            if(elegido == 'otro' and i != ' '):
                argument += i
            if(i == ' ' and elegido == 'otro'):
                argument += '&nbsp;'
            if(len(line) == contChar and literal_elegido != 'm-string'):
                argument += '<br>'
    archivo.close()
    # salida
    archivo = open(path_salida, 'w')
    html = "<!DOCTYPE html><html><body>"+argument+"</body></html>"
    archivo.write(html)
    archivo.close()
    
lexer(path_entrada)

# Segun los cálculos, la complejidad de mi programa es de O(n**2)
# debido a que solo hay dos ciclos 'for' usandose a la vez.