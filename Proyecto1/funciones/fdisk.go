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
//size es el valor que debe tener la particion
func CrearParticion(size int, path string, tipo string, fit string, name string) {

	if PuedoAgregarParticion(path) {

		if strings.ToLower(tipo) == "p" || tipo == "" { // Particion primaria

			HayEspacio, Start := EspacioDisponible(size, path)

			if HayEspacio {

				fmt.Printf("La particion iniciara en el byte: %d\n", Start)

				Pindice := IndiceParticion(path)
				CrearPrimariaOExtendida(Pindice, Start, size, path, fit, name, tipo)

			} else {

				fmt.Println("Operación fallida. No hay espacio disponible para nueva particion.")

			}

		} else if strings.ToLower(tipo) == "e" { // Particion extendida

			if !ExisteExtendida(path) {

				HayEspacio, Start := EspacioDisponible(size, path)

				if HayEspacio {

					fmt.Printf("La particion iniciara en el byte: %d\n", Start)

					Pindice := IndiceParticion(path)
					CrearPrimariaOExtendida(Pindice, Start, size, path, fit, name, tipo)

				} else {

					fmt.Println("Operación fallida. No hay espacio disponible para nueva particion.")

				}

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

//EspacioDisponible function, en caso de revolver TRUE, el valor entero es el byte de inicio para la nueva particion
func EspacioDisponible(size int, path string) (bool, int) {

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
	file.Seek(0, 0)

	for i := DiskSize; i <= int(Disco1.Msize)-size; i++ {
		vacio := true
		for x := 0; x < 4; x++ {

			if Disco1.Mpartitions[x].Psize > 0 {
				if i >= int(Disco1.Mpartitions[x].Pstart) && i <= int(Disco1.Mpartitions[x].Pstart)+int(Disco1.Mpartitions[x].Psize)-1 {
					vacio = false
				} else if i+size-1 >= int(Disco1.Mpartitions[x].Pstart) && i+size-1 <= int(Disco1.Mpartitions[x].Pstart)+int(Disco1.Mpartitions[x].Psize)-1 {
					vacio = false
				} else if i <= int(Disco1.Mpartitions[x].Pstart) && i+size-1 >= int(Disco1.Mpartitions[x].Pstart)+int(Disco1.Mpartitions[x].Psize)-1 {
					vacio = false
				} else if i == int(Disco1.Mpartitions[x].Pstart)+int(Disco1.Mpartitions[x].Psize)-1 {
					vacio = false
				} else if i+size-1 == int(Disco1.Mpartitions[x].Pstart) {
					vacio = false
				}
			}
		}

		if vacio {
			file.Close()
			return true, i
		}

	}

	file.Close()
	return false, 0
}

//PuedoAgregarParticion function
//siguiendo la teoria de particiones, verifica si hay 3 o menos particiones en el disco
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
//verifica si ya existe o no una particion extendida
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

//IndiceParticion function, busca una posicion para almacenar la info de la nueva particion
// en el arreglo de structs (particiones) del MBR.
func IndiceParticion(path string) int {
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
			return i
		}
	}
	file.Close()
	return 0

}

//CrearPrimariaOExtendida function
//indiceMBR es la posicion en el arreglo de structs del MBR, start es el parametro Pstart de c/ particion
func CrearPrimariaOExtendida(indiceMBR int, start int, size int, path string, fit string, name string, tipo string) {

	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		file.Close()
	}

	//Disco1 sera el apuntador al struct MBR temporal
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
	//Seteando los atributos de la nueva particion en un struct del arreglo del MBR
	Disco1.Mpartitions[indiceMBR].Pstart = uint32(start)
	Disco1.Mpartitions[indiceMBR].Psize = uint32(size)
	var chars [16]byte
	copy(chars[:], name)
	copy(Disco1.Mpartitions[indiceMBR].Pname[:], chars[:])
	Disco1.Mpartitions[indiceMBR].Pstatus = 'D' // D = desactivada , A = activada

	if strings.ToLower(tipo) == "p" || tipo == "" {
		Disco1.Mpartitions[indiceMBR].Ptype = 'P'
	} else if strings.ToLower(tipo) == "e" {
		Disco1.Mpartitions[indiceMBR].Ptype = 'E'
	}

	if strings.ToLower(fit) == "wf" || fit == "" {
		Disco1.Mpartitions[indiceMBR].Pfit = 'W'
	} else if strings.ToLower(fit) == "bf" {
		Disco1.Mpartitions[indiceMBR].Pfit = 'B'
	} else if strings.ToLower(fit) == "ff" {
		Disco1.Mpartitions[indiceMBR].Pfit = 'F'
	}

	//Re-escribiendo MBR en el archivo binario (disco)
	file.Seek(0, 0)
	m1 := &Disco1
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, m1)
	escribirBytes(file, binario.Bytes())

	if strings.ToLower(tipo) == "e" {
		//CREAR Y ALMACENAR EBR
		e := estructuras.EBR{}
		e.Enext = -1
		file.Seek(int64(start), 0)
		ebr1 := &e
		var binario1 bytes.Buffer
		binary.Write(&binario1, binary.BigEndian, ebr1)
		escribirBytes(file, binario1.Bytes())
	}

	file.Close()
}

//EspacioDisponibleExtendida function, en caso de revolver TRUE, el valor entero es el byte de inicio para la nueva particion
func EspacioDisponibleExtendida(size int, path string, ebrStart int) (bool, int) {

	file, err := os.Open(path)
	if err != nil { //validar que no sea nulo.
		panic(err)
	}

	file.Close()
	return false, 0
}
