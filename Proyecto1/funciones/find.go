package funciones

import (
	"path/filepath"
	"strings"

	"github.com/doun/terminal/color"
)

//EjecutarFind function
func EjecutarFind(id string, path string, name string) {

	if sesionActiva {

		if path != "" && id != "" && name != "" {

			if strings.HasPrefix(path, "/") {

				if len(name) <= 20 {

					extension := filepath.Ext(name)

					if strings.ToLower(extension) == ".txt" || strings.ToLower(extension) == ".pdf" || strings.ToLower(extension) == ".mia" || strings.ToLower(extension) == ".dsk" || strings.ToLower(extension) == ".sh" {

						FindFile(id, path, name)

					} else {
						FindDir(id, path, name)
					}

				} else {

					color.Println("@{!r} El nombre del archivo o carpeta no puede tener más de 20 caracteres")

				}

			} else {
				color.Println("@{!r}Path incorrecto, debe iniciar con @{!y}/")
			}

		} else {
			color.Println("@{!r}Faltan parámetros obligatorios para la funcion REN.")
		}

	} else {
		color.Println("@{!r}Se necesita de una sesión activa para ejecutar la función MKDIR.")
	}

}

//FindFile busca un archivo
func FindFile(id string, path string, name string) {

}

//FindDir busca una carpeta
func FindDir(id string, path string, name string) {
	if last := len(path) - 1; last >= 0 && path[last] == '/' {
		path = path[:last]
	}

}
