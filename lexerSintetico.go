//
// lexerSintetico.go
// El programa analizará de forma concurrente una sintaxis de archivos Python
// Marcelo Eduardo Guillen Castillo A00831137
//

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		fmt.Println("Error al seleccionar el directorio!")
	}
	// Lectura de archivos
	for _, f := range files {
		// Leer los archivos de entrada y almacenarlos en un arreglo de buffers
		lexerSintetico(os.Args[1] + (f.Name()))
	}

}

// Analisis y coloreo de sintaxis
func lexerSintetico(filepath string) {
	// Creacion del buffer del archivo
	fileBuffer, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("No puedo leer el archivo!")
		os.Exit(1)
	}
	// Lectura del archivo
	dato := bufio.NewScanner(strings.NewReader(string(fileBuffer)))
	dato.Split(bufio.ScanLines)
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
	unPunto, numCientifico, unNegativo := true, true, false
	for dato.Scan() {
		linea := dato.Text()
		linea += " "
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
				} else if partOfArray(delimitadores_lista[:], string(value)) {
					cabezal = tokens[5]
					temporal_word += string(value)
				}
			} else if partOfArray(delimitadores_lista[:], string(value)) {
				cabezal = tokens[5]
				temporal_word += string(value)
			}

			// Literales 4
			if cabezal == tokens[4] {
				if cabezal_tipo_lateral == tipo_literal[3] { // es un numero
					if unPunto && string(value) == "." {
						unPunto = !unPunto
						temporal_word += string(value)
					} else if (string(value) == "E" || string(value) == "e") && numCientifico {
						numCientifico = !numCientifico
						temporal_word += string(value)
						unNegativo = !unNegativo
						unPunto = !unPunto
					} else if !isDecimal(string(value)) { // termina de leer numeros
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: #9370DB;'>" + temporal_word + "</span>"
						temporal_word = ""
						if (string(value) == "\"" || string(value) == "'") && cabezal == tokens[0] { // es un string 1
							cabezal = tokens[4]
							temporal_word += string(value)
							cabezal_tipo_lateral = tipo_literal[1]
						} else if partOfArray(delimitadores_lista[:], string(value)) {
							cabezal = tokens[5]
							temporal_word += string(value)
						} else if string(value) == "#" && cabezal == tokens[0] {
							temporal_word += string(value)
							cabezal = tokens[3]
						}
					} else if string(value) == "-" && unNegativo {
						unNegativo = !unNegativo
						temporal_word += string(value)
					} else {
						unNegativo = !unNegativo
						temporal_word += string(value)
					}
				} else if cabezal_tipo_lateral == tipo_literal[1] { // es un simple string
					if (temporal_word[0] == '"' && temporal_word[len(temporal_word)-1] == '"') || i == len(linea)-1 { // termina string con ""
						cabezal = tokens[0]
						cabezal_tipo_lateral = tipo_literal[0]
						argument += "<span style='color: magenta; background-color: #2E8B57'>" + temporal_word + "</span>"
						temporal_word = ""
					}
				} else if cabezal_tipo_lateral == tipo_literal[2] { // es un super string o de triple comillas

				}
			} else if isDecimal(string(value)) && cabezal == tokens[0] { // es un numero 3
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
					if partOfArray(delimitadores_lista[:], string(value)) { // delimitadores
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
					} else if partOfArray(delimitadores_lista[:], string(value)) {
						cabezal = tokens[5]
						temporal_word += string(value)
					} else if string(value) == "#" && cabezal == tokens[0] {
						temporal_word += string(value)
						cabezal = tokens[3]
					}
				} else {
					temporal_word += string(value)
				}
			} else if partOfArray(operadores_lista[:], string(value)) {
				cabezal = tokens[2]
				temporal_word += string(value)
			}

			// Ningun token
			if string(value) == " " {
				argument += "&nbsp;"
			} else if cabezal == tokens[0] && string(value) != " " {
				argument += string(value)
			}
		}
		argument += "<br>\n"
	}
	html = "<!DOCTYPE html><html><title>Python Sintaxis</title></html><body>" + argument + "</body>"
	writeFile(html, filepath)
}

func writeFile(html, name string) {
	// Nombre del html final
	salidaName := name[:len(name)-len(filepath.Ext(name))] // Nombre del archivo salida
	// Escritura del archivo a convertir a html
	file, err := os.Create(string(filepath.Base(salidaName)) + ".html")
	if err != nil {
		fmt.Println("Error al crear el archivo destino!!")
	}
	_, err2 := file.WriteString(html)
	if err2 != nil {
		fmt.Println("Error al escribir sobre el archivo!!")
	}
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
