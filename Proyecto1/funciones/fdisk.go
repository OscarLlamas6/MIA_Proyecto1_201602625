package funciones

import (
	"Proyecto1/estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

//EjecutarFDisk function
func EjecutarFDisk(size string, unit string, path string, tipo string, fit string, delete string, name string, add string) {

	valorSize := 0
	valorBytes := 0

	if delete == "" && add == "" { //Fdisk normal (crea particiones)

		if size != "" && path != "" && name != "" {

			if strings.HasSuffix(strings.ToLower(path), ".dsk") {

				if fileExists(path) {

					if i, _ := strconv.Atoi(size); i > 0 {

						valorSize = i

						if strings.ToLower(unit) == "k" || unit == "" {
							valorBytes = 1024
						} else if strings.ToLower(unit) == "b" {
							valorBytes = 1
						} else if strings.ToLower(unit) == "m" {
							valorBytes = 1024 * 1024
						}

						valorReal := valorSize * valorBytes
						CrearParticion(valorReal, path, tipo, fit, name)

					} else {
						fmt.Println("El size debe ser mayor que cero.")
					}

				} else {

					fmt.Println("El disco especificado no existe.")
				}

			} else {

				fmt.Println("La ruta debe especificar un archivo con extension '.dsk'.")
			}

		} else {
			fmt.Println("Faltan parámetros obligatorios en la función FDISK")
		}

	} else if delete == "" && add != "" { // Fdisk para agregar o quitar espacio de una particion

	} else if delete != "" && add == "" { // Fdisk para eliminar una particion

	} else {
		fmt.Println("Los parámetros '-delete' y '-add' no pueden venir en la misma instruccion.")
	}

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//CrearParticion fuction
func CrearParticion(size int, path string, tipo string, fit string, name string) {

	if PuedoAgregarParticion(path) {

		if strings.ToLower(tipo) == "p" || tipo == "" { // Particion primaria
			fmt.Println("Creando particion primaria")
			//CALCULAR SI HAY ESPACIO DISPONIBLE
		} else if strings.ToLower(tipo) == "e" { // Particion extendida

			if !ExisteExtendida(path) {
				fmt.Println("Creando particion extendida")
				//CALCULAR SI HAY ESPACIO DISPONIBLE
			} else {
				fmt.Println("El disco ha alcanzado el limite de particiones extendidas.")
			}

		} else if strings.ToLower(tipo) == "l" { // Particion logica
			if ExisteExtendida(path) {
				fmt.Println("Creando particion logica")
				//CALCULAR SI HAY ESPACIO DISPONIBLE
			} else {
				fmt.Println("No se pudo crear la partición lógica porque el disco no tiene partición extendida..")
			}

		}

	} else {
		fmt.Println("El disco ha alcanzado el limite de particiones.")
	}

}

//PuedoAgregarParticion function
func PuedoAgregarParticion(path string) bool {
	file, err := os.Open(path)
	if err != nil { //validar que no sea nulo.
		panic(err)
	}

	Disco1 := estructuras.MBR{}
	//Obtenemos el tamanio del mbr
	DiskSize := int(unsafe.Sizeof(Disco1))
	file.Seek(0, 0)
	//Lee la cantidad de <size> bytes del archivo
	DiskData := leerBytes(file, DiskSize)
	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(DiskData)

	//Decodificamos y guardamos en la variable Disco1
	err = binary.Read(buffer, binary.BigEndian, &Disco1)
	if err != nil {
		file.Close()
		panic(err)
	}
	for i := 0; i < 4; i++ {
		if Disco1.Mpartitions[i].Psize == 0 {
			file.Close()
			return true
		}
	}

	file.Close()
	return false
}

//ExisteExtendida function
func ExisteExtendida(path string) bool {
	file, err := os.Open(path)
	if err != nil { //validar que no sea nulo.
		panic(err)
	}

	Disco1 := estructuras.MBR{}
	//Obtenemos el tamanio del mbr
	DiskSize := int(unsafe.Sizeof(Disco1))
	file.Seek(0, 0)
	//Lee la cantidad de <size> bytes del archivo
	DiskData := leerBytes(file, DiskSize)
	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(DiskData)

	//Decodificamos y guardamos en la variable Disco1
	err = binary.Read(buffer, binary.BigEndian, &Disco1)
	if err != nil {
		file.Close()
		panic(err)
	}
	for i := 0; i < 4; i++ {
		if Disco1.Mpartitions[i].Ptype == 'e' || Disco1.Mpartitions[i].Ptype == 'E' {
			file.Close()
			return true
		}
	}

	file.Close()
	return false
}
