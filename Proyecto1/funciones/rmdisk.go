package funciones

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

//EjecutarRmDisk function
func EjecutarRmDisk(path string) {

	if path != "" {
		if strings.HasSuffix(strings.ToLower(path), ".dsk") {

			if fileExists(path) {

				fmt.Println("¿Está segur@ que desea borrar este disco?")

				pedir := true
				linea := ""

				for pedir {
					reader := bufio.NewReader(os.Stdin)
					input, _ := reader.ReadString('\n')

					if runtime.GOOS == "windows" {
						input = strings.TrimRight(input, "\r\n")
					} else {
						input = strings.TrimRight(input, "\n")
					}

					if strings.ToLower(input) == "n" || strings.ToLower(input) == "y" {
						linea = input
						pedir = false
					}

				}

				if strings.ToLower(linea) == "y" {
					err := os.Remove(path)

					if err != nil {
						fmt.Println("Error al borrar disco.")
						fmt.Println(err)
					}
					fmt.Println("Disco borrado con éxito.")
				}

			} else {

				fmt.Println("El disco especificado no existe.")
			}

		} else {

			fmt.Println("La ruta debe especificar un archivo con extension '.dsk'.")
		}

	} else {
		fmt.Println("La ruta no puede ser una cadena vacia.")
	}
}
