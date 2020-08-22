package funciones

import "fmt"

//EjecutarUnmount function
func EjecutarUnmount(lista *[]string) {
	for _, id := range *lista {
		fmt.Println(id)
	}
}
