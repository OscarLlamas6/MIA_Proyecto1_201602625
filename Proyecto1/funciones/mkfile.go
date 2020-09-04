package funciones

import (
	"Proyecto1/estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"

	"github.com/doun/terminal/color"
)

//EjecutarMkfile inicia la creación de un nuevo grupo
func EjecutarMkfile(id string, path string, size string, cont string, p string) {

	if sesionActiva {

		if path != "" && id != "" {

			if IDYaRegistrado(id) {

				NameAux, PathAux := GetDatosPart(id)

				if Existe, Indice := ExisteParticion(PathAux, NameAux); Existe {

					fileMBR, err2 := os.OpenFile(PathAux, os.O_RDWR, 0666)
					if err2 != nil {
						fmt.Println(err2)
						fileMBR.Close()
					}

					// Change permissions Linux.
					err2 = os.Chmod(PathAux, 0666)
					if err2 != nil {
						log.Println(err2)
					}

					//Leemos el MBR
					Disco1 := estructuras.MBR{}
					DiskSize := int(unsafe.Sizeof(Disco1))
					DiskData := leerBytes(fileMBR, DiskSize)
					buffer := bytes.NewBuffer(DiskData)
					err := binary.Read(buffer, binary.BigEndian, &Disco1)
					if err != nil {
						fileMBR.Close()
						fmt.Println(err)
						return
					}

					//LEER EL SUPERBLOQUE
					InicioParticion := Disco1.Mpartitions[Indice].Pstart
					fileMBR.Seek(int64(InicioParticion+1), 0)
					SB1 := estructuras.Superblock{}
					SBsize := int(unsafe.Sizeof(SB1))
					SBData := leerBytes(fileMBR, SBsize)
					buffer2 := bytes.NewBuffer(SBData)
					err = binary.Read(buffer2, binary.BigEndian, &SB1)
					if err != nil {
						fileMBR.Close()
						fmt.Println(err)
						return
					}

					if SB1.MontajesCount > 0 {

						//NOS POSICIONAMOS DONDE EMPIEZA EL STRCUT DE LA CARPETA ROOT (primer struct AVD)
						ApuntadorAVD := SB1.InicioAVDS
						//Vamos a comparar Padres e Hijos
						carpetas := strings.Split(path, "/")
						i := 1
						PathCorrecto := true
						for i < len(carpetas)-1 {

							if TieneSub, ApuntadorSiguiente := ExisteSub(carpetas[i], int(ApuntadorAVD), PathAux); TieneSub {

								//Si entramos a esta parte, significa que el padre si contiene al hijo (subdirectorio)
								//El hijo sería otro padre en el path o directamente será el padre de la carpeta que queremos crear
								//Por lo tanto leeremos otro AVD con el resultado de "APuntadorSiguiente" y seguiremos.
								ApuntadorAVD = int32(ApuntadorSiguiente)
								i++
								PathCorrecto = true

							} else {

								if p != "" {
									//CREAR DIRECTORIO
									//Si entramos a esta parte significa que el directorio requerido no existe Y que en el comando MKDIR
									//se especificó el parámetro de recursividad, es decir debemos crear el directorio (el padre)

									if SB1.FreeAVDS > 0 && SB1.FreeDDS > 0 {

										CrearDirectorio(fileMBR, &SB1, int(ApuntadorAVD), carpetas[i])
										SB1.FirstFreeAVD = SB1.InicioAVDS + (int32(GetBitmap(fileMBR, int(SB1.InicioBitmapAVDS), int(SB1.TotalAVDS))))
										SB1.FirstFreeDD = SB1.InicioDDS + (int32(GetBitmap(fileMBR, int(SB1.InicioBitMapDDS), int(SB1.TotalDDS))))
										SB1.FreeAVDS = SB1.FreeAVDS - int32(1)
										SB1.FreeDDS = SB1.FreeDDS - int32(1)
										fileMBR.Seek(int64(InicioParticion+1), 0)
										//Reescribiendo el Superbloque
										sb1 := &SB1
										var binario1 bytes.Buffer
										binary.Write(&binario1, binary.BigEndian, sb1)
										escribirBytes(fileMBR, binario1.Bytes())
										//Reescribir el Backup del Superbloque
										fileMBR.Seek(int64(SB1.InicioBitacora+(SB1.SizeBitacora*SB1.TotalBitacoras)+1), 0)
										sb2 := &SB1
										var binario2 bytes.Buffer
										binary.Write(&binario2, binary.BigEndian, sb2)
										escribirBytes(fileMBR, binario2.Bytes())
										color.Printf("@{!m}La carpeta @{!y}%v @{!m}fue creada con éxito\n", carpetas[i])

									} else {
										PathCorrecto = false
										color.Println("@{!r} Ya no hay espacio en el sistema de archivos para crear directorio. Acción fallida.")
										break
									}

								} else {
									PathCorrecto = false
									break
								}
							}
						}

						if PathCorrecto {
							//Si se llega a este punto es porque si existian los padres, o se crearon correctamente y podemos
							//escribir el archivo
							//En caso que todos los padres ya existieran
							//Primero verificamos si ya existe el archivo para no repetir nombres
							//En este punto APuntadorAVD apuntará al directorio padre del archivo
							if YaExiste := ExisteFile(carpetas[len(carpetas)-1], int(ApuntadorAVD), PathAux); !YaExiste {

								if SB1.FreeInodos > 0 && SB1.FreeBloques > 0 {

									CrearFile(fileMBR, &SB1, int(ApuntadorAVD), carpetas[len(carpetas)-1])

									//Setear los nuevos parametros del SuperBloque

									fileMBR.Seek(int64(InicioParticion+1), 0)
									//Reescribiendo el Superbloque
									sb1 := &SB1
									var binario1 bytes.Buffer
									binary.Write(&binario1, binary.BigEndian, sb1)
									escribirBytes(fileMBR, binario1.Bytes())
									//Reescribir el Backup del Superbloque
									fileMBR.Seek(int64(SB1.InicioBitacora+(SB1.SizeBitacora*SB1.TotalBitacoras)+1), 0)
									sb2 := &SB1
									var binario2 bytes.Buffer
									binary.Write(&binario2, binary.BigEndian, sb2)
									escribirBytes(fileMBR, binario2.Bytes())

									color.Printf("@{!m}La carpeta @{!y}%v @{!m}fue creada con éxito\n", carpetas[len(carpetas)-1])

								} else {
									color.Println("@{!r} Ya no hay espacio en el sistema de archivos para crear el archivo. Acción fallida.")
								}

							} else {
								color.Println("@{!r} El archivo ya existe")
							}
						} else {
							color.Println("@{!r} Error, una o más carpetas padre no existen.")
						}

					} else {
						color.Println("@{!r} La partición indicada no ha sido formateada.")
					}

					fileMBR.Close()

				} else if ExisteL, IndiceL := ExisteParticionLogica(PathAux, NameAux); ExisteL {

					fileMBR, err := os.Open(PathAux)
					if err != nil { //validar que no sea nulo.
						panic(err)
					}

					EBRAux := estructuras.EBR{}
					EBRSize := int(unsafe.Sizeof(EBRAux))

					//LEER EL SUPERBLOQUE
					InicioParticion := IndiceL + EBRSize
					fileMBR.Seek(int64(InicioParticion+1), 0)
					SB1 := estructuras.Superblock{}
					SBsize := int(unsafe.Sizeof(SB1))
					SBData := leerBytes(fileMBR, SBsize)
					buffer2 := bytes.NewBuffer(SBData)
					err = binary.Read(buffer2, binary.BigEndian, &SB1)
					if err != nil {
						fileMBR.Close()
						fmt.Println(err)
						return
					}

					if SB1.MontajesCount > 0 {

					} else {
						color.Println("@{!r} La partición indicada no ha sido formateada.")
					}

					fileMBR.Close()

				}

			} else {
				color.Printf("@{!r}No hay ninguna partición montada con el id: @{!y}%v\n", id)
			}

		} else {
			color.Println("@{!r}Faltan parámetros obligatorios para la funcion MKDIR.")
		}
	} else {
		color.Println("@{!r}Se necesita de una sesión activa para ejecutar la función MKDIR.")
	}

}

//ExisteFile verifica si el archivo existe o no
func ExisteFile(nombre string, inicioAVD int, path string) bool {

	//LEER AVD
	file, err := os.Open(path)
	if err != nil { //validar que no sea nulo.
		panic(err)
	}

	//CREAMOS UN STRUCT AVD TEMPORAL
	AVDAux := estructuras.AVD{}
	SizeAVD := int(unsafe.Sizeof(AVDAux))
	file.Seek(int64(inicioAVD+1), 0)
	AnteriorData := leerBytes(file, int(SizeAVD))
	buffer2 := bytes.NewBuffer(AnteriorData)
	err = binary.Read(buffer2, binary.BigEndian, &AVDAux)
	if err != nil {
		file.Close()
		fmt.Println(err)
		return false
	}

	//AHORA DEBEMOS LEER EL DETALLE DIRECTORIO DE DICHO AVD
	DDAux := estructuras.DD{}
	PosicionDD := AVDAux.ApuntadorDD
	SizeDD := int(unsafe.Sizeof(DDAux))
	file.Seek(int64(PosicionDD+1), 0)
	DDData := leerBytes(file, int(SizeDD))
	bufferDD := bytes.NewBuffer(DDData)
	err = binary.Read(bufferDD, binary.BigEndian, &DDAux)
	if err != nil {
		file.Close()
		fmt.Println(err)
		return false
	}

	Continuar := true
	//Recorremos el struct DD, y si el apuntador indirecto a apunta a otro DD tambien lo recorremos
	//en caso que no se encuentre el archivo
	for Continuar {
		//Iteramos en las 5 posiciones del arreglo de archivos que tiene el DD
		for i := 0; i < 5; i++ {
			//Validamos que el apuntador al inodo si esté apuntando a algo
			if DDAux.DDFiles[i].ApuntadorInodo > 0 {
				//Comparamos el nombre del archivo con el nombre del archivo que queremos verificar si existe
				//si existe el archivo retornamos true
				var chars [20]byte
				copy(chars[:], nombre)

				if string(DDAux.DDFiles[i].Name[:]) == string(chars[:]) {
					file.Close()
					return true
				}

			}

		}
		//Si el archivo no está en el arreglo de archivos
		//verificamos si el DD actual apunta hacia otro DD

		if DDAux.ApuntadorDD > 0 {

			//Leemos el DD (que se considera contiguo)
			file.Seek(int64(DDAux.ApuntadorDD+int32(1)), 0)
			DDData = leerBytes(file, int(SizeDD))
			bufferDD = bytes.NewBuffer(DDData)
			err = binary.Read(bufferDD, binary.BigEndian, &DDAux)
			if err != nil {
				file.Close()
				fmt.Println(err)
				return false
			}

		} else {
			//Si ya no apunta a otro DD y llegamos a esta parte, cancelamos el ciclo FOR
			Continuar = false
		}
	}

	//De llegar a esta parte significa que el archivo NO EXISTE en el directorio
	file.Close()
	return false
}

//CrearFile crea un archivo en el directorio especificado (AVDPadre)
func CrearFile(file *os.File, sb *estructuras.Superblock, AVDPadre int, nombre string) {

}
