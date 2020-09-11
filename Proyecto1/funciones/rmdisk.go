package funciones

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/doun/terminal/color"
)

//EjecutarRmDisk function
func EjecutarRmDisk(path string) {

	if path != "" {
		if strings.HasSuffix(strings.ToLower(path), ".dsk") {

			if fileExists(path) {

				fileName := filepath.Base(path)

				color.Printf("@{!yM}¿Está segur@@ que desea borrar el disco %v?̣[Y/n]\n", fileName)

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
						color.Println("@{!r}Error al borrar disco.")
						fmt.Println(err)
					}
					color.Println("@{!c}Disco borrado con éxito.")
				}
			} else {
				color.Println("@{!r}El disco especificado no existe.")
			}

		} else {
			color.Println("@{!r}La ruta debe especificar un archivo con extension '.dsk'.")
		}

	} else {
		color.Println("@{!r}La ruta no puede ser una cadena vacia.")
	}
}
