package funciones

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//EjecutarMkDisk function
func EjecutarMkDisk(size string, path string, name string, unit string) {

	valorSize := 0
	valorBytes := 0

	if size != "" && path != "" && name != "" {

		if strings.HasSuffix(strings.ToLower(name), ".dsk") {

			if i, err := strconv.Atoi(size); i > 0 {

				valorSize = i
				if err := ensureDir(path); err == nil {

					fmt.Println("Ejecutando MKDISK")

					if strings.ToLower(unit) == "m" || unit == "" || strings.ToLower(unit) == "k" {

						fullName := path + name

						if strings.ToLower(unit) == "m" || unit == "" {
							valorBytes = 1024 * 1024
						} else {
							valorBytes = 1024
						}

						valorReal := valorSize * valorBytes

						file, err := os.Create(fullName) //Crea un nuevo archivo
						if err != nil {
							panic(err)
						}

						// Change permissions Linux.
						err = os.Chmod(fullName, 0777)
						if err != nil {
							log.Println(err)
						}

						data := make([]byte, valorReal) //-size=2 -unit=K
						file.Write(data)                //Escribir datos como un arreglo de bytes
						file.Close()

						/////TOCA ESCRIBIR EL MBR EN EL NUEVO DISCO

					} else {
						fmt.Println("Parámetro 'unit' inválido")
					}

				} else {
					fmt.Println("Directorio inválido")
					panic(err)
				}

			} else {
				fmt.Println("El size debe ser mayor que cero.")
				panic(err)
			}

		} else {
			fmt.Println("El nombre debe contener la extension '.dsk'")
		}

	} else {
		fmt.Println("Faltan parámetros obligatorio en la función MKDISK")
	}
}

func ensureDir(dirName string) error {

	err := os.Mkdir(dirName, 0777)

	if err == nil || os.IsExist(err) {
		return nil
	}
	return err

}
