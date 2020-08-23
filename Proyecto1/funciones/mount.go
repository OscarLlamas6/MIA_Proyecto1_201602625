package funciones

import (
	"Proyecto1/estructuras"
	"fmt"
)

var (
	//Discos donde se almacenaran los discos que tienen al menos una partición
	Discos []*estructuras.MD
	//PMList lista de todas las particiones
	PMList []*estructuras.PM
)

const abc = "abcdefghijklmnopqrstuvwxyz"

//EjecutarMount function
func EjecutarMount(path string, name string) {

	if fileExists(path) {

		existe, _ := ExisteParticion(path, name)
		existel, _ := ExisteParticionLogica(path, name)

		if existe || existel {

			if !EsExtendida(path, name) {

				if DiscoRegistrado, i := DiscoYaRegistrado(path); DiscoRegistrado {

					if ParticionRegistrada := ParticionYaRegistrada(path, name); !ParticionRegistrada {
						Discos[i].MDcount++
						id := "vd"
						id += getABC(i + 1)
						num := fmt.Sprint(Discos[i].MDcount)
						id += num
						Discos[i].Particiones = append(Discos[i].Particiones, id)

						newPM := new(estructuras.PM)
						newPM.PMid = id
						newPM.PMname = name
						newPM.PMpath = path
						PMList = append(PMList, newPM)
						fmt.Println(id)
					} else {
						fmt.Println("Esta particion ya ha sido montada")
					}

				} else {

					newReg := new(estructuras.MD)
					newReg.MDcount = 1
					newReg.MDocupado = 1
					newReg.MDpath = path
					Discos = append(Discos, newReg)
					id := "vd"
					id += getABC(len(Discos))
					id += "1"

					Discos[len(Discos)-1].Particiones = append(Discos[len(Discos)-1].Particiones, id)

					fmt.Println(id)

					newPM := new(estructuras.PM)
					newPM.PMid = id
					newPM.PMname = name
					newPM.PMpath = path
					PMList = append(PMList, newPM)
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

	if len(Discos) > 0 {
		for i := 0; i < len(Discos); i++ {
			if Discos[i].MDpath == path {
				return true, i
			}
		}
	}
	return false, 0
}

//ParticionYaRegistrada verifica si la partición ya ha sido montada con aterioridad
func ParticionYaRegistrada(path string, name string) bool {

	if len(PMList) > 0 {
		for i := 0; i < len(PMList); i++ {
			if PMList[i].PMpath == path && PMList[i].PMname == name {
				return true
			}
		}
	}
	return false
}

func getABC(i int) string {
	return abc[i-1 : i]
}
