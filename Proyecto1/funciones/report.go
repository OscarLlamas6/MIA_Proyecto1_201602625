package funciones

import (
	"Proyecto1/estructuras"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"github.com/doun/terminal/color"
)

//EjecutarReporte verifica el tipo de reporte segun el parametro NOMBRE
func EjecutarReporte(nombre string, path string, ruta string, id string) {

	if path != "" && nombre != "" && id != "" {
		if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil { //verificamos que se pueda construir el path
			fmt.Printf("Path invalido")
		} else {

			if IDYaRegistrado(id) { //verificamos que el id si exista, osea que haya una particion montada con ese id
				if strings.ToLower(nombre) == "mbr" {
					ReporteMBR(path, ruta, id)
				} else if strings.ToLower(nombre) == "disk" {
					ReporteDisk(path, ruta, id)
				} else if strings.ToLower(nombre) == "bm_arbdir" {
					ReporteBitmapAVD(path, ruta, id)
				} else if strings.ToLower(nombre) == "bm_detdir" {
					ReporteBitmapDD(path, ruta, id)
				} else if strings.ToLower(nombre) == "bm_inode" {
					ReporteBitmapInode(path, ruta, id)
				} else if strings.ToLower(nombre) == "bm_block" {
					ReporteBitmapBloque(path, ruta, id)
				} else if strings.ToLower(nombre) == "sb" {
					ReporteSB(path, ruta, id)
				}
			} else {
				color.Printf("@{!r}No hay ninguna partición montada con el id: @{!y}%v\n", id)
			}
		}
	} else {
		color.Println("@{!r}Faltan parámetros obligatorios para la funcion REP.")
	}

}

//ReporteMBR crea el reporte del mbr
func ReporteMBR(path string, ruta string, id string) {

	extension := filepath.Ext(path)

	if strings.ToLower(extension) == ".pdf" || strings.ToLower(extension) == ".jpg" || strings.ToLower(extension) == ".png" {

		file, err := os.OpenFile("codigo.dot", os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}
		// Change permissions Linux.
		err = os.Chmod("codigo.dot", 0666)
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}

		file.Truncate(0)
		file.Seek(0, 0)

		_, err = file.WriteString("digraph H {\n node [ shape=plain] \n table [ label = <\n  <table border='1' cellborder='1'>\n   <tr><td>Nombre</td><td>Valor</td></tr>\n")
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}

		//LEER Y RECORRER EL MBR
		_, PathAux := GetDatosPart(id)
		fileMBR, err2 := os.Open(PathAux)
		if err2 != nil { //validar que no sea nulo.
			panic(err2)
		}
		Disco1 := estructuras.MBR{}
		DiskSize := int(unsafe.Sizeof(Disco1))
		DiskData := leerBytes(fileMBR, DiskSize)
		buffer := bytes.NewBuffer(DiskData)
		err = binary.Read(buffer, binary.BigEndian, &Disco1)
		if err != nil {
			fileMBR.Close()
			fmt.Println(err)
			return
		}
		fileMBR.Close()

		w := bufio.NewWriter(file)

		fmt.Fprintf(w, "   <tr><td>MBR_Tamanio</td><td>%v</td></tr>\n", Disco1.Msize)
		fmt.Fprintf(w, "   <tr><td>MBR_Fecha_Creación</td><td>%v</td></tr>\n", string(Disco1.Mdate[:len(Disco1.Mdate)-1]))
		fmt.Fprintf(w, "   <tr><td>MBR_Disk_Signature</td><td>%v</td></tr>\n", Disco1.Msignature)

		PartNum := 1
		for i := 0; i < 4; i++ {
			if Disco1.Mpartitions[i].Psize > 0 {
				fmt.Fprintf(w, "   <tr><td>Part_%d_Status</td><td>%v</td></tr>\n", PartNum, string(Disco1.Mpartitions[i].Pstatus))
				fmt.Fprintf(w, "   <tr><td>Part_%d_Type</td><td>%v</td></tr>\n", PartNum, string(Disco1.Mpartitions[i].Ptype))
				fmt.Fprintf(w, "   <tr><td>Part_%d_Fit</td><td>%v</td></tr>\n", PartNum, string(Disco1.Mpartitions[i].Pfit))
				fmt.Fprintf(w, "   <tr><td>Part_%d_Start</td><td>%v</td></tr>\n", PartNum, Disco1.Mpartitions[i].Pstart)
				n := bytes.Index(Disco1.Mpartitions[i].Pname[:], []byte{0})
				fmt.Fprintf(w, "   <tr><td>Part_%d_Name</td><td>%v</td></tr>\n", PartNum, string(Disco1.Mpartitions[i].Pname[:n]))
				PartNum++
			}
		}

		w.Flush()
		////////////////////
		_, err = file.WriteString("  </table>\n > ]\n}")
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}

		file.Close()

		extT := "-T"

		switch strings.ToLower(extension) {
		case ".png":
			extT += "png"
		case ".pdf":
			extT += "pdf"
		case ".jpg":
			extT += "jpg"
		default:

		}

		if runtime.GOOS == "windows" {
			cmd := exec.Command("dot", extT, "-o", path, "codigo.dot") //Windows example, its tested
			//cmd.Stdout = os.Stdout
			cmd.Run()
		} else {
			cmd := exec.Command("dot", extT, "-o", path, "codigo.dot") //Linux example, its tested
			//cmd.Stdout = os.Stdout
			cmd.Run()
		}
	} else {
		color.Println("@{!r}El reporte MBR solo puede generar archivos con extensión .png, .jpg ó .pdf.")
	}

}

//ReporteDisk crea el reporte de las particiones del disco
func ReporteDisk(path string, ruta string, id string) {

	extension := filepath.Ext(path)

	if strings.ToLower(extension) == ".pdf" || strings.ToLower(extension) == ".jpg" || strings.ToLower(extension) == ".png" {

		file, err := os.OpenFile("codigo.dot", os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}
		// Change permissions Linux.
		err = os.Chmod("codigo.dot", 0666)
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}

		file.Truncate(0)
		file.Seek(0, 0)

		_, err = file.WriteString("digraph D {\nparent [\n shape=plaintext\n label=<\n<table border='1' cellborder='1'>\n <tr>\n  <td rowspan=\"2\" bgcolor=\"orange\">MBR</td>\n")
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}

		//LEER Y RECORRER EL MBR
		_, PathAux := GetDatosPart(id)
		fileMBR, err2 := os.Open(PathAux)
		if err2 != nil { //validar que no sea nulo.
			panic(err2)
		}
		Disco1 := estructuras.MBR{}
		DiskSize := int(unsafe.Sizeof(Disco1))
		DiskData := leerBytes(fileMBR, DiskSize)
		buffer := bytes.NewBuffer(DiskData)
		err = binary.Read(buffer, binary.BigEndian, &Disco1)
		if err != nil {
			fileMBR.Close()
			fmt.Println(err)
			return
		}
		fileMBR.Close()

		w := bufio.NewWriter(file)
		HayExtendida := false
		nLogicas := 0
		for i := 0; i < 4; i++ {
			if Disco1.Mpartitions[i].Psize > 0 {

				if Disco1.Mpartitions[i].Ptype == 'e' || Disco1.Mpartitions[i].Ptype == 'E' {
					HayExtendida = true
					nLogicas = CantidadLogicas(PathAux)
					Columnas := 0
					if nLogicas > 0 {
						Columnas = 2 * nLogicas
					} else {
						Columnas = 2
					}

					if i > 0 {

						if Disco1.Mpartitions[i-1].Psize == 0 {
							n := bytes.Index(Disco1.Mpartitions[i].Pname[:], []byte{0})
							fmt.Fprintf(w, "  <td colspan=\"%d\" bgcolor=\"#0cff04\">%v - Extendida</td>\n", Columnas, string(Disco1.Mpartitions[i].Pname[:n]))
						} else {

							if (Disco1.Mpartitions[i-1].Pstart + Disco1.Mpartitions[i-1].Psize) == Disco1.Mpartitions[i].Pstart {
								n := bytes.Index(Disco1.Mpartitions[i].Pname[:], []byte{0})
								fmt.Fprintf(w, "  <td colspan=\"%d\" bgcolor=\"#0cff04\">%v - Extendida</td>\n", Columnas, string(Disco1.Mpartitions[i].Pname[:n]))
							} else {
								fmt.Fprint(w, "  <td rowspan=\"2\" bgcolor=\"skyblue\">Libre</td>")
								n := bytes.Index(Disco1.Mpartitions[i].Pname[:], []byte{0})
								fmt.Fprintf(w, "  <td colspan=\"%d\" bgcolor=\"#0cff04\">%v - Extendida</td>\n", Columnas, string(Disco1.Mpartitions[i].Pname[:n]))
							}

						}

					} else {
						n := bytes.Index(Disco1.Mpartitions[i].Pname[:], []byte{0})
						fmt.Fprintf(w, "  <td colspan=\"%d\" bgcolor=\"#0cff04\">%v - Extendida</td>\n", Columnas, string(Disco1.Mpartitions[i].Pname[:n]))
					}

				} else {

					if i > 0 {

						if Disco1.Mpartitions[i-1].Psize == 0 {
							n := bytes.Index(Disco1.Mpartitions[i].Pname[:], []byte{0})
							fmt.Fprintf(w, "  <td rowspan=\"2\" bgcolor=\"yellow\">%v - Primaria</td>\n", string(Disco1.Mpartitions[i].Pname[:n]))
						} else {

							if (Disco1.Mpartitions[i-1].Pstart + Disco1.Mpartitions[i-1].Psize) == Disco1.Mpartitions[i].Pstart {
								n := bytes.Index(Disco1.Mpartitions[i].Pname[:], []byte{0})
								fmt.Fprintf(w, "  <td rowspan=\"2\" bgcolor=\"yellow\">%v - Primaria</td>\n", string(Disco1.Mpartitions[i].Pname[:n]))
							} else {
								fmt.Fprint(w, "  <td rowspan=\"2\" bgcolor=\"skyblue\">Libre</td>")
								n := bytes.Index(Disco1.Mpartitions[i].Pname[:], []byte{0})
								fmt.Fprintf(w, "  <td rowspan=\"2\" bgcolor=\"yellow\">%v - Primaria</td>\n", string(Disco1.Mpartitions[i].Pname[:n]))
							}
						}

					} else {
						n := bytes.Index(Disco1.Mpartitions[i].Pname[:], []byte{0})
						fmt.Fprintf(w, "  <td rowspan=\"2\" bgcolor=\"yellow\">%v - Primaria</td>\n", string(Disco1.Mpartitions[i].Pname[:n]))
					}
				}

				if i == 3 {
					if Disco1.Mpartitions[i].Pstart+Disco1.Mpartitions[i].Psize != Disco1.Msize-1 {
						fmt.Fprint(w, "  <td rowspan=\"2\" bgcolor=\"skyblue\">Libre</td>\n")
					}

				}

			} else {

				if i > 0 {
					if Disco1.Mpartitions[i-1].Psize != 0 {
						fmt.Fprint(w, "  <td rowspan=\"2\" bgcolor=\"skyblue\">Libre</td>\n")
					}
				} else {
					fmt.Fprint(w, "  <td rowspan=\"2\" bgcolor=\"skyblue\">Libre</td>\n")
				}

			}
		}
		fmt.Fprint(w, " </tr>\n")

		if HayExtendida {
			fmt.Fprint(w, " <tr>\n")

			////////////////////RECORRER EBRS

			indiceExt, _ := IndiceExtendida(PathAux)

			fileEBR, err := os.OpenFile(PathAux, os.O_RDWR, 0666)
			if err != nil {
				fmt.Println(err)
				fileEBR.Close()
			}

			EBRaux := estructuras.EBR{}
			EBRSize := int(unsafe.Sizeof(EBRaux))
			fileEBR.Seek(int64(indiceExt)+1, 0)
			EBRData := leerBytes(fileEBR, EBRSize)
			buffer := bytes.NewBuffer(EBRData)

			err = binary.Read(buffer, binary.BigEndian, &EBRaux)
			if err != nil {
				fileEBR.Close()
				panic(err)
			}

			Continuar := true

			for Continuar {

				if EBRaux.Esize > 0 {
					fmt.Fprint(w, "  <td>EBR</td>\n")
					fmt.Fprint(w, "  <td>LOGICA</td>\n")
				} else {
					fmt.Fprint(w, "  <td>EBR</td>\n")
					fmt.Fprint(w, "  <td>LIBRE</td>\n")
				}

				if EBRaux.Enext != -1 {
					//Si hay otro EBR a la derecha lo leemos y volvemos al inicio del FOR
					fileEBR.Seek(int64(EBRaux.Enext)+1, 0)
					EBRData := leerBytes(fileEBR, EBRSize)
					buffer := bytes.NewBuffer(EBRData)
					err = binary.Read(buffer, binary.BigEndian, &EBRaux)
					if err != nil {
						fileEBR.Close()
						panic(err)
					}
				} else {
					//Si no cancelamos el loop
					Continuar = false
				}

			}
			fileEBR.Close()

			////////////////////////////////////
			fmt.Fprint(w, " </tr>\n")
		}

		w.Flush()

		////////////////////
		_, err = file.WriteString("  </table>\n >]\n}")
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}

		file.Close()

		extT := "-T"

		switch strings.ToLower(extension) {
		case ".png":
			extT += "png"
		case ".pdf":
			extT += "pdf"
		case ".jpg":
			extT += "jpg"
		default:

		}

		if runtime.GOOS == "windows" {
			cmd := exec.Command("dot", extT, "-o", path, "codigo.dot") //Windows example, its tested
			//cmd.Stdout = os.Stdout
			cmd.Run()
		} else {
			cmd := exec.Command("dot", extT, "-o", path, "codigo.dot") //Linux example, its tested
			//cmd.Stdout = os.Stdout
			cmd.Run()
		}
	} else {
		color.Println("@{!r}El reporte DISK solo puede generar archivos con extensión .png, .jpg ó .pdf.")
	}
}

//ReporteSB crea el reporte del super bloque
func ReporteSB(path string, ruta string, id string) {

	extension := filepath.Ext(path)

	if strings.ToLower(extension) == ".pdf" || strings.ToLower(extension) == ".jpg" || strings.ToLower(extension) == ".png" {

		NameAux, PathAux := GetDatosPart(id)

		if Existe, Indice := ExisteParticion(PathAux, NameAux); Existe {

			//LEER Y RECORRER EL MBR
			fileMBR, err2 := os.Open(PathAux)
			if err2 != nil { //validar que no sea nulo.
				panic(err2)
			}
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
			fileMBR.Close()

			if SB1.MontajesCount > 0 {

				file, err := os.OpenFile("codigo.dot", os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}
				// Change permissions Linux.
				err = os.Chmod("codigo.dot", 0666)
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Truncate(0)
				file.Seek(0, 0)

				_, err = file.WriteString("digraph H {\n node [ shape=plain] \n table [ label = <\n  <table border='1' cellborder='1'>\n   <tr><td>Atributo</td><td>Valor</td></tr>\n")
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				w := bufio.NewWriter(file)

				///////////SETEAR DATOS DEL SUPERBLOQUE
				n := bytes.Index(SB1.Name[:], []byte{0})
				fmt.Fprintf(w, "   <tr><td>Nombre</td><td>%v</td></tr>\n", string(SB1.Name[:n]))
				fmt.Fprintf(w, "   <tr><td>Total AVDs</td><td>%v</td></tr>\n", SB1.TotalAVDS)
				fmt.Fprintf(w, "   <tr><td>Total DDs</td><td>%v</td></tr>\n", SB1.TotalDDS)
				fmt.Fprintf(w, "   <tr><td>Total Inodos</td><td>%v</td></tr>\n", SB1.TotalInodos)
				fmt.Fprintf(w, "   <tr><td>Total Bloques</td><td>%v</td></tr>\n", SB1.TotalBloques)
				fmt.Fprintf(w, "   <tr><td>Total Bitacoras</td><td>%v</td></tr>\n", SB1.TotalBitacoras)
				fmt.Fprintf(w, "   <tr><td>Free AVDs</td><td>%v</td></tr>\n", SB1.FreeAVDS)
				fmt.Fprintf(w, "   <tr><td>Free DDs</td><td>%v</td></tr>\n", SB1.FreeDDS)
				fmt.Fprintf(w, "   <tr><td>Free Inodos</td><td>%v</td></tr>\n", SB1.FreeInodos)
				fmt.Fprintf(w, "   <tr><td>Free Bloques</td><td>%v</td></tr>\n", SB1.FreeBloques)
				fmt.Fprintf(w, "   <tr><td>Free Bitacoras</td><td>%v</td></tr>\n", SB1.FreeBitacoras)
				n = bytes.Index(SB1.DateCreacion[:], []byte{0})
				fmt.Fprintf(w, "   <tr><td>Fecha creación</td><td>%v</td></tr>\n", string(SB1.DateCreacion[:n]))
				n = bytes.Index(SB1.DateLastMount[:], []byte{0})
				fmt.Fprintf(w, "   <tr><td>Fecha creación</td><td>%v</td></tr>\n", string(SB1.DateLastMount[:n]))
				fmt.Fprintf(w, "   <tr><td>No. Montajes</td><td>%v</td></tr>\n", SB1.MontajesCount)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bitmap AVDs</td><td>%v</td></tr>\n", SB1.InicioBitmapAVDS)
				fmt.Fprintf(w, "   <tr><td>Apuntador AVDs</td><td>%v</td></tr>\n", SB1.InicioAVDS)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bitmap DDs</td><td>%v</td></tr>\n", SB1.InicioBitMapDDS)
				fmt.Fprintf(w, "   <tr><td>Apuntador DDs</td><td>%v</td></tr>\n", SB1.InicioDDS)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bitmap Inodos</td><td>%v</td></tr>\n", SB1.InicioBitmapInodos)
				fmt.Fprintf(w, "   <tr><td>Apuntador inodos</td><td>%v</td></tr>\n", SB1.InicioInodos)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bitmap bloques</td><td>%v</td></tr>\n", SB1.InicioBitmapBloques)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bloques</td><td>%v</td></tr>\n", SB1.InicioBloques)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bitacoras</td><td>%v</td></tr>\n", SB1.InicioBitacora)
				fmt.Fprintf(w, "   <tr><td>Size struct AVD</td><td>%v</td></tr>\n", SB1.SizeAVD)
				fmt.Fprintf(w, "   <tr><td>Size struct DD</td><td>%v</td></tr>\n", SB1.SizeDD)
				fmt.Fprintf(w, "   <tr><td>Size struct Inodo</td><td>%v</td></tr>\n", SB1.SizeInodo)
				fmt.Fprintf(w, "   <tr><td>Size struct Bloque</td><td>%v</td></tr>\n", SB1.SizeBloque)
				fmt.Fprintf(w, "   <tr><td>Byte primer AVD libre</td><td>%v</td></tr>\n", SB1.FirstFreeAVD)
				fmt.Fprintf(w, "   <tr><td>Byte primer DD libre</td><td>%v</td></tr>\n", SB1.FirstFreeDD)
				fmt.Fprintf(w, "   <tr><td>Byte primer Inodo libre</td><td>%v</td></tr>\n", SB1.FirstFreeInodo)
				fmt.Fprintf(w, "   <tr><td>Byte primer Bloque libre</td><td>%v</td></tr>\n", SB1.FirstFreeBloque)
				fmt.Fprintf(w, "   <tr><td>Magic Number :D </td><td>%v</td></tr>\n", SB1.MagicNum)
				////////////////////

				w.Flush()

				_, err = file.WriteString("  </table>\n > ]\n}")
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Close()

				extT := "-T"

				switch strings.ToLower(extension) {
				case ".png":
					extT += "png"
				case ".pdf":
					extT += "pdf"
				case ".jpg":
					extT += "jpg"
				default:

				}

				if runtime.GOOS == "windows" {
					cmd := exec.Command("dot", extT, "-o", path, "codigo.dot") //Windows example, its tested
					//cmd.Stdout = os.Stdout
					cmd.Run()
				} else {
					cmd := exec.Command("dot", extT, "-o", path, "codigo.dot") //Linux example, its tested
					//cmd.Stdout = os.Stdout
					cmd.Run()
				}

			} else {
				color.Println("@{!r} La partición indicada no ha sido formateada.")
			}

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
			fileMBR.Close()
			if SB1.MontajesCount > 0 {

				file, err := os.OpenFile("codigo.dot", os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}
				// Change permissions Linux.
				err = os.Chmod("codigo.dot", 0666)
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Truncate(0)
				file.Seek(0, 0)

				_, err = file.WriteString("digraph H {\n node [ shape=plain] \n table [ label = <\n  <table border='1' cellborder='1'>\n   <tr><td>Atributo</td><td>Valor</td></tr>\n")
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				w := bufio.NewWriter(file)

				///////////SETEAR DATOS DEL SUPERBLOQUE
				n := bytes.Index(SB1.Name[:], []byte{0})
				fmt.Fprintf(w, "   <tr><td>Nombre</td><td>%v</td></tr>\n", string(SB1.Name[:n]))
				fmt.Fprintf(w, "   <tr><td>Total AVDs</td><td>%v</td></tr>\n", SB1.TotalAVDS)
				fmt.Fprintf(w, "   <tr><td>Total DDs</td><td>%v</td></tr>\n", SB1.TotalDDS)
				fmt.Fprintf(w, "   <tr><td>Total Inodos</td><td>%v</td></tr>\n", SB1.TotalInodos)
				fmt.Fprintf(w, "   <tr><td>Total Bloques</td><td>%v</td></tr>\n", SB1.TotalBloques)
				fmt.Fprintf(w, "   <tr><td>Total Bitacoras</td><td>%v</td></tr>\n", SB1.TotalBitacoras)
				fmt.Fprintf(w, "   <tr><td>Free AVDs</td><td>%v</td></tr>\n", SB1.FreeAVDS)
				fmt.Fprintf(w, "   <tr><td>Free DDs</td><td>%v</td></tr>\n", SB1.FreeDDS)
				fmt.Fprintf(w, "   <tr><td>Free Inodos</td><td>%v</td></tr>\n", SB1.FreeInodos)
				fmt.Fprintf(w, "   <tr><td>Free Bloques</td><td>%v</td></tr>\n", SB1.FreeBloques)
				fmt.Fprintf(w, "   <tr><td>Free Bitacoras</td><td>%v</td></tr>\n", SB1.FreeBitacoras)
				n = bytes.Index(SB1.DateCreacion[:], []byte{0})
				fmt.Fprintf(w, "   <tr><td>Fecha creación</td><td>%v</td></tr>\n", string(SB1.DateCreacion[:n]))
				n = bytes.Index(SB1.DateLastMount[:], []byte{0})
				fmt.Fprintf(w, "   <tr><td>Fecha Modificación</td><td>%v</td></tr>\n", string(SB1.DateLastMount[:n]))
				fmt.Fprintf(w, "   <tr><td>No. Montajes</td><td>%v</td></tr>\n", SB1.MontajesCount)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bitmap AVDs</td><td>%v</td></tr>\n", SB1.InicioBitmapAVDS)
				fmt.Fprintf(w, "   <tr><td>Apuntador AVDs</td><td>%v</td></tr>\n", SB1.InicioAVDS)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bitmap DDs</td><td>%v</td></tr>\n", SB1.InicioBitMapDDS)
				fmt.Fprintf(w, "   <tr><td>Apuntador DDs</td><td>%v</td></tr>\n", SB1.InicioDDS)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bitmap Inodos</td><td>%v</td></tr>\n", SB1.InicioBitmapInodos)
				fmt.Fprintf(w, "   <tr><td>Apuntador inodos</td><td>%v</td></tr>\n", SB1.InicioInodos)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bitmap bloques</td><td>%v</td></tr>\n", SB1.InicioBitmapBloques)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bloques</td><td>%v</td></tr>\n", SB1.InicioBloques)
				fmt.Fprintf(w, "   <tr><td>Apuntador Bitacoras</td><td>%v</td></tr>\n", SB1.InicioBitacora)
				fmt.Fprintf(w, "   <tr><td>Size struct AVD</td><td>%v</td></tr>\n", SB1.SizeAVD)
				fmt.Fprintf(w, "   <tr><td>Size struct DD</td><td>%v</td></tr>\n", SB1.SizeDD)
				fmt.Fprintf(w, "   <tr><td>Size struct Inodo</td><td>%v</td></tr>\n", SB1.SizeInodo)
				fmt.Fprintf(w, "   <tr><td>Size struct Bloque</td><td>%v</td></tr>\n", SB1.SizeBloque)
				fmt.Fprintf(w, "   <tr><td>Byte primer AVD libre</td><td>%v</td></tr>\n", SB1.FirstFreeAVD)
				fmt.Fprintf(w, "   <tr><td>Byte primer DD libre</td><td>%v</td></tr>\n", SB1.FirstFreeDD)
				fmt.Fprintf(w, "   <tr><td>Byte primer Inodo libre</td><td>%v</td></tr>\n", SB1.FirstFreeInodo)
				fmt.Fprintf(w, "   <tr><td>Byte primer Bloque libre</td><td>%v</td></tr>\n", SB1.FirstFreeBloque)
				fmt.Fprintf(w, "   <tr><td>Magic Number :D </td><td>%v</td></tr>\n", SB1.MagicNum)
				////////////////////

				w.Flush()

				_, err = file.WriteString("  </table>\n > ]\n}")
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Close()

				extT := "-T"

				switch strings.ToLower(extension) {
				case ".png":
					extT += "png"
				case ".pdf":
					extT += "pdf"
				case ".jpg":
					extT += "jpg"
				default:

				}

				if runtime.GOOS == "windows" {
					cmd := exec.Command("dot", extT, "-o", path, "codigo.dot") //Windows example, its tested
					//cmd.Stdout = os.Stdout
					cmd.Run()
				} else {
					cmd := exec.Command("dot", extT, "-o", path, "codigo.dot") //Linux example, its tested
					//cmd.Stdout = os.Stdout
					cmd.Run()
				}

			} else {
				color.Println("@{!r} La partición indicada no ha sido formateada.")
			}

		}

	} else {
		color.Println("@{!r}El reporte SB solo puede generar archivos con extensión .png, .jpg ó .pdf.")
	}

}

//ReporteBitmapAVD crea el reporte del Bitmap AVD
func ReporteBitmapAVD(path string, ruta string, id string) {
	extension := filepath.Ext(path)

	if strings.ToLower(extension) == ".txt" {

		NameAux, PathAux := GetDatosPart(id)

		if Existe, Indice := ExisteParticion(PathAux, NameAux); Existe {

			//LEER Y RECORRER EL MBR
			fileMBR, err2 := os.Open(PathAux)
			if err2 != nil { //validar que no sea nulo.
				panic(err2)
			}
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

				file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}
				// Change permissions Linux.
				err = os.Chmod(path, 0666)
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Truncate(0)
				file.Seek(0, 0)
				w := bufio.NewWriter(file)

				fileMBR.Seek(int64(SB1.InicioBitmapAVDS+1), 0)
				BitmapData := leerBytes(fileMBR, int(SB1.TotalAVDS))
				contador := 0
				for _, b := range BitmapData {
					if b == 1 {
						fmt.Fprint(w, "1	")
					} else {
						fmt.Fprint(w, "0	")
					}
					contador++
					if contador == 20 {
						fmt.Fprint(w, "\n")
						contador = 0
					}
				}

				w.Flush()
				file.Close()

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

				file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}
				// Change permissions Linux.
				err = os.Chmod(path, 0666)
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Truncate(0)
				file.Seek(0, 0)
				w := bufio.NewWriter(file)

				fileMBR.Seek(int64(SB1.InicioBitmapAVDS+1), 0)
				BitmapData := leerBytes(fileMBR, int(SB1.TotalAVDS))
				contador := 0
				for _, b := range BitmapData {
					if b == 1 {
						fmt.Fprint(w, "1	")
					} else {
						fmt.Fprint(w, "0	")
					}
					contador++
					if contador == 20 {
						fmt.Fprint(w, "\n")
						contador = 0
					}
				}

				w.Flush()
				file.Close()

			} else {
				color.Println("@{!r} La partición indicada no ha sido formateada.")
			}

			fileMBR.Close()
		}

	} else {
		color.Println("@{!r}El reporte BM_ARBDIR solo puede generar archivos con extensión .txt.")
	}
}

//ReporteBitmapDD crea el reporte del Bitmap AVD
func ReporteBitmapDD(path string, ruta string, id string) {
	extension := filepath.Ext(path)

	if strings.ToLower(extension) == ".txt" {

		NameAux, PathAux := GetDatosPart(id)

		if Existe, Indice := ExisteParticion(PathAux, NameAux); Existe {

			//LEER Y RECORRER EL MBR
			fileMBR, err2 := os.Open(PathAux)
			if err2 != nil { //validar que no sea nulo.
				panic(err2)
			}
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

				file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}
				// Change permissions Linux.
				err = os.Chmod(path, 0666)
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Truncate(0)
				file.Seek(0, 0)
				w := bufio.NewWriter(file)

				fileMBR.Seek(int64(SB1.InicioBitMapDDS+1), 0)
				BitmapData := leerBytes(fileMBR, int(SB1.TotalDDS))
				contador := 0
				for _, b := range BitmapData {
					if b == 1 {
						fmt.Fprint(w, "1	")
					} else {
						fmt.Fprint(w, "0	")
					}
					contador++
					if contador == 20 {
						fmt.Fprint(w, "\n")
						contador = 0
					}
				}

				w.Flush()
				file.Close()

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

				file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}
				// Change permissions Linux.
				err = os.Chmod(path, 0666)
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Truncate(0)
				file.Seek(0, 0)
				w := bufio.NewWriter(file)

				fileMBR.Seek(int64(SB1.InicioBitMapDDS+1), 0)
				BitmapData := leerBytes(fileMBR, int(SB1.TotalDDS))
				contador := 0
				for _, b := range BitmapData {
					if b == 1 {
						fmt.Fprint(w, "1	")
					} else {
						fmt.Fprint(w, "0	")
					}
					contador++
					if contador == 20 {
						fmt.Fprint(w, "\n")
						contador = 0
					}
				}

				w.Flush()
				file.Close()

			} else {
				color.Println("@{!r} La partición indicada no ha sido formateada.")
			}

			fileMBR.Close()
		}

	} else {
		color.Println("@{!r}El reporte BM_DETDIR solo puede generar archivos con extensión .txt.")
	}
}

//ReporteBitmapInode crea el reporte del Bitmap AVD
func ReporteBitmapInode(path string, ruta string, id string) {
	extension := filepath.Ext(path)

	if strings.ToLower(extension) == ".txt" {

		NameAux, PathAux := GetDatosPart(id)

		if Existe, Indice := ExisteParticion(PathAux, NameAux); Existe {

			//LEER Y RECORRER EL MBR
			fileMBR, err2 := os.Open(PathAux)
			if err2 != nil { //validar que no sea nulo.
				panic(err2)
			}
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

				file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}
				// Change permissions Linux.
				err = os.Chmod(path, 0666)
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Truncate(0)
				file.Seek(0, 0)
				w := bufio.NewWriter(file)

				fileMBR.Seek(int64(SB1.InicioBitmapInodos+1), 0)
				BitmapData := leerBytes(fileMBR, int(SB1.TotalInodos))
				contador := 0
				for _, b := range BitmapData {
					if b == 1 {
						fmt.Fprint(w, "1	")
					} else {
						fmt.Fprint(w, "0	")
					}
					contador++
					if contador == 20 {
						fmt.Fprint(w, "\n")
						contador = 0
					}
				}

				w.Flush()
				file.Close()

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

				file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}
				// Change permissions Linux.
				err = os.Chmod(path, 0666)
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Truncate(0)
				file.Seek(0, 0)
				w := bufio.NewWriter(file)

				fileMBR.Seek(int64(SB1.InicioBitmapInodos+1), 0)
				BitmapData := leerBytes(fileMBR, int(SB1.TotalInodos))
				contador := 0
				for _, b := range BitmapData {
					if b == 1 {
						fmt.Fprint(w, "1	")
					} else {
						fmt.Fprint(w, "0	")
					}
					contador++
					if contador == 20 {
						fmt.Fprint(w, "\n")
						contador = 0
					}
				}

				w.Flush()
				file.Close()

			} else {
				color.Println("@{!r} La partición indicada no ha sido formateada.")
			}

			fileMBR.Close()
		}

	} else {
		color.Println("@{!r}El reporte BM_INODE solo puede generar archivos con extensión .txt.")
	}
}

//ReporteBitmapBloque crea el reporte del Bitmap AVD
func ReporteBitmapBloque(path string, ruta string, id string) {
	extension := filepath.Ext(path)

	if strings.ToLower(extension) == ".txt" {

		NameAux, PathAux := GetDatosPart(id)

		if Existe, Indice := ExisteParticion(PathAux, NameAux); Existe {

			//LEER Y RECORRER EL MBR
			fileMBR, err2 := os.Open(PathAux)
			if err2 != nil { //validar que no sea nulo.
				panic(err2)
			}
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

				file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}
				// Change permissions Linux.
				err = os.Chmod(path, 0666)
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Truncate(0)
				file.Seek(0, 0)
				w := bufio.NewWriter(file)

				fileMBR.Seek(int64(SB1.InicioBitmapBloques+1), 0)
				BitmapData := leerBytes(fileMBR, int(SB1.TotalBloques))
				contador := 0
				for _, b := range BitmapData {
					if b == 1 {
						fmt.Fprint(w, "1	")
					} else {
						fmt.Fprint(w, "0	")
					}
					contador++
					if contador == 20 {
						fmt.Fprint(w, "\n")
						contador = 0
					}
				}

				w.Flush()
				file.Close()

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

				file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666) //Crea un nuevo archivo
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}
				// Change permissions Linux.
				err = os.Chmod(path, 0666)
				if err != nil {
					fmt.Println(err)
					file.Close()
					return
				}

				file.Truncate(0)
				file.Seek(0, 0)
				w := bufio.NewWriter(file)

				fileMBR.Seek(int64(SB1.InicioBitmapBloques+1), 0)
				BitmapData := leerBytes(fileMBR, int(SB1.TotalBloques))
				contador := 0
				for _, b := range BitmapData {
					if b == 1 {
						fmt.Fprint(w, "1	")
					} else {
						fmt.Fprint(w, "0	")
					}
					contador++
					if contador == 20 {
						fmt.Fprint(w, "\n")
						contador = 0
					}
				}

				w.Flush()
				file.Close()

			} else {
				color.Println("@{!r} La partición indicada no ha sido formateada.")
			}

			fileMBR.Close()
		}

	} else {
		color.Println("@{!r}El reporte BM_BLOCK solo puede generar archivos con extensión .txt.")
	}
}
