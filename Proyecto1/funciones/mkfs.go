package funciones

import (
	"Proyecto1/estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"github.com/doun/terminal/color"
)

//EjecutarMkfs inicia el formateo de una particion
func EjecutarMkfs(id string, tipo string, add string, unit string) {

	if id != "" {

		if IDYaRegistrado(id) {

			color.Println("@{!c}Listo para montar sistema de archivos")

			NameAux, PathAux := GetDatosPart(id)

			if Existe, Indice := ExisteParticion(PathAux, NameAux); Existe { //SIGNIFICA QUE ES PRIMARIA, NO PUEDE SER EXTENDIDA PORQUE CUANDO EN LA FUNCION MOUNT SE VERIFICAR QUE NO SEA EXTENDIDA ANTES DE MONTAR.

				PartStart, PartSize := GetStartAndSize(PathAux, Indice)

				color.Printf("@{!m}La particion inicia en el byte: %d\n", PartStart)
				color.Printf("@{!m}La particion tiene un size de: %d\n", PartSize)

				Formatear(PartStart, PartSize, tipo, PathAux)

			} else if ExisteL, _ := ExisteParticionLogica(PathAux, NameAux); ExisteL { //SIGNIFICA QUE ES LOGICA (NO ES REQUISITO)

			}

		} else {
			color.Printf("@{!r}No hay ninguna partición montada con el id: %v\n", id)
		}

	} else {
		color.Println("@{!r}Falta parámetro -id, obligatorio para ejecutar la función Mkfs.")
	}

}

//Formatear ejecuta el formateo Full O Fast según se requiera
func Formatear(PartStart int, PartSize int, tipo string, path string) {

	sizeSuperbloque := int32(unsafe.Sizeof(estructuras.Superblock{}))
	startPart := int32(PartStart)
	sizePart := int32(PartSize)

	if (sizePart - (2 * sizeSuperbloque)) > 0 {

		//Obteniendo sizes de cada structs

		sizeAVD := int32(unsafe.Sizeof(estructuras.AVD{}))
		sizeDD := int32(unsafe.Sizeof(estructuras.DD{}))
		sizeInodo := int32(unsafe.Sizeof(estructuras.Inodo{}))
		sizeBloque := int32(unsafe.Sizeof(estructuras.BloqueDatos{}))
		sizeBitacora := int32(unsafe.Sizeof(estructuras.Bitacora{}))

		//Calculando número de estructuras

		NumEstructuras := (sizePart - (2 * sizeSuperbloque)) / (27 + sizeAVD + sizeDD + (5*sizeInodo + (20 * sizeBloque) + sizeBitacora))

		cantidadAVDS := NumEstructuras
		cantidadDDS := NumEstructuras
		cantidadInodos := int32(5 * NumEstructuras)
		cantidadBloques := int32(4 * cantidadInodos)
		cantidadBitacoras := NumEstructuras

		//Seteando Superbloque

		sb := estructuras.Superblock{}
		var chars [16]byte
		VirtualName := filepath.Base(path)
		copy(chars[:], VirtualName)
		copy(sb.Name[:], chars[:])
		sb.TotalAVDS = cantidadAVDS
		sb.TotalDDS = cantidadDDS
		sb.TotalInodos = cantidadInodos
		sb.TotalBloques = cantidadBloques
		sb.TotalBitacoras = cantidadBitacoras
		sb.FreeAVDS = cantidadAVDS - 1
		sb.FreeDDS = cantidadDDS - 1
		sb.FreeInodos = cantidadInodos
		sb.FreeBloques = cantidadBloques
		sb.FreeBitacoras = cantidadBitacoras
		t := time.Now()
		var charsDate [20]byte
		cadena := t.Format("2006-01-02 15:04:05")
		copy(charsDate[:], cadena)
		copy(sb.DateCreacion[:], charsDate[:])
		copy(sb.DateLastMount[:], charsDate[:])
		sb.MontajesCount = 1
		sb.InicioBitmapAVDS = startPart + sizeSuperbloque
		sb.InicioAVDS = sb.InicioBitmapAVDS + cantidadAVDS
		sb.InicioBitMapDDS = sb.InicioAVDS + (sizeAVD * cantidadAVDS)
		sb.InicioDDS = sb.InicioBitMapDDS + cantidadDDS
		sb.InicioBitmapInodos = sb.InicioDDS + (sizeDD * cantidadDDS)
		sb.InicioInodos = sb.InicioBitmapInodos + cantidadInodos
		sb.InicioBitmapBloques = sb.InicioInodos + (sizeInodo * cantidadInodos)
		sb.InicioBloques = sb.InicioBitmapBloques + cantidadBloques
		sb.InicioBitacora = sb.InicioBloques + (sizeBloque * cantidadBloques)
		sb.SizeAVD = sizeAVD
		sb.SizeDD = sizeDD
		sb.SizeInodo = sizeInodo
		sb.SizeBloque = sizeBloque
		sb.SizeBitacora = sizeBitacora
		sb.FirstFreeAVD = sb.InicioAVDS + sb.SizeAVD //Le sumamos un sizeAVD porque vamos a crear la carpeta "/"
		sb.FirstFreeDD = sb.InicioDDS + sb.SizeDD    //Le sumamos un sizeDD porque vamos a crear el DD de la carpeta "/"
		sb.FirstFreeInodo = sb.InicioInodos
		sb.FirstFreeBloque = sb.InicioBloques
		sb.MagicNum = 201602625

		file, err := os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			fmt.Println(err)
			file.Close()
		}

		file.Seek(int64(PartStart+1), 0)
		if strings.ToLower(tipo) == "full" {
			data := make([]byte, sizePart)
			file.Write(data)
			file.Seek(int64(PartStart+1), 0)
		}
		//Escribiendo el Superbloque
		sb1 := &sb
		var binario1 bytes.Buffer
		binary.Write(&binario1, binary.BigEndian, sb1)
		escribirBytes(file, binario1.Bytes())

		//Escribiendo el Bitmap de AVDS
		file.Seek(int64(sb.InicioBitmapAVDS+1), 0)
		data := make([]byte, cantidadAVDS)
		file.Write(data)

		//Escribiendo el Bitmap de DDS
		file.Seek(int64(sb.InicioBitMapDDS+1), 0)
		data = make([]byte, cantidadDDS)
		file.Write(data)

		//Escribiendo de Bitmap de Inodos
		file.Seek(int64(sb.InicioBitmapInodos+1), 0)
		data = make([]byte, cantidadInodos)
		file.Write(data)

		//Escribiendo el Bitmap de bloques
		file.Seek(int64(sb.InicioBitmapBloques+1), 0)
		data = make([]byte, cantidadBloques)
		file.Write(data)

		//Escribir el Backup del Superbloque
		file.Seek(int64((sb.InicioBitacora+(sizeBitacora*cantidadBitacoras))+1), 0)
		sb2 := &sb
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, sb2)
		escribirBytes(file, binario2.Bytes())

		//Creando folder / y users.txt
		//Bitmap de AVD (la primera posición es la carpeta "/")
		//Escribiendo un 1 en la primera posicion del bitmap
		file.Seek(int64(sb.InicioBitmapAVDS+1), 0)
		data = []byte{0x01}
		file.Write(data)
		//Seteando valores de la carpeta root "/
		//Creamos nueva estructura AVD, la cual será escrita en su posición correspondiente
		AVDaux := estructuras.AVD{}
		t = time.Now()
		cadena = t.Format("2006-01-02 15:04:05")
		copy(charsDate[:], cadena)
		//Seteando fecha de creacion
		copy(AVDaux.FechaCreacion[:], charsDate[:])
		var ArrayNombre [20]byte
		nombreDir := "/"
		copy(ArrayNombre[:], nombreDir)
		//Seteando nombre del directorio "/"
		copy(AVDaux.NombreDir[:], ArrayNombre[:])
		//La primera estructura AVD apuntará al primer Detalle de Directorio
		//Seteando el apuntador de su DD, en este caso es InicioDDS
		//al ser el primer DD que se usar
		AVDaux.ApuntadorDD = sb.InicioDDS
		nombrePropietario := "root"
		copy(ArrayNombre[:], nombrePropietario)
		//Seteando nombre del propietario, en este caso la raiz pertenece al id "root"
		copy(AVDaux.Proper[:], ArrayNombre[:])
		//APuntadorAVD y los 6 apuntadores a subdirectorios no se setean en este momento
		//se hará conforme se vayan creando subdirectorios :)

		//Ahora toca escribir el struct AVD en su posición correspondiente
		file.Seek(int64(sb.InicioAVDS+1), 0)
		avdp := &AVDaux
		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, avdp)
		escribirBytes(file, binario3.Bytes())

		//Bitmap de DD (la primera posición es el DD de la carpeta "/")
		//Escribiendo un 1 en la primera posicion del bitmap
		file.Seek(int64(sb.InicioBitMapDDS+1), 0)
		data = []byte{0x01}
		file.Write(data)

		//Creamos una estructura DD para la carpeta "/"
		DDaux := estructuras.DD{}

		//No le seteamos nada, porque no se han creado archivos
		//Solo es el struct vació correspondiente para la carpeta "/"
		//Se le irán seteando valores conforme se creen archivos en este directorio :)

		//Ahora toca escribir el struct DD en su posición correspondiente
		file.Seek(int64(sb.InicioDDS+1), 0)
		ddp := &DDaux
		var binario4 bytes.Buffer
		binary.Write(&binario4, binary.BigEndian, ddp)
		escribirBytes(file, binario4.Bytes())

		file.Close()

		color.Println("@{!c}Sistema de archivos LWH instalado correctamente.")

	} else {
		color.Println("@{!r}El tamaño de la partición es insuficiente para montar el sistema de archivos LWH.")
	}
}
