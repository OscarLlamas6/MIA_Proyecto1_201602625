package main

import (
	"Proyecto1/analizadores"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//Pausa fuc
func Pausa() {
	fmt.Print("Ejecución pausada. Presiona 'Enter' para continuar...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

//LimpiarPantalla fuction
func LimpiarPantalla() {

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	fmt.Println("------------ SISTEMA DE ARCHIVOS LWH | Dev. By Oscar Llamas ------------")
}

func main() {

	continuar := true
	LimpiarPantalla()

	for continuar {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">>")
		input, _ := reader.ReadString('\n')

		if runtime.GOOS == "windows" {
			input = strings.TrimRight(input, "\r\n")
		} else {
			input = strings.TrimRight(input, "\n")
		}

		if strings.HasSuffix(input, "\\*") {
			pedir := true
			linea := ""
			for pedir {
				linea, _ = reader.ReadString('\n')
				if runtime.GOOS == "windows" {
					linea = strings.TrimRight(linea, "\r\n")
				} else {
					linea = strings.TrimRight(linea, "\n")
				}

				input += " "
				input += linea
				if !strings.HasSuffix(input, "\\*") {
					pedir = false
				}
			}
		}

		if strings.ToLower(input) == "pause" {
			Pausa()
		} else if strings.ToLower(input) == "exit" {
			continuar = false
		} else if strings.ToLower(input) == "clear" {
			LimpiarPantalla()
		} else {
			analizadores.Lexico(input)
		}
	}

	//comando := "mdisk create path->10M"
	/*filename := "hola.txt"
	filebuffer, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	analizadores.Lexico(string(filebuffer), &errorLexico)
	*/
	/*if !errorLexico {
		fmt.Println("Analisis lexico exitoso")
	} else {
		fmt.Println("Analisis lexico no exitoso")
	}*/

}
