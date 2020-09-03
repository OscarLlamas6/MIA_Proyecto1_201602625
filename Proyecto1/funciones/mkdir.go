package funciones

import (
	"Proyecto1/estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/doun/terminal/color"
)

//EjecutarMkdir function
func EjecutarMkdir(id string, path string, p string) {

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
						i := 0
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
							//escribir el directorio hijo (DIR objectivo)
							//En caso que todos los padres ya existieran
							//Primero verificamos si ya existe el directorio para no repetir nombres
							//En este punto APuntadorAVD apuntará al padre más cercano
							if YaExiste, _ := ExisteSub(carpetas[len(carpetas)-1], int(ApuntadorAVD), PathAux); !YaExiste {

								if SB1.FreeAVDS > 0 && SB1.FreeDDS > 0 {

									CrearDirectorio(fileMBR, &SB1, int(ApuntadorAVD), carpetas[len(carpetas)-1])
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

								} else {
									color.Println("@{!r} Ya no hay espacio en el sistema de archivos para crear directorio. Acción fallida.")
								}

							} else {
								color.Println("@{!r} El directorio ya existe")
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

//ExisteSub verifica si existe el hijo en el padre, retorna el apuntador al hijo (al subdirectorio)
func ExisteSub(nombre string, inicioAVD int, path string) (bool, int) {

	//LEER AVD
	file, err := os.Open(path)
	if err != nil { //validar que no sea nulo.
		panic(err)
	}

	//CREAMOS UN STRUCT TEMPORAL
	AVDAux := estructuras.AVD{}
	SizeAVD := int(unsafe.Sizeof(AVDAux))
	file.Seek(int64(inicioAVD+1), 0)
	AnteriorData := leerBytes(file, int(SizeAVD))
	buffer2 := bytes.NewBuffer(AnteriorData)
	err = binary.Read(buffer2, binary.BigEndian, &AVDAux)
	if err != nil {
		file.Close()
		fmt.Println(err)
		return false, 0
	}

	Continuar := true
	//Recorremos el struct y si el apuntador indirecto a punta a otro AVD tambien lo recorreremos en caso que no se encuentre
	//el directorio
	for Continuar {
		//Iteramos en las 6 posiciones del arreglo de subdirectoios (apuntadores)
		for i := 0; i < 6; i++ {
			//Validamos que el apuntador si esté apuntando a algo
			if AVDAux.ApuntadorSubs[i] > 0 {
				//Con el valor del apuntador leemos un struct AVD
				AVDHijo := estructuras.AVD{}
				file.Seek(int64(AVDAux.ApuntadorSubs[i]+int32(1)), 0)
				HijoData := leerBytes(file, int(SizeAVD))
				buffer := bytes.NewBuffer(HijoData)
				err = binary.Read(buffer, binary.BigEndian, &AVDHijo)
				if err != nil {
					file.Close()
					fmt.Println(err)
					return false, 0
				}
				//Comparamos el nombre del AVD leido con el nombre del directorio que queremos verificar si existe
				//si existe el directorio retornamos true y el byte donde está dicho AVD
				var chars [20]byte
				copy(chars[:], nombre)

				if string(AVDHijo.NombreDir[:]) == string(chars[:]) {
					file.Close()
					return true, int(AVDAux.ApuntadorSubs[i])
				}
			}

		}
		//Si el directorio no está en el arreglo de apuntadores directos
		//verificamos si el AVD actual apunta hacia otro AVD con otros 6 apuntadores
		if AVDAux.ApuntadorAVD > 0 {
			//Leemos el AVD (que se considera contiguo)
			file.Seek(int64(AVDAux.ApuntadorAVD+int32(1)), 0)
			AnteriorData = leerBytes(file, int(SizeAVD))
			buffer2 := bytes.NewBuffer(AnteriorData)
			err = binary.Read(buffer2, binary.BigEndian, &AVDAux)
			if err != nil {
				file.Close()
				fmt.Println(err)
				return false, 0
			}

		} else {
			//Si ya no apunta a otro AVD y llegamos a esta parte, cancelamos el ciclo FOR
			Continuar = false
		}

	}
	//De llegar a esta parte significa que el subdirectorio NO EXISTE en el directorio
	file.Close()
	return false, 0
}

//CrearDirectorio crea la carpeta
func CrearDirectorio(file *os.File, sb *estructuras.Superblock, AVDPadre int, nombre string) {

	//Buscamos la posicion en el bitmap para el nuevo AVD
	PosicionEnBitmapAVD := GetBitmap(file, int(sb.InicioBitmapAVDS), int(sb.TotalAVDS))
	//Calculamos la posicion del byte en el bitmap AVD
	BitmapPos := int(sb.InicioBitmapAVDS) + PosicionEnBitmapAVD
	//Escribimos un 1 en esa posición del bitmap
	file.Seek(int64(BitmapPos+1), 0)
	data := []byte{0x01}
	file.Write(data)
	//Calculamos la posicion del byte del nuevo AVD
	AVDPos := int(sb.InicioAVDS) + (int(sb.SizeAVD) * (PosicionEnBitmapAVD))
	//Creamos el nuevo AVD
	newAVD := estructuras.AVD{}
	t := time.Now()
	cadena := t.Format("2006-01-02 15:04:05")
	var charsDate [20]byte
	copy(charsDate[:], cadena)
	//Seteando fecha de creacion
	copy(newAVD.FechaCreacion[:], charsDate[:])
	var ArrayNombre [20]byte
	//Seteamos el nombre del nuevo directorio
	nombreDir := nombre
	copy(ArrayNombre[:], nombreDir)
	//Seteando nombre del directorio "/"
	copy(newAVD.NombreDir[:], ArrayNombre[:])
	nombrePropietario := idSesion
	copy(ArrayNombre[:], nombrePropietario)
	//Seteando nombre del propietario, en este caso pertenece al id del usuario en curso
	copy(newAVD.Proper[:], ArrayNombre[:])
	//Ahora hay que buscar un bitmap libre para el nuevo DD, y escribir el nuevo DD
	PosicionEnBitmapDD := GetBitmap(file, int(sb.InicioBitMapDDS), int(sb.TotalDDS))
	//Calculamos la posicion del byte en el bitmap DD
	BitmapPos = int(sb.InicioBitMapDDS) + PosicionEnBitmapDD
	//Seteamos el byte donde iniciara el nuevo struct DD
	DDPos := int(sb.InicioDDS) + (int(sb.SizeDD) * (PosicionEnBitmapDD))
	//Seteamos el apuntador de su Detalle Directorio al nuevo AVD
	newAVD.ApuntadorDD = int32(DDPos)
	newAVD.Permisos = 664
	//Ahora toca escribir el nuevo AVD en su posición correspondiente
	file.Seek(int64(AVDPos+1), 0)
	avdp := &newAVD
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, avdp)
	escribirBytes(file, binario3.Bytes())
	//Escribimos un 1 en esa posición del bitmap DD
	file.Seek(int64(BitmapPos+1), 0)
	data = []byte{0x01}
	file.Write(data)
	//Creamos un nuevo DD
	DDaux := estructuras.DD{}
	//Ahora toca escribir el nuevo DD en su posición correspondiente
	file.Seek(int64(DDPos+1), 0)
	ddp := &DDaux
	var binario4 bytes.Buffer
	binary.Write(&binario4, binary.BigEndian, ddp)
	escribirBytes(file, binario4.Bytes())

	//En este punto ya está creada la nueva carpeta con su respectivo DD
	//Ahora toca setear el apuntador al AVDPadre

	//LEEMOS EL AVD PADRE
	AVDAux := estructuras.AVD{}
	PosPadre := AVDPadre
	file.Seek(int64(AVDPadre+1), 0)
	PadreData := leerBytes(file, int(sb.SizeAVD))
	buffer5 := bytes.NewBuffer(PadreData)
	err := binary.Read(buffer5, binary.BigEndian, &AVDAux)
	if err != nil {
		file.Close()
		fmt.Println(err)
	}

	Continuar := true
	//Recorremos el struct y si el apuntador indirecto a punta a otro AVD tambien lo recorreremos en caso que
	//todos los apuntadores esten ocupados
	for Continuar {
		//Iteramos en las 6 posiciones del arreglo de subdirectoios (apuntadores)
		for i := 0; i < 6; i++ {
			//Validamos que el apuntador no este apuntando a algo
			if AVDAux.ApuntadorSubs[i] == 0 {
				AVDAux.ApuntadorSubs[i] = int32(AVDPos)
				break
			}

		}
		//Si todos los apuntadores en el arreglo están ocupados (apuntando a un AVD)
		//verificamos si el AVD actual apunta hacia otro AVD con otros 6 apuntadores

		if AVDAux.ApuntadorAVD > 0 {
			//Leemos el AVD (que se considera contiguo)
			file.Seek(int64(AVDAux.ApuntadorAVD+int32(1)), 0)
			PosPadre = int(AVDAux.ApuntadorAVD)
			PadreData = leerBytes(file, int(sb.SizeAVD))
			buffer2 := bytes.NewBuffer(PadreData)
			err = binary.Read(buffer2, binary.BigEndian, &AVDAux)
			if err != nil {
				file.Close()
				fmt.Println(err)
				return
			}
		} else {
			//Si llega a este punto significa que aun no se ha asignado el apuntador
			//por lo tanto hay que crear un nuevo AVD y enlazarlo con AVDaux

			//Buscamos la posicion en el bitmap para el nuevo AVD
			PosicionEnBitmapAVD := GetBitmap(file, int(sb.InicioBitmapAVDS), int(sb.TotalAVDS))
			//Calculamos la posicion del byte en el bitmap AVD
			BitmapPos := int(sb.InicioBitmapAVDS) + PosicionEnBitmapAVD
			//Escribimos un 1 en esa posición del bitmap
			file.Seek(int64(BitmapPos+1), 0)
			data := []byte{0x01}
			file.Write(data)
			//Calculamos la posicion del byte del nuevo AVD
			AVDPos2 := int(sb.InicioAVDS) + (int(sb.SizeAVD) * (PosicionEnBitmapAVD))
			//Creamos el nuevo AVD
			newAVD2 := estructuras.AVD{}
			//Como este AVD es una extensión de AVDaux, los atributos serán los mismos
			//además este nuevo AVD no tiene DD, solo sirve para usar sus 6 apuntadores a subs
			copy(newAVD2.FechaCreacion[:], AVDAux.FechaCreacion[:])
			copy(newAVD2.NombreDir[:], AVDAux.NombreDir[:])
			copy(newAVD2.Proper[:], AVDAux.Proper[:])
			newAVD2.ApuntadorSubs[0] = int32(AVDPos)
			newAVD2.Permisos = AVDAux.Permisos
			//Ahora toca escribir el nuevo AVD en su posición correspondiente
			file.Seek(int64(AVDPos2+1), 0)
			avdp := &newAVD2
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, avdp)
			escribirBytes(file, binario3.Bytes())

			AVDAux.ApuntadorAVD = int32(AVDPos2)
			Continuar = false
			break
		}

	}

	//Reescribimos el AVD Padre
	file.Seek(int64(PosPadre+1), 0)
	appadre := &AVDAux
	var binario6 bytes.Buffer
	binary.Write(&binario6, binary.BigEndian, appadre)
	escribirBytes(file, binario6.Bytes())

}

//GetBitmap busca el primer bitmap libre para el nuevo diretorio
func GetBitmap(file *os.File, BitmapStart int, BitmapSize int) int {

	file.Seek(int64(BitmapStart+1), 0)
	BitmapData := leerBytes(file, BitmapSize)
	for i, b := range BitmapData {
		if b == 0 {
			return i
		}
	}

	return -1
}
