package funciones

import (
	"fmt"
)

var (
	//Discos donde se almacenaran los discos que tienen al menos una partición
	Discos []string
)

const abc = "abcdefghijklmnopqrstuvwxyz"

//EjecutarMount function
func EjecutarMount(path string, name string) {

	if fileExists(path) {

		existe, _ := ExisteParticion(path, name)
		existel, _ := ExisteParticionLogica(path, name)

		if existe || existel {

			if !EsExtendida(path, name) {

				if DiscoRegistrado, _ := DiscoYaRegistrado(path); !DiscoRegistrado {

					fmt.Println("disco registrado")
					fmt.Println(len(Discos))

				} else {
					fmt.Println("disco no registrado")
				}

			} else {
				fmt.Println("No se puede montar porque es una partición extendida.")
			}

		} else {
			fmt.Println("El disco especificado no tiene ninguna partición con ese nombre.")
		}

	} else {
		fmt.Println("El disco especificado no existe.")
	}
}

//DiscoYaRegistrado verifica si ese disco ya tiene alguna otra particion montada, para asignar nueva letra
func DiscoYaRegistrado(path string) (bool, int) {
	Discos = append(Discos, "hola")
	return false, 0
}

//ParticionYaRegistrada verifica si la partición ya ha sido montada con aterioridad
func ParticionYaRegistrada(indice int, name string) bool {

	return false
}

func getABC(i int) string {
	return abc[i-1 : i]
}

//IndiceDisponible busca un nuevo espacio en el arreglo para asignar una letra al disco
func IndiceDisponible() int {

	return -1
}
