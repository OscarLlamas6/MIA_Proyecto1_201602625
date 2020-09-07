package funciones

import (
	"Proyecto1/estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/doun/terminal/color"
)

//EjecutarMkfile inicia la creación de un nuevo grupo
func EjecutarMkfile(id string, path string, size string, cont string, p string) {

	if sesionActiva {

		if path != "" && id != "" {

			if strings.HasPrefix(path, "/") {

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
							//CREAMOS UN STRUCT TEMPORAL
							AVDAux := estructuras.AVD{}
							SizeAVD := int(unsafe.Sizeof(AVDAux))
							fileMBR.Seek(int64(ApuntadorAVD+1), 0)
							AnteriorData := leerBytes(fileMBR, int(SizeAVD))
							buffer2 := bytes.NewBuffer(AnteriorData)
							err = binary.Read(buffer2, binary.BigEndian, &AVDAux)
							if err != nil {
								fileMBR.Close()
								fmt.Println(err)
								return
							}

							var NombreAnterior [20]byte
							copy(NombreAnterior[:], AVDAux.NombreDir[:])
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
									fileMBR.Seek(int64(ApuntadorAVD+1), 0)
									AnteriorData = leerBytes(fileMBR, int(SizeAVD))
									buffer2 = bytes.NewBuffer(AnteriorData)
									err = binary.Read(buffer2, binary.BigEndian, &AVDAux)
									if err != nil {
										fileMBR.Close()
										fmt.Println(err)
										return
									}
									copy(NombreAnterior[:], AVDAux.NombreDir[:])
									i++
									PathCorrecto = true

								} else {

									color.Printf("@{!g}La carpeta @{!y}%v @{!g}no existe\n", carpetas[i])

									if p != "" {
										//CREAR DIRECTORIO
										//Si entramos a esta parte significa que el directorio requerido no existe Y que en el comando MKDIR
										//se especificó el parámetro de recursividad, es decir debemos crear el directorio (el padre)

										if SB1.FreeAVDS > 0 && SB1.FreeDDS > 0 {

											if sesionRoot || EscrituraPropietarioDir(&AVDAux) || EscrituraGrupoDir(&AVDAux) || EscrituraOtrosDir(&AVDAux) {

												if len(carpetas[i]) <= 20 {

													color.Printf("@{!g}Creando carpeta @{!y}%v\n", carpetas[i])

													CrearDirectorio(fileMBR, &SB1, int(ApuntadorAVD), carpetas[i])
													SB1.FirstFreeAVD = SB1.InicioAVDS + (int32(GetBitmap(fileMBR, int(SB1.InicioBitmapAVDS), int(SB1.TotalAVDS))))
													SB1.FirstFreeDD = SB1.InicioDDS + (int32(GetBitmap(fileMBR, int(SB1.InicioBitMapDDS), int(SB1.TotalDDS))))
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
													color.Println("@{!r} El nombre de la carpeta no puede tener más de 20 caracteres")
													break
												}

											} else {
												PathCorrecto = false
												color.Printf("@{!r} El usuario @{!y}%v @{!m}no tiene permisos de escritura en la carpeta @{!y}%v.\n", idSesion, string(NombreAnterior[:]))
												break
											}

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
								if YaExiste := ExisteFile(carpetas[len(carpetas)-1], int(ApuntadorAVD), PathAux); !YaExiste {

									if SB1.FreeInodos > 0 && SB1.FreeBloques > 0 {

										if sesionRoot || EscrituraPropietarioDir(&AVDAux) || EscrituraGrupoDir(&AVDAux) || EscrituraOtrosDir(&AVDAux) {

											if EsCorrecto, Fsize := SizeCorrecto(size); EsCorrecto {

												if len(carpetas[len(carpetas)-1]) <= 20 {

													copy(NombreAnterior[:], AVDAux.NombreDir[:])
													CrearFile(fileMBR, &SB1, int(AVDAux.ApuntadorDD), carpetas[len(carpetas)-1], Fsize, cont, size)

													//Setear nuevas propiedades del superblock
													SB1.FirstFreeInodo = SB1.InicioInodos + (int32(GetBitmap(fileMBR, int(SB1.InicioBitmapInodos), int(SB1.TotalInodos))))
													SB1.FirstFreeBloque = SB1.InicioBloques + (int32(GetBitmap(fileMBR, int(SB1.InicioBitmapBloques), int(SB1.TotalBloques))))

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
													color.Printf("@{!m}El archivo @{!y}%v @{!m}fue creado con éxito\n", carpetas[len(carpetas)-1])

												} else {
													PathCorrecto = false
													color.Println("@{!r} El nombre del archivo no puede tener más de 20 caracteres")
												}

											} else {
												color.Println("@{!r}El size debe ser mayor que cero.")
											}

										} else {
											PathCorrecto = false
											color.Printf("@{!r} El usuario @{!y}%v @{!m}no tiene permisos de escritura en la carpeta @{!y}%v.\n", idSesion, string(NombreAnterior[:]))
										}

									} else {
										color.Println("@{!r} Ya no hay espacio en el sistema de archivos para crear el archivo. Acción fallida.")
									}

								} else {
									color.Printf("@{!r}El archivo @{!y}%v @{!r}ya existe en la carpeta @{!y}%v.\n", carpetas[len(carpetas)-1], string(NombreAnterior[:]))
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

							//NOS POSICIONAMOS DONDE EMPIEZA EL STRCUT DE LA CARPETA ROOT (primer struct AVD)
							ApuntadorAVD := SB1.InicioAVDS
							//CREAMOS UN STRUCT TEMPORAL
							AVDAux := estructuras.AVD{}
							SizeAVD := int(unsafe.Sizeof(AVDAux))
							fileMBR.Seek(int64(ApuntadorAVD+1), 0)
							AnteriorData := leerBytes(fileMBR, int(SizeAVD))
							buffer2 := bytes.NewBuffer(AnteriorData)
							err = binary.Read(buffer2, binary.BigEndian, &AVDAux)
							if err != nil {
								fileMBR.Close()
								fmt.Println(err)
								return
							}

							var NombreAnterior [20]byte
							copy(NombreAnterior[:], AVDAux.NombreDir[:])
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
									fileMBR.Seek(int64(ApuntadorAVD+1), 0)
									AnteriorData = leerBytes(fileMBR, int(SizeAVD))
									buffer2 = bytes.NewBuffer(AnteriorData)
									err = binary.Read(buffer2, binary.BigEndian, &AVDAux)
									if err != nil {
										fileMBR.Close()
										fmt.Println(err)
										return
									}
									copy(NombreAnterior[:], AVDAux.NombreDir[:])
									i++
									PathCorrecto = true

								} else {

									color.Printf("@{!g}La carpeta @{!y}%v @{!g}no existe\n", carpetas[i])

									if p != "" {
										//CREAR DIRECTORIO
										//Si entramos a esta parte significa que el directorio requerido no existe Y que en el comando MKDIR
										//se especificó el parámetro de recursividad, es decir debemos crear el directorio (el padre)

										if SB1.FreeAVDS > 0 && SB1.FreeDDS > 0 {

											if sesionRoot || EscrituraPropietarioDir(&AVDAux) || EscrituraGrupoDir(&AVDAux) || EscrituraOtrosDir(&AVDAux) {

												if len(carpetas[i]) <= 20 {

													color.Printf("@{!g}Creando carpeta @{!y}%v\n", carpetas[i])

													CrearDirectorio(fileMBR, &SB1, int(ApuntadorAVD), carpetas[i])
													SB1.FirstFreeAVD = SB1.InicioAVDS + (int32(GetBitmap(fileMBR, int(SB1.InicioBitmapAVDS), int(SB1.TotalAVDS))))
													SB1.FirstFreeDD = SB1.InicioDDS + (int32(GetBitmap(fileMBR, int(SB1.InicioBitMapDDS), int(SB1.TotalDDS))))

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
													color.Println("@{!r} El nombre de la carpeta no puede tener más de 20 caracteres")
													break
												}

											} else {
												PathCorrecto = false
												color.Printf("@{!r} El usuario @{!y}%v @{!m}no tiene permisos de escritura en la carpeta @{!y}%v.\n", idSesion, string(NombreAnterior[:]))
												break
											}

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
								if YaExiste := ExisteFile(carpetas[len(carpetas)-1], int(ApuntadorAVD), PathAux); !YaExiste {

									if SB1.FreeInodos > 0 && SB1.FreeBloques > 0 {

										if sesionRoot || EscrituraPropietarioDir(&AVDAux) || EscrituraGrupoDir(&AVDAux) || EscrituraOtrosDir(&AVDAux) {

											if EsCorrecto, Fsize := SizeCorrecto(size); EsCorrecto {

												if len(carpetas[len(carpetas)-1]) <= 20 {

													copy(NombreAnterior[:], AVDAux.NombreDir[:])
													CrearFile(fileMBR, &SB1, int(AVDAux.ApuntadorDD), carpetas[len(carpetas)-1], Fsize, cont, size)

													//Setear nuevas propiedades del superblock
													SB1.FirstFreeInodo = SB1.InicioInodos + (int32(GetBitmap(fileMBR, int(SB1.InicioBitmapInodos), int(SB1.TotalInodos))))
													SB1.FirstFreeBloque = SB1.InicioBloques + (int32(GetBitmap(fileMBR, int(SB1.InicioBitmapBloques), int(SB1.TotalBloques))))

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
													color.Printf("@{!m}El archivo @{!y}%v @{!m}fue creado con éxito\n", carpetas[len(carpetas)-1])

												} else {
													PathCorrecto = false
													color.Println("@{!r} El nombre del archivo no puede tener más de 20 caracteres")
												}

											} else {
												color.Println("@{!r}El size debe ser mayor que cero.")
											}

										} else {
											PathCorrecto = false
											color.Printf("@{!r} El usuario @{!y}%v @{!m}no tiene permisos de escritura en la carpeta @{!y}%v.\n", idSesion, string(NombreAnterior[:]))
										}

									} else {
										color.Println("@{!r} Ya no hay espacio en el sistema de archivos para crear el archivo. Acción fallida.")
									}

								} else {
									color.Printf("@{!r}El archivo @{!y}%v @{!r}ya existe en la carpeta @{!y}%v.\n", carpetas[len(carpetas)-1], string(NombreAnterior[:]))
								}
							} else {
								color.Println("@{!r} Error, una o más carpetas padre no existen.")
							}

						} else {
							color.Println("@{!r} La partición indicada no ha sido formateada.")
						}

						fileMBR.Close()

					}

				} else {
					color.Printf("@{!r}No hay ninguna partición montada con el id: @{!y}%v\n", id)
				}

			} else {
				color.Println("@{!r}Path incorrecto, debe iniciar con @{!y}/")
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

//CrearFile crea un archivo en el directorio especificado (DDPadre)
func CrearFile(file *os.File, sb *estructuras.Superblock, DDPadre int, nombre string, size int, cont string, cadenaSize string) {

	//Buscamos la posición en el bitmap para el nuevo Inodo
	PosicionEnBitmapInodo := GetBitmap(file, int(sb.InicioBitmapInodos), int(sb.TotalInodos))
	//Calculamos la posicion del byte en el bitmap Inodos
	BitmapPos := int(sb.InicioBitmapInodos) + PosicionEnBitmapInodo
	//Escribimos un 1 en esa posición del bitmap
	file.Seek(int64(BitmapPos+1), 0)
	data := []byte{0x01}
	file.Write(data)
	//Calculamos la posición del byte del nuevo Inodo
	InodoPos := int(sb.InicioInodos) + (int(sb.SizeInodo) * (PosicionEnBitmapInodo))
	//Creamos el nuevo Inodo
	newInodo := estructuras.Inodo{}
	//Seteando nombre del propietario, en este caso pertenece al id del usuario en curso
	var ArrayProper [20]byte
	nombrePropietario := idSesion
	copy(ArrayProper[:], nombrePropietario)
	copy(newInodo.Proper[:], ArrayProper[:])
	//Seteando nombre del grupo, en este caso pertenece al id del grupo en curso
	var ArrayGrupo [20]byte
	nombreGrupo := idGrupo
	copy(ArrayGrupo[:], nombreGrupo)
	copy(newInodo.Proper[:], ArrayGrupo[:])
	newInodo.NumeroInodo = int32(int32(sb.TotalInodos)-int32(sb.FreeInodos)) + 1
	newInodo.PermisoU = 6
	newInodo.PermisoG = 6
	newInodo.PermisoO = 4

	//Seteamos SB
	sb.FreeInodos = sb.FreeInodos - 1

	//En este punto ya está creado el nuevo inodo
	//Ahora toca setear el apuntador al DD

	//LEEMOS EL DD PADRE
	DDAux := estructuras.DD{}
	PosPadre := DDPadre
	file.Seek(int64(PosPadre+1), 0)
	PadreData := leerBytes(file, int(sb.SizeDD))
	buffer := bytes.NewBuffer(PadreData)
	err := binary.Read(buffer, binary.BigEndian, &DDAux)
	if err != nil {
		file.Close()
		fmt.Println(err)
	}

	Continuar := true

	//Recorremos el struct y el apuntador indirecto indirecto apunta a otro DD lo recorremos en caso
	//que todos los apuntadores ya estén apuntando a un inodo
	for Continuar {
		//Iteramos en las 5 posiciones del arreglo de archivos que tiene el DD
		for i := 0; i < 5; i++ {
			//Validamos que el apuntador al inodo si esté apuntando a algo
			if DDAux.DDFiles[i].ApuntadorInodo == 0 {
				//Seteamos los datos del nuevo archivo
				var DDnombre [20]byte
				copy(DDnombre[:], nombre)
				copy(DDAux.DDFiles[i].Name[:], DDnombre[:])
				var chars [20]byte
				t := time.Now()
				cadena := t.Format("2006-01-02 15:04:05")
				copy(chars[:], cadena)
				copy(DDAux.DDFiles[i].FechaCreacion[:], chars[:])
				copy(DDAux.DDFiles[i].FechaModificacion[:], chars[:])
				DDAux.DDFiles[i].ApuntadorInodo = int32(InodoPos)
				Continuar = false
				break
			}

		}

		if Continuar == false {
			break
		}

		//Si todos los apuntadores en el arreglo están ocupados (apuntando a un inodo)
		//verificamos si el DD actual apunta hacia otro DD con otros 5 apuntadores

		if DDAux.ApuntadorDD > 0 {

			//Leemos el DD (que se considera contiguo)
			file.Seek(int64(DDAux.ApuntadorDD+int32(1)), 0)
			PosPadre = int(DDAux.ApuntadorDD)
			PadreData = leerBytes(file, int(sb.SizeDD))
			buffer = bytes.NewBuffer(PadreData)
			err = binary.Read(buffer, binary.BigEndian, &DDAux)
			if err != nil {
				file.Close()
				fmt.Println(err)
				return
			}

		} else {

			//Si llega a este punto significa que aun no se ha asignado el apuntador
			//por lo tanto hay que crear un nuevo DD y enlazarlo con DDaux

			//Ahora hay que buscar un bitmap libre para el nuevo DD, y escribir el nuevo DD
			PosicionEnBitmapDD := GetBitmap(file, int(sb.InicioBitMapDDS), int(sb.TotalDDS))
			//Calculamos la posicion del byte en el bitmap DD
			BitmapPos = int(sb.InicioBitMapDDS) + PosicionEnBitmapDD
			//Escribimos un 1 en esa posición del bitmap
			file.Seek(int64(BitmapPos+1), 0)
			data := []byte{0x01}
			file.Write(data)
			//Seteamos el byte donde iniciara el nuevo struct DD
			DDPos := int(sb.InicioDDS) + (int(sb.SizeDD) * (PosicionEnBitmapDD))
			//Creamos el nuevo DD
			newDD2 := estructuras.DD{}
			//Como este DD está nuevo, tenemos la certeza que la posición cero en el arreglo está desocupada
			var DDnombre [20]byte
			copy(DDnombre[:], nombre)
			copy(newDD2.DDFiles[0].Name[:], DDnombre[:])
			var chars [20]byte
			t := time.Now()
			cadena := t.Format("2006-01-02 15:04:05")
			copy(chars[:], cadena)
			copy(newDD2.DDFiles[0].FechaCreacion[:], chars[:])
			copy(newDD2.DDFiles[0].FechaModificacion[:], chars[:])
			newDD2.DDFiles[0].ApuntadorInodo = int32(InodoPos)

			//Actualizamos el SB
			sb.FreeDDS = sb.FreeDDS - 1

			//Ahora toca escribir el nuevo DD en su posición correspondiente
			file.Seek(int64(DDPos+1), 0)
			ddp := &newDD2
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, ddp)
			escribirBytes(file, binario2.Bytes())

			DDAux.ApuntadorDD = int32(DDPos)
			Continuar = false
			break
		}
	}

	//Reescribimos el DD Padre
	file.Seek(int64(PosPadre+1), 0)
	ApPadre := &DDAux
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, ApPadre)
	escribirBytes(file, binario.Bytes())

	//En este punto el inodo correspondiente al nuevo archivo
	//Ya está seteado al apuntador del DD correspondiente
	//Ahora toca escribir el contenido en los bloques de datos
	//y enlazarlos con newInodo, cabe recalcar que se debe validar
	//la cantidad de bloques necesarios, para crear un nuevo inodo de ser necesario

	if cadenaSize == "" && cont == "" { //Si no vienen nungun parámetro ni size ni cont, el tamaño de 0

		//Si el size es cero, el archivo no contiene datos, por lo tanto no se le asigna ningun bloque
		//y podemos proceder a guardarlo inmediatamente

		newInodo.FileSize = 0
		newInodo.NumeroBloques = 0
		newInodo.ApuntadorIndirecto = 0

		//Ahora toca escribir el struct Inodo en su posición correspondiente
		file.Seek(int64(InodoPos+1), 0)
		inodop := &newInodo
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, inodop)
		escribirBytes(file, binario.Bytes())

	} else {

		if size == 0 {

			//Si el size es cero, el archivo no contiene datos, por lo tanto no se le asigna ningun bloque
			//y podemos proceder a guardarlo inmediatamente

			newInodo.FileSize = 0
			newInodo.NumeroBloques = 0
			newInodo.ApuntadorIndirecto = 0

			//Ahora toca escribir el struct Inodo en su posición correspondiente
			file.Seek(int64(InodoPos+1), 0)
			inodop := &newInodo
			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, inodop)
			escribirBytes(file, binario.Bytes())

		} else {

			contenido := ""
			FileSize := 0

			if cont == "" { //No hay contenido (setear el abcdario)

				contenido = getNTimesabc(size)
				FileSize = size

			} else if cadenaSize == "" && cont != "" {

				contenido = cont
				FileSize = len(contenido)

			} else {

				if size > len(cont) { //size es mayor que el tamaño de cont, completar cont con el abcdario hasta llegar a size

					FileSize = size
					contenido = cont             //contenido se le concatena el contenido enviado como parametro, que tiene un tamaño menor a size
					r := size - len(cont)        //calculamos cuantos caraceres hay que agregarle a contenido para cumplir con el size
					contenido += getNTimesabc(r) //a contenido le enviamos los caraćteres del abcdario necesarios para llegar al size

				} else if size == len(cont) { //si size es igual al tamaño de cont, llenamos los bloques con cont

					FileSize = size
					contenido = cont

				} else if size < len(cont) { //si size es menor que el tamaño de cont, cortamos cont hasta el tamaño de size

					FileSize = size
					contenido = contenido[:size]

				}

			}

			//En este momento el contenido cumple con el size requerido y está listo para ser seteado
			//FileSize será la variable a utilizar para setear el atributo en el Inodo

			//Dividimos la cantidad de caracteres que tiene contenido dentro de 25
			//si el resultado es un numero decimal lo aproximamos al entero más cercano a la derecha
			//esto para saber cuando bloques de datos vamos a usar
			//Ejemplo: Si tenemos 52 caracteres, al dividirlo dentro de 25 obtenemos 2.08
			//ese 0.08 nos obliga a tomar un 3er bloque, entonces la funcion Roundf, aproximaria 2.08 a 3
			resultado := float64(len(contenido)) / float64(25)
			resultado = Roundf(resultado)
			//esta variable nos ayudara al corrimiento de los caracteres
			x := 0
			indx := 0
			for i := 0; i < int(resultado); i++ {

				//Creamos un bloque datos, escribimos en su bitmap, creamos el struct, el asignamos los datos desde x hasta x+25
				//si el index llega más de 3, significa que necesitaremos otro inodo

				//Buscamos la posicion en el bitmap para el nuevo BloqueDatos
				PosicionEnBitmapBloque := GetBitmap(file, int(sb.InicioBitmapBloques), int(sb.TotalBloques))
				//Calculamos la posicion del byte en el bitmap BloqueDatos
				BitmapPos := int(sb.InicioBitmapBloques) + PosicionEnBitmapBloque
				//Escribimos un 1 en esa posición del bitmap
				file.Seek(int64(BitmapPos+1), 0)
				data := []byte{0x01}
				file.Write(data)
				//Calculamos la posición del byte del nuevo Bloque Datos
				BloquePos := int(sb.InicioBloques) + (int(sb.SizeBloque) * (PosicionEnBitmapBloque))
				//Creamos el nuevo Bloque
				newBloque := estructuras.BloqueDatos{}
				//Le pasamos el contenido desde x hasta x+25

				if i != int(resultado-1) {
					copy(newBloque.Data[:], contenido[x:x+25])
				} else {
					copy(newBloque.Data[:], contenido[x:])
				}

				//Ahora toca escribir el struct BloqueDatos en su posición correspondiente
				file.Seek(int64(BloquePos+1), 0)
				bloquep := &newBloque
				var binario6 bytes.Buffer
				binary.Write(&binario6, binary.BigEndian, bloquep)
				escribirBytes(file, binario6.Bytes())

				//Actualizamos el SB
				sb.FreeBloques = sb.FreeBloques - 1

				//Asignamos al inodo el apuntador al bloque que creamos
				newInodo.ApuntadoresBloques[indx] = int32(BloquePos)

				x += 25
				indx++
				if indx > 3 && (i+1) < int(resultado) {
					//Si el index es mayor a 3 y estamos seguros que daremos otra iteración más
					//eso significa que necesitamos calcula la posición para un nuevo inodo
					//esa posicion la seteamos en el apuntador indirecto de nuestro newInodo
					//escribimos nuestro newinodo en el archivo, seguido de esto, creamos otro inodo
					//y newInodo apuntaria a un nuevoInodo, finalmente reseteamos indx a 0
					//para que pueda comenzar a apuntar desde la primera posición en el arreglo de apuntadores
					//de bloques del nuevo Inodo

					//Buscamos la posición en el bitmap para el nuevo Inodo
					PosicionEnBitmapInodo := GetBitmap(file, int(sb.InicioBitmapInodos), int(sb.TotalInodos))
					//Calculamos la posicion del byte en el bitmap Inodos
					BitmapPos := int(sb.InicioBitmapInodos) + PosicionEnBitmapInodo
					//Escribimos un 1 en esa posición del bitmap
					file.Seek(int64(BitmapPos+1), 0)
					data := []byte{0x01}
					file.Write(data)
					//Calculamos la posición del byte del nuevo Inodo
					Inodo2Pos := int(sb.InicioInodos) + (int(sb.SizeInodo) * (PosicionEnBitmapInodo))

					//Asignamos la posición del nuevo inodo a nuestro newInodo original

					newInodo.ApuntadorIndirecto = int32(Inodo2Pos)
					newInodo.FileSize = int32(FileSize)
					newInodo.NumeroBloques = int32(resultado)

					//Ahora toca escribir el struct Inodo en su posición correspondiente
					file.Seek(int64(InodoPos+1), 0)
					inodop := &newInodo
					var binario bytes.Buffer
					binary.Write(&binario, binary.BigEndian, inodop)
					escribirBytes(file, binario.Bytes())

					//Ahora creamos el nuevo inodo
					newInodo = estructuras.Inodo{}
					InodoPos = Inodo2Pos

					//Seteando nombre del propietario, en este caso pertenece al id del usuario en curso
					var ArrayProper [20]byte
					nombrePropietario = idSesion
					copy(ArrayProper[:], nombrePropietario)
					copy(newInodo.Proper[:], ArrayProper[:])
					//Seteando nombre del grupo, en este caso pertenece al id del grupo en curso
					var ArrayGrupo [20]byte
					nombreGrupo = idGrupo
					copy(ArrayGrupo[:], nombreGrupo)
					copy(newInodo.Proper[:], ArrayGrupo[:])
					newInodo.NumeroInodo = int32(int32(sb.TotalInodos)-int32(sb.FreeInodos)) + 1
					newInodo.PermisoU = 6
					newInodo.PermisoG = 6
					newInodo.PermisoO = 4

					//Seteamos SB
					sb.FreeInodos = sb.FreeInodos - 1

					indx = 0
				}
			}

			newInodo.FileSize = int32(FileSize)
			newInodo.NumeroBloques = int32(resultado)

			//Ahora toca escribir el ultimo strcut Inodo creado en su posición correspondiente
			file.Seek(int64(InodoPos+1), 0)
			inodop := &newInodo
			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, inodop)
			escribirBytes(file, binario.Bytes())

		}

	}

}

func getNTimesabc(c int) string {
	resultado := ""
	x := 1
	for i := 0; i < c; i++ {
		resultado += getABC(x)
		x++
		if x > 26 {
			x = 1
		}
	}
	return resultado

}

//SizeCorrecto verifica si el size es correcto en caso de venir como parametro
func SizeCorrecto(size string) (bool, int) {

	if size != "" {

		if Fsize, _ := strconv.Atoi(size); Fsize >= 0 {
			return true, Fsize
		}
		return false, 0

	}

	return true, 0

}
