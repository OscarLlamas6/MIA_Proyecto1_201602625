package funciones

import "github.com/doun/terminal/color"

var (
	sesionActiva, sesionRoot bool = false, false
)

//EjecutarLogout termina la sesión en caso que si haya una sesión activa
func EjecutarLogout() {

	if sesionActiva {
		sesionActiva = false
		sesionRoot = false
	} else {
		color.Println("@{!r}	No hay ninguna sesión activa actualmente.")
	}

}
