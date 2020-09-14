package funciones

import (
	"strings"

	"github.com/doun/terminal/color"
)

//EjecutarEdit function
func EjecutarEdit(id string, path string, size string, cont string) {

	if sesionActiva {

		if path != "" && id != "" {

			if strings.HasPrefix(path, "/") {

			} else {
				color.Println("@{!r}Path incorrecto, debe iniciar con @{!y}/")
			}

		} else {
			color.Println("@{!r}Faltan parámetros obligatorios para la funcion REN.")
		}

	} else {
		color.Println("@{!r}Se necesita de una sesión activa para ejecutar la función EDIT.")
	}

}
