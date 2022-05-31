//
// lexerSintetico.go
// El programa analizará de forma concurrente una sintaxis de archivos Python
// Marcelo Eduardo Guillen Castillo A00831137
//

package main

import (
	"fmt"
	"os"
)

func main() {
	// Lectura de archivos
	for i, value := range os.Args {
		if i > 0 {
			fmt.Println(i, " --> ", value)
		} else if len(os.Args) == 1 {
			fmt.Println("Seleccione los archivos por favor")
		}
	}

}
