//
// lexerSintetico.go
// El programa analizará de forma concurrente una sintaxis de archivos Python
// Marcelo Eduardo Guillen Castillo A00831137
/*
La complejidad de mi algoritmo si fuera secuencial sería de O(n^3), pero
ahora de manera paralela sería O(n^3)/files, donde files representa la cantidad de archivos a leer
*/

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		fmt.Println("Error al seleccionar el directorio!")
	}
	c := make(chan bool, 10)

	// Lectura de archivos secuencial
	start := time.Now()
	for _, f := range files {
		// Leer los archivos de entrada y almacenarlos en un arreglo de buffers
		if string((os.Args[1])[len(os.Args[1])-1]) != "\\" {
			html := lexerSintetico_secuencial(os.Args[1] + "\\" + (f.Name()))
			fileName := string(filepath.Base(f.Name()))
			writeFile(html, os.Args[2]+"\\"+"A00831137_"+fileName)
		} else {
			html := lexerSintetico_secuencial(os.Args[1] + (f.Name()))
			fileName := string(filepath.Base(f.Name()))
			writeFile(html, os.Args[2]+"A00831137_"+fileName)
		}
	}
	durationNew := time.Since(start)
	fmt.Println("Time duration from secuencial algorithm: " + durationNew.String())
	// Lectura de archivos paralela
	size := 0
	startNew := time.Now()
	for _, f := range files {
		// Leer los archivos de entrada y almacenarlos en un arreglo de buffers
		if string((os.Args[1])[len(os.Args[1])-1]) != "\\" {
			go lexerSintetico(os.Args[1]+"\\"+(f.Name()), c)
		} else {
			go lexerSintetico(os.Args[1]+(f.Name()), c)
		}
		size++
	}
	// Escritura de archivos paralela
	for i := 0; i < size; i++ {
		ok := <-c
		if !ok {
			close(c)
			break
		}
	}
	duration := time.Since(startNew)
	fmt.Println("Time duration from paralel algorithm: " + duration.String())
}

// Analisis y coloreo de sintaxis
func lexerSintetico(file string, c chan bool) {
	// Creacion del buffer del archivo
	fileBuffer, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("No puedo leer el archivo!")
		os.Exit(1)
	}
	// Lectura del archivo
	dato := bufio.NewScanner(strings.NewReader(string(fileBuffer)))
	dato.Split(bufio.ScanLines)
	// informacion del archivo
	fileSize := cantidad_lineas(string(fileBuffer))
	contador_lineas := 0
	// Variables a utilizar
	var argument string
	var html string
	var temporal_word string
	//					0			1					2			3				4			5
	tokens := [6]string{"ninguno", "identificadores", "operadores", "comentarios", "laterales", "delimitador"}
	cabezal := tokens[0]
	// Lista de delimitadores
	delimitadores_lista := [6]string{"(", "{", "[", "]", "}", ")"}
	// Lista de palabras clave
	palabras_clave := [30]string{"and", "del", "for", "is", "raise",
		"assert", "elif", "from", "lambda", "return",
		"break", "else", "global", "not", "try",
		"class", "except", "if", "or", "while",
		"continue", "exec", "import", "pass", "with",
		"def", "finally", "in", "print", "yield"}
	// Lista de operadores
	operadores_lista := [13]string{"&", "|", "^", "~", "+", "*", "/", "%", "<", ">", "!", "=", "-"}
	// Literales
	//							0		1			2				3
	tipo_literal := [4]string{"ninguno", "string", "super-string", "entero"}
	cabezal_tipo_lateral := tipo_literal[0]
	unPunto, numCientifico := true, true
	unNegativo := false
	for dato.Scan() {
		linea := dato.Text()
		linea += " "
		contador_lineas++
		linea = strings.ReplaceAll(linea, "\t", "    ")
		for i, value := range linea {
			// Aqui se hace el analisis de sintaxis
			// el texto a leer es con: string(value)

			// Comentarios
			if cabezal == tokens[3] {
				if i == len(dato.Text()) { // Termina la linea, por ende acaba comentario
					argument += "<span style='color: green'>" + temporal_word + "</span>"
					temporal_word = ""
					cabezal = tokens[0]
				} else {
					temporal_word += string(value)
				}
			} else if string(value) == "#" && cabezal == tokens[0] {
				temporal_word += string(value)
				cabezal = tokens[3]
			}

			// Delimitadores
			if cabezal == tokens[5] { // Termina de leer delimitadores
				cabezal = tokens[0]
				argument += "<span style='color: #FFA500;+'>" + temporal_word + "</span>"
				temporal_word = ""
				if string(value) == "#" && cabezal == tokens[0] {
					temporal_word += string(value)
					cabezal = tokens[3]
				} else if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
					cabezal = tokens[5]
					temporal_word += string(value)
				}
			} else if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] { // delimitador
				cabezal = tokens[5]
				temporal_word += string(value)
			}

			// Literales 4
			if cabezal == tokens[4] {
				if cabezal_tipo_lateral == tipo_literal[3] { // es un numero
					if unPunto && string(value) == "." {
						unPunto = false
						temporal_word += string(value)
					} else if string(value) == "-" && unNegativo && !unPunto && (temporal_word[len(temporal_word)-1] == 'e' || temporal_word[len(temporal_word)-1] == 'E') { // negativo para cientifico
						unNegativo = false
						temporal_word += string(value)
					} else if (string(value) == "E" || string(value) == "e") && numCientifico {
						numCientifico = false
						temporal_word += string(value)
						unNegativo = true
						unPunto = false
					} else if !isDecimal(string(value)) { // termina de leer numeros
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: #9370DB;'>" + temporal_word + "</span>"
						temporal_word = ""
						unNegativo = false
						numCientifico = true
						unPunto = true
						if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[1]
						} else if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
							cabezal = tokens[5]
							temporal_word += string(value)
						} else if string(value) == "#" && cabezal == tokens[0] {
							temporal_word += string(value)
							cabezal = tokens[3]
						}
					} else {
						unNegativo = false
						temporal_word += string(value)
					}
				} else if cabezal_tipo_lateral == tipo_literal[1] { // es un simple string
					if len(temporal_word) >= 2 && temporal_word[0:2] == "\"\"" && string(value) == "\"" { // inicio de super-string
						cabezal_tipo_lateral = tipo_literal[2]
						cabezal = tokens[4]
						temporal_word += string(value)
					} else if len(temporal_word) >= 2 && temporal_word[0:2] == "''" && string(value) == "'" {
						cabezal_tipo_lateral = tipo_literal[2]
						cabezal = tokens[4]
						temporal_word += string(value)
					} else if (len(temporal_word) >= 2 && temporal_word[0] == '"' && temporal_word[len(temporal_word)-1] == '"') || i == len(linea)-1 { // termina string con ""
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: magenta; background-color: #2E8B57'>" + temporal_word + "</span>"
						temporal_word = ""
						if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
							cabezal = tokens[5]
							temporal_word += string(value)
						} else if string(value) == "#" && cabezal == tokens[0] {
							temporal_word += string(value)
							cabezal = tokens[3]
						} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[3]
						}
					} else if (len(temporal_word) >= 2 && temporal_word[0] == '\'' && temporal_word[len(temporal_word)-1] == '\'') || i == len(linea)-1 { // termina string con ''
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: magenta; background-color: #2E8B57'>" + temporal_word + "</span>"
						temporal_word = ""
						if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
							cabezal = tokens[5]
							temporal_word += string(value)
						} else if string(value) == "#" && cabezal == tokens[0] {
							temporal_word += string(value)
							cabezal = tokens[3]
						} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[3]
						}
					} else {
						if string(value) == " " {
							temporal_word += "&nbsp;"
						} else if string(value) == "&" {
							temporal_word += "&amp;"
						} else if string(value) == "<" {
							temporal_word += "&lt;"
						} else if string(value) == ">" {
							temporal_word += "&gt;"
						} else {
							temporal_word += string(value)
						}
					}
				} else if cabezal_tipo_lateral == tipo_literal[2] { // es un super string o de triple comillas
					if len(temporal_word) >= 6 && (temporal_word[:3] == "\"\"\"" &&
						temporal_word[(len(temporal_word)-3):] == "\"\"\"") { // termina de leer super-string
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: magenta; background-color: #2E8B57'>" + temporal_word + "</span>"
						temporal_word = ""
						if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
							cabezal = tokens[5]
							temporal_word += string(value)
						} else if string(value) == "#" && cabezal == tokens[0] {
							temporal_word += string(value)
							cabezal = tokens[3]
						} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[3]
						} else if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[1]
						}
					} else if len(temporal_word) >= 6 && (temporal_word[:3] == "'''" &&
						temporal_word[(len(temporal_word)-3):] == "'''") { // termina de leer super-string
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: magenta; background-color: #2E8B57'>" + temporal_word + "</span>"
						temporal_word = ""
						if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
							cabezal = tokens[5]
							temporal_word += string(value)
						} else if string(value) == "#" && cabezal == tokens[0] {
							temporal_word += string(value)
							cabezal = tokens[3]
						} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[3]
						} else if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[1]
						}
					} else if i == len(linea)-1 && contador_lineas == fileSize { // final del documento
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: magenta; background-color: #2E8B57'>" + temporal_word + "</span>"
						temporal_word = ""
					} else {
						if string(value) == " " {
							temporal_word += "&nbsp;"
						} else if string(value) == "&" {
							temporal_word += "&amp;"
						} else if string(value) == "<" {
							temporal_word += "&lt;"
						} else if string(value) == ">" {
							temporal_word += "&gt;"
						} else {
							temporal_word += string(value)
						}
					}
				}
			} else if (isDecimal(string(value))) && cabezal == tokens[0] { // es un numero 3
				cabezal = tokens[4]
				temporal_word += string(value)
				cabezal_tipo_lateral = tipo_literal[3]
			} else if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
				cabezal = tokens[4]
				temporal_word += string(value)
				cabezal_tipo_lateral = tipo_literal[1]
			}

			// identificadores 1
			if cabezal == tokens[1] {
				if !(isAlphabet(string(value)) || isDecimal(string(value)) || string(value) == "_") { // termina de leer identificadores
					if partOfArray(palabras_clave[:], temporal_word) {
						temporal_word = "<span style='color: crimson;'>" + temporal_word + "</span>"
					}
					argument += "<span style='color: blue;'>" + temporal_word + "</span>"
					cabezal = tokens[0]
					temporal_word = ""
					if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] { // delimitadores
						cabezal = tokens[5]
						temporal_word += string(value)
					} else if string(value) == "#" && cabezal == tokens[0] { // comentario
						temporal_word += string(value)
						cabezal = tokens[3]
					} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
						cabezal = tokens[4]
						temporal_word += string(value)
						cabezal_tipo_lateral = tipo_literal[3]
					} else if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
						cabezal = tokens[4]
						temporal_word += string(value)
						cabezal_tipo_lateral = tipo_literal[1]
					}
				} else {
					temporal_word += string(value)
				}
			} else if isAlphabet(string(value)) && cabezal == tokens[0] {
				cabezal = tokens[1]
				temporal_word += string(value)
			}

			// operadores 2
			if cabezal == tokens[2] {
				if !(partOfArray(operadores_lista[:], string(value))) { // termina de leer operadores
					cabezal = tokens[0]
					argument += "<span style='color: violet;'>" + temporal_word + "</span>"
					temporal_word = ""
					if isAlphabet(string(value)) && cabezal == tokens[0] {
						cabezal = tokens[1]
						temporal_word += string(value)
					} else if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
						cabezal = tokens[5]
						temporal_word += string(value)
					} else if string(value) == "#" && cabezal == tokens[0] {
						temporal_word += string(value)
						cabezal = tokens[3]
					} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
						cabezal = tokens[4]
						temporal_word += string(value)
						cabezal_tipo_lateral = tipo_literal[3]
					} else if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
						cabezal = tokens[4]
						temporal_word += string(value)
						cabezal_tipo_lateral = tipo_literal[1]
					}
				} else {
					temporal_word += string(value)
				}
			} else if partOfArray(operadores_lista[:], string(value)) && cabezal == tokens[0] {
				cabezal = tokens[2]
				temporal_word += string(value)
			}

			// Ningun token
			if string(value) != " " && cabezal == tokens[0] {
				argument += string(value)
			} else if string(value) == " " && cabezal == tokens[0] {
				argument += "&nbsp;"
			} else if string(value) == "&" && cabezal == tokens[0] {
				argument += "&amp;"
			} else if string(value) == "<" && cabezal == tokens[0] {
				argument += "&lt;"
			} else if string(value) == ">" && cabezal == tokens[0] {
				argument += "&gt;"
			}
		}
		if cabezal_tipo_lateral != tipo_literal[2] {
			argument += "<br>\n"
		} else {
			temporal_word += "<br>\n"
		}
	}
	html = "<!DOCTYPE html><html><title>Python Sintaxis</title></html><body>" + argument + "</body>"
	//  string(filepath.Base(file))
	if string((os.Args[2])[len(os.Args[2])-1]) != "\\" {
		writeFile_paralel(html, os.Args[2]+"\\"+"A00831137_"+string(filepath.Base(file)), c)
	} else {
		writeFile_paralel(html, os.Args[2]+"A00831137_"+string(filepath.Base(file)), c)
	}
}

func lexerSintetico_secuencial(file string) string {
	// Creacion del buffer del archivo
	fileBuffer, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("No puedo leer el archivo!")
		os.Exit(1)
	}
	// Lectura del archivo
	dato := bufio.NewScanner(strings.NewReader(string(fileBuffer)))
	dato.Split(bufio.ScanLines)
	// informacion del archivo
	fileSize := cantidad_lineas(string(fileBuffer))
	contador_lineas := 0
	// Variables a utilizar
	var argument string
	var html string
	var temporal_word string
	//					0			1					2			3				4			5
	tokens := [6]string{"ninguno", "identificadores", "operadores", "comentarios", "laterales", "delimitador"}
	cabezal := tokens[0]
	// Lista de delimitadores
	delimitadores_lista := [6]string{"(", "{", "[", "]", "}", ")"}
	// Lista de palabras clave
	palabras_clave := [30]string{"and", "del", "for", "is", "raise",
		"assert", "elif", "from", "lambda", "return",
		"break", "else", "global", "not", "try",
		"class", "except", "if", "or", "while",
		"continue", "exec", "import", "pass", "with",
		"def", "finally", "in", "print", "yield"}
	// Lista de operadores
	operadores_lista := [13]string{"&", "|", "^", "~", "+", "*", "/", "%", "<", ">", "!", "=", "-"}
	// Literales
	//							0		1			2				3
	tipo_literal := [4]string{"ninguno", "string", "super-string", "entero"}
	cabezal_tipo_lateral := tipo_literal[0]
	unPunto, numCientifico := true, true
	unNegativo := false
	for dato.Scan() {
		linea := dato.Text()
		linea += " "
		contador_lineas++
		linea = strings.ReplaceAll(linea, "\t", "    ")
		for i, value := range linea {
			// Aqui se hace el analisis de sintaxis
			// el texto a leer es con: string(value)

			// Comentarios
			if cabezal == tokens[3] {
				if i == len(dato.Text()) { // Termina la linea, por ende acaba comentario
					argument += "<span style='color: green'>" + temporal_word + "</span>"
					temporal_word = ""
					cabezal = tokens[0]
				} else {
					temporal_word += string(value)
				}
			} else if string(value) == "#" && cabezal == tokens[0] {
				temporal_word += string(value)
				cabezal = tokens[3]
			}

			// Delimitadores
			if cabezal == tokens[5] { // Termina de leer delimitadores
				cabezal = tokens[0]
				argument += "<span style='color: #FFA500;+'>" + temporal_word + "</span>"
				temporal_word = ""
				if string(value) == "#" && cabezal == tokens[0] {
					temporal_word += string(value)
					cabezal = tokens[3]
				} else if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
					cabezal = tokens[5]
					temporal_word += string(value)
				}
			} else if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] { // delimitador
				cabezal = tokens[5]
				temporal_word += string(value)
			}

			// Literales 4
			if cabezal == tokens[4] {
				if cabezal_tipo_lateral == tipo_literal[3] { // es un numero
					if unPunto && string(value) == "." {
						unPunto = false
						temporal_word += string(value)
					} else if string(value) == "-" && unNegativo && !unPunto && (temporal_word[len(temporal_word)-1] == 'e' || temporal_word[len(temporal_word)-1] == 'E') { // negativo para cientifico
						unNegativo = false
						temporal_word += string(value)
					} else if (string(value) == "E" || string(value) == "e") && numCientifico {
						numCientifico = false
						temporal_word += string(value)
						unNegativo = true
						unPunto = false
					} else if !isDecimal(string(value)) { // termina de leer numeros
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: #9370DB;'>" + temporal_word + "</span>"
						temporal_word = ""
						unNegativo = false
						numCientifico = true
						unPunto = true
						if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[1]
						} else if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
							cabezal = tokens[5]
							temporal_word += string(value)
						} else if string(value) == "#" && cabezal == tokens[0] {
							temporal_word += string(value)
							cabezal = tokens[3]
						}
					} else {
						unNegativo = false
						temporal_word += string(value)
					}
				} else if cabezal_tipo_lateral == tipo_literal[1] { // es un simple string
					if len(temporal_word) >= 2 && temporal_word[0:2] == "\"\"" && string(value) == "\"" { // inicio de super-string
						cabezal_tipo_lateral = tipo_literal[2]
						cabezal = tokens[4]
						temporal_word += string(value)
					} else if len(temporal_word) >= 2 && temporal_word[0:2] == "''" && string(value) == "'" {
						cabezal_tipo_lateral = tipo_literal[2]
						cabezal = tokens[4]
						temporal_word += string(value)
					} else if (len(temporal_word) >= 2 && temporal_word[0] == '"' && temporal_word[len(temporal_word)-1] == '"') || i == len(linea)-1 { // termina string con ""
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: magenta; background-color: #2E8B57'>" + temporal_word + "</span>"
						temporal_word = ""
						if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
							cabezal = tokens[5]
							temporal_word += string(value)
						} else if string(value) == "#" && cabezal == tokens[0] {
							temporal_word += string(value)
							cabezal = tokens[3]
						} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[3]
						}
					} else if (len(temporal_word) >= 2 && temporal_word[0] == '\'' && temporal_word[len(temporal_word)-1] == '\'') || i == len(linea)-1 { // termina string con ''
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: magenta; background-color: #2E8B57'>" + temporal_word + "</span>"
						temporal_word = ""
						if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
							cabezal = tokens[5]
							temporal_word += string(value)
						} else if string(value) == "#" && cabezal == tokens[0] {
							temporal_word += string(value)
							cabezal = tokens[3]
						} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[3]
						}
					} else {
						if string(value) == " " {
							temporal_word += "&nbsp;"
						} else if string(value) == "&" {
							temporal_word += "&amp;"
						} else if string(value) == "<" {
							temporal_word += "&lt;"
						} else if string(value) == ">" {
							temporal_word += "&gt;"
						} else {
							temporal_word += string(value)
						}
					}
				} else if cabezal_tipo_lateral == tipo_literal[2] { // es un super string o de triple comillas
					if len(temporal_word) >= 6 && (temporal_word[:3] == "\"\"\"" &&
						temporal_word[(len(temporal_word)-3):] == "\"\"\"") { // termina de leer super-string
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: magenta; background-color: #2E8B57'>" + temporal_word + "</span>"
						temporal_word = ""
						if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
							cabezal = tokens[5]
							temporal_word += string(value)
						} else if string(value) == "#" && cabezal == tokens[0] {
							temporal_word += string(value)
							cabezal = tokens[3]
						} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[3]
						} else if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[1]
						}
					} else if len(temporal_word) >= 6 && (temporal_word[:3] == "'''" &&
						temporal_word[(len(temporal_word)-3):] == "'''") { // termina de leer super-string
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: magenta; background-color: #2E8B57'>" + temporal_word + "</span>"
						temporal_word = ""
						if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
							cabezal = tokens[5]
							temporal_word += string(value)
						} else if string(value) == "#" && cabezal == tokens[0] {
							temporal_word += string(value)
							cabezal = tokens[3]
						} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[3]
						} else if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[1]
						}
					} else if i == len(linea)-1 && contador_lineas == fileSize { // final del documento
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: magenta; background-color: #2E8B57'>" + temporal_word + "</span>"
						temporal_word = ""
					} else {
						if string(value) == " " {
							temporal_word += "&nbsp;"
						} else if string(value) == "&" {
							temporal_word += "&amp;"
						} else if string(value) == "<" {
							temporal_word += "&lt;"
						} else if string(value) == ">" {
							temporal_word += "&gt;"
						} else {
							temporal_word += string(value)
						}
					}
				}
			} else if (isDecimal(string(value))) && cabezal == tokens[0] { // es un numero 3
				cabezal = tokens[4]
				temporal_word += string(value)
				cabezal_tipo_lateral = tipo_literal[3]
			} else if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
				cabezal = tokens[4]
				temporal_word += string(value)
				cabezal_tipo_lateral = tipo_literal[1]
			}

			// identificadores 1
			if cabezal == tokens[1] {
				if !(isAlphabet(string(value)) || isDecimal(string(value)) || string(value) == "_") { // termina de leer identificadores
					if partOfArray(palabras_clave[:], temporal_word) {
						temporal_word = "<span style='color: crimson;'>" + temporal_word + "</span>"
					}
					argument += "<span style='color: blue;'>" + temporal_word + "</span>"
					cabezal = tokens[0]
					temporal_word = ""
					if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] { // delimitadores
						cabezal = tokens[5]
						temporal_word += string(value)
					} else if string(value) == "#" && cabezal == tokens[0] { // comentario
						temporal_word += string(value)
						cabezal = tokens[3]
					} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
						cabezal = tokens[4]
						temporal_word += string(value)
						cabezal_tipo_lateral = tipo_literal[3]
					} else if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
						cabezal = tokens[4]
						temporal_word += string(value)
						cabezal_tipo_lateral = tipo_literal[1]
					}
				} else {
					temporal_word += string(value)
				}
			} else if isAlphabet(string(value)) && cabezal == tokens[0] {
				cabezal = tokens[1]
				temporal_word += string(value)
			}

			// operadores 2
			if cabezal == tokens[2] {
				if !(partOfArray(operadores_lista[:], string(value))) { // termina de leer operadores
					cabezal = tokens[0]
					argument += "<span style='color: violet;'>" + temporal_word + "</span>"
					temporal_word = ""
					if isAlphabet(string(value)) && cabezal == tokens[0] {
						cabezal = tokens[1]
						temporal_word += string(value)
					} else if partOfArray(delimitadores_lista[:], string(value)) && cabezal == tokens[0] {
						cabezal = tokens[5]
						temporal_word += string(value)
					} else if string(value) == "#" && cabezal == tokens[0] {
						temporal_word += string(value)
						cabezal = tokens[3]
					} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
						cabezal = tokens[4]
						temporal_word += string(value)
						cabezal_tipo_lateral = tipo_literal[3]
					} else if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
						cabezal = tokens[4]
						temporal_word += string(value)
						cabezal_tipo_lateral = tipo_literal[1]
					}
				} else {
					temporal_word += string(value)
				}
			} else if partOfArray(operadores_lista[:], string(value)) && cabezal == tokens[0] {
				cabezal = tokens[2]
				temporal_word += string(value)
			}

			// Ningun token
			if string(value) != " " && cabezal == tokens[0] {
				argument += string(value)
			} else if string(value) == " " && cabezal == tokens[0] {
				argument += "&nbsp;"
			} else if string(value) == "&" && cabezal == tokens[0] {
				argument += "&amp;"
			} else if string(value) == "<" && cabezal == tokens[0] {
				argument += "&lt;"
			} else if string(value) == ">" && cabezal == tokens[0] {
				argument += "&gt;"
			}
		}
		if cabezal_tipo_lateral != tipo_literal[2] {
			argument += "<br>\n"
		} else {
			temporal_word += "<br>\n"
		}
	}
	html = "<!DOCTYPE html><html><title>Python Sintaxis</title></html><body>" + argument + "</body>"
	return html
}

func writeFile(html, name string) {
	// Nombre del html final
	salidaName := name[:len(name)-len(filepath.Ext(name))] // Nombre del archivo salida
	// Escritura del archivo a convertir a html
	file, err := os.Create(string(salidaName) + ".html")
	if err != nil {
		fmt.Println("Error al crear el archivo destino!!")
	}
	_, err2 := file.WriteString(html)
	if err2 != nil {
		fmt.Println("Error al escribir sobre el archivo!!")
	}
}

func writeFile_paralel(html, name string, c chan bool) {
	// Nombre del html final
	salidaName := name[:len(name)-len(filepath.Ext(name))] // Nombre del archivo salida
	// Escritura del archivo a convertir a html
	file, err := os.Create(string(salidaName) + ".html")
	if err != nil {
		fmt.Println("Error al crear el archivo destino!!")
	}
	_, err2 := file.WriteString(html)
	if err2 != nil {
		fmt.Println("Error al escribir sobre el archivo!!")
	}
	c <- true
}

func partOfArray(array []string, ele string) bool {
	for _, value := range array {
		if value == ele {
			return true
		}
	}
	return false
}

func isAlphabet(letter string) bool {
	return (letter >= "a" && letter <= "z") || (letter >= "A" && letter <= "Z")
}

func isDecimal(integer string) bool {
	_, err := strconv.Atoi(integer)
	return err == nil
}

func cantidad_lineas(fileBuffer string) int {
	dato := bufio.NewScanner(strings.NewReader(fileBuffer))
	dato.Split(bufio.ScanLines)
	cont := 0
	for dato.Scan() {
		cont++
	}
	return cont
}
