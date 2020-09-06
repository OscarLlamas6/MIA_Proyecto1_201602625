package funciones

import (
	"Proyecto1/estructuras"
	"bytes"
	"fmt"
)

//GenerarAVD devuelve un AVD seteado en formato string
func GenerarAVD(NoAVD int, AVDaux *estructuras.AVD) string {

	n := bytes.Index(AVDaux.NombreDir[:], []byte{0})
	Directorio := string(AVDaux.NombreDir[:n])

	n = bytes.Index(AVDaux.FechaCreacion[:], []byte{0})
	Fecha := string(AVDaux.FechaCreacion[:n])

	n = bytes.Index(AVDaux.Proper[:], []byte{0})
	Propietario := string(AVDaux.Proper[:n])

	n = bytes.Index(AVDaux.Grupo[:], []byte{0})
	Grupo := string(AVDaux.Grupo[:n])

	cadena := fmt.Sprintf(`AVD%v [label=<
	<TABLE BORDER="1"  cellpadding="2"   CELLBORDER="1" CELLSPACING="4" BGCOLOR="blue4" color = 'black'>            
	   <TR> 
		   <TD bgcolor='purple' colspan="2"><font color='white' point-size='13'>AVD: %s</font></TD>
	   </TR>
	   <TR> 
		   <TD bgcolor="slateblue2" >Fecha creación</TD>
		   <TD bgcolor='slateblue2' > %s </TD>
	   </TR>
	   <TR>
		   <TD bgcolor='slateblue2' >Propietario</TD>
		   <TD bgcolor='slateblue2' > %s </TD>
	   </TR>
	   <TR>
		   <TD bgcolor='slateblue2' >Grupo</TD>
		   <TD bgcolor='slateblue2' > %s </TD>
	   </TR>
	   <TR>
		   <TD bgcolor='slateblue2' >Permisos</TD>
		   <TD bgcolor='slateblue2' > %v%v%v </TD>
	   </TR>
	   <TR>
		   <TD  bgcolor='green1' >Subdirectorios[0]</TD>
		   <TD  bgcolor='green1' PORT="0"> %v </TD>
	   </TR>
	   <TR>
		   <TD  bgcolor='green1' >Subdirectorios[1]</TD>
		   <TD  bgcolor='green1' PORT="1"> %v </TD>
	   </TR>
	   <TR>
		   <TD  bgcolor='green1' >Subdirectorios[2]</TD>
		   <TD  bgcolor='green1' PORT="2"> %v </TD>
	   </TR>
	   <TR>
		   <TD  bgcolor='green1' >Subdirectorios[3]</TD>
		   <TD  bgcolor='green1' PORT="3"> %v </TD>
	   </TR>
	   <TR>
		   <TD  bgcolor='green1' >Subdirectorios[4]</TD>
		   <TD  bgcolor='green1' PORT="4"> %v </TD>
	   </TR>
	   <TR>
		   <TD  bgcolor='green1' >Subdirectorios[5]</TD>
		   <TD  bgcolor='green1' PORT="5"> %v </TD>
	   </TR>
	   <TR>
		   <TD  bgcolor='hotpink' >ApuntadorDD</TD>
		   <TD  bgcolor='hotpink' PORT="6"> %v </TD>
	   </TR>
	   <TR>
		   <TD  bgcolor='orange' >ApuntadorAVD</TD>
		   <TD  bgcolor='orange' PORT="7"> %v</TD>
	   </TR>
   </TABLE>
	>];

	`, NoAVD, Directorio, Fecha, Propietario, Grupo, AVDaux.PermisoU, AVDaux.PermisoG, AVDaux.PermisoO, AVDaux.ApuntadorSubs[0], AVDaux.ApuntadorSubs[1], AVDaux.ApuntadorSubs[2], AVDaux.ApuntadorSubs[3], AVDaux.ApuntadorSubs[4], AVDaux.ApuntadorSubs[5], AVDaux.ApuntadorDD, AVDaux.ApuntadorAVD)

	return cadena
}

//GenerarDD devuelve un DD seteado en formato string
func GenerarDD(NoDD int, DDaux *estructuras.DD, carpeta string) string {

	// i = 0

	n := bytes.Index(DDaux.DDFiles[0].Name[:], []byte{0})
	Nombre0 := string(DDaux.DDFiles[0].Name[:n])

	n = bytes.Index(DDaux.DDFiles[0].FechaCreacion[:], []byte{0})
	Fc0 := string(DDaux.DDFiles[0].FechaCreacion[:n])

	n = bytes.Index(DDaux.DDFiles[0].FechaModificacion[:], []byte{0})
	Fm0 := string(DDaux.DDFiles[0].FechaModificacion[:n])

	// i = 1

	n = bytes.Index(DDaux.DDFiles[1].Name[:], []byte{0})
	Nombre1 := string(DDaux.DDFiles[1].Name[:n])

	n = bytes.Index(DDaux.DDFiles[1].FechaCreacion[:], []byte{0})
	Fc1 := string(DDaux.DDFiles[1].FechaCreacion[:n])

	n = bytes.Index(DDaux.DDFiles[1].FechaModificacion[:], []byte{0})
	Fm1 := string(DDaux.DDFiles[1].FechaModificacion[:n])

	// i = 2

	n = bytes.Index(DDaux.DDFiles[2].Name[:], []byte{0})
	Nombre2 := string(DDaux.DDFiles[2].Name[:n])

	n = bytes.Index(DDaux.DDFiles[2].FechaCreacion[:], []byte{0})
	Fc2 := string(DDaux.DDFiles[2].FechaCreacion[:n])

	n = bytes.Index(DDaux.DDFiles[2].FechaModificacion[:], []byte{0})
	Fm2 := string(DDaux.DDFiles[2].FechaModificacion[:n])

	// i = 3

	n = bytes.Index(DDaux.DDFiles[3].Name[:], []byte{0})
	Nombre3 := string(DDaux.DDFiles[3].Name[:n])

	n = bytes.Index(DDaux.DDFiles[3].FechaCreacion[:], []byte{0})
	Fc3 := string(DDaux.DDFiles[3].FechaCreacion[:n])

	n = bytes.Index(DDaux.DDFiles[3].FechaModificacion[:], []byte{0})
	Fm3 := string(DDaux.DDFiles[3].FechaModificacion[:n])

	// i = 4

	n = bytes.Index(DDaux.DDFiles[4].Name[:], []byte{0})
	Nombre4 := string(DDaux.DDFiles[4].Name[:n])

	n = bytes.Index(DDaux.DDFiles[4].FechaCreacion[:], []byte{0})
	Fc4 := string(DDaux.DDFiles[4].FechaCreacion[:n])

	n = bytes.Index(DDaux.DDFiles[4].FechaModificacion[:], []byte{0})
	Fm4 := string(DDaux.DDFiles[4].FechaModificacion[:n])

	cadena := fmt.Sprintf(`DD%v [label=<
	<TABLE BORDER="1"  cellpadding="2"   CELLBORDER="1" CELLSPACING="4" BGCOLOR="blue4" color = 'black'>            
	   <TR> 
		   <TD bgcolor='purple' colspan="2"><font color='white' point-size='13'>DD: %s</font></TD>
	   </TR>
	   <TR>
		   <TD bgcolor='khaki' >[0].Nombre</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
	   <TR>
		   <TD bgcolor='khaki' >[0].FechaCreacion</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
		<TR>
		   <TD bgcolor='khaki' >[0].FechaModificación</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
		<TR>
		   <TD bgcolor='coral' >[0].ApuntadorInodo</TD>
		   <TD bgcolor='coral' PORT="0" > %v </TD>
	   </TR>
		<TR>
		   <TD bgcolor='khaki' >[1].Nombre</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
	   <TR>
		   <TD bgcolor='khaki' >[1].FechaCreacion</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
		<TR>
		   <TD bgcolor='khaki' >[1].FechaModificación</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
		<TR>
		   <TD bgcolor='coral' >[1].ApuntadorInodo</TD>
		   <TD bgcolor='coral' PORT="1" > %v </TD>
	   </TR>
		<TR>
		   <TD bgcolor='khaki' >[2].Nombre</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
	   <TR>
		   <TD bgcolor='khaki' >[2].FechaCreacion</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
		<TR>
		   <TD bgcolor='khaki' >[2].FechaModificación</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
		<TR>
		   <TD bgcolor='coral' >[2].ApuntadorInodo</TD>
		   <TD bgcolor='coral' PORT="2" > %v </TD>
	   </TR>
		<TR>
		   <TD bgcolor='khaki' >[3].Nombre</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
	   <TR>
		   <TD bgcolor='khaki' >[3].FechaCreacion</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
		<TR>
		   <TD bgcolor='khaki' >[3].FechaModificación</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
		<TR>
		   <TD bgcolor='coral' >[3].ApuntadorInodo</TD>
		   <TD bgcolor='coral' PORT="3" > %v </TD>
	   </TR>
		<TR>
		   <TD bgcolor='khaki' >[4].Nombre</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
	   <TR>
		   <TD bgcolor='khaki' >[4].FechaCreacion</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
		<TR>
		   <TD bgcolor='khaki' >[4].FechaModificación</TD>
		   <TD bgcolor='khaki' > %s </TD>
	   </TR>
		<TR>
		   <TD bgcolor='coral' >[4].ApuntadorInodo</TD>
		   <TD bgcolor='coral' PORT="4" > %v </TD>
	   </TR>
	   <TR>
		   <TD  bgcolor='springgreen' >ApuntadorDD</TD>
		   <TD  bgcolor='springgreen' PORT="5"> %v </TD>
	   </TR>

   </TABLE>
	>];

	`, NoDD, carpeta, Nombre0, Fc0, Fm0, DDaux.DDFiles[0].ApuntadorInodo, Nombre1, Fc1, Fm1, DDaux.DDFiles[1].ApuntadorInodo, Nombre2, Fc2, Fm2, DDaux.DDFiles[2].ApuntadorInodo, Nombre3, Fc3, Fm3, DDaux.DDFiles[3].ApuntadorInodo, Nombre4, Fc4, Fm4, DDaux.DDFiles[4].ApuntadorInodo, DDaux.ApuntadorDD)

	return cadena

}

//GenerarInodo devuelve un Inodo seteado en formato string
func GenerarInodo(NoInodo int, InodoAux *estructuras.Inodo) string {

	n := bytes.Index(InodoAux.Proper[:], []byte{0})
	Propietario := string(InodoAux.Proper[:n])

	n = bytes.Index(InodoAux.Grupo[:], []byte{0})
	Grupo := string(InodoAux.Grupo[:n])

	cadena := fmt.Sprintf(`Inodo%v [label=<
	<TABLE BORDER="1"  cellpadding="2"   CELLBORDER="1" CELLSPACING="4" BGCOLOR="blue4" color = 'black'>            
	   <TR>
	   <TD bgcolor='purple' colspan="2"><font color='white' point-size='13'>Inodo %v</font></TD>
	   </TR>
	   <TR> 
		   <TD bgcolor='yellow' >Propietario</TD>
		   <TD bgcolor='yellow' > %s </TD>
	   </TR>
	   <TR> 
		   <TD bgcolor='yellow' >Grupo</TD>
		   <TD bgcolor='yellow' > %s </TD>
	   </TR>
	   <TR> 
		   <TD bgcolor='yellow' >File Size</TD>
		   <TD bgcolor='yellow' > %v </TD>
	   </TR>
	   <TR> 
		   <TD bgcolor='yellow' >Numero de bloques</TD>
		   <TD bgcolor='yellow' > %v </TD>
	   </TR>
	   <TR> 
		   <TD bgcolor='yellow' >Permisos</TD>
		   <TD bgcolor='yellow' > %v%v%v </TD>
	   </TR>
	   <TR> 
		   <TD bgcolor='mistyrose' >Bloques[0]</TD>
		   <TD bgcolor='mistyrose' PORT="0" > %v </TD>
	   </TR>
	   <TR> 
		   <TD bgcolor='mistyrose' >Bloques[1]</TD>
		   <TD bgcolor='mistyrose' PORT="1" > %v </TD>
	   </TR>
	   <TR> 
		   <TD bgcolor='mistyrose' >Bloques[2]</TD>
		   <TD bgcolor='mistyrose' PORT="2" > %v </TD>
	   </TR>
	   <TR> 
		   <TD bgcolor='mistyrose' >Bloques[3]</TD>
		   <TD bgcolor='mistyrose' PORT="3" > %v </TD>
	   </TR>
	   <TR> 
		   <TD bgcolor='lime' >ApuntadorInodo</TD>
		   <TD bgcolor='lime' PORT="4" > %v </TD>
	   </TR>

   	</TABLE>
   >];
   
	`, NoInodo, InodoAux.NumeroInodo, Propietario, Grupo, InodoAux.FileSize, InodoAux.NumeroBloques, InodoAux.PermisoU, InodoAux.PermisoG, InodoAux.PermisoO, InodoAux.ApuntadoresBloques[0], InodoAux.ApuntadoresBloques[1], InodoAux.ApuntadoresBloques[2], InodoAux.ApuntadoresBloques[3], InodoAux.ApuntadorIndirecto)

	return cadena
}

//GenerarBloque devuelve un Bloque de datos seteado en formato string
func GenerarBloque(NoBloque, Bloqueaux *estructuras.BloqueDatos) string {

	n := bytes.Index(Bloqueaux.Data[:], []byte{0})
	contenido := string(Bloqueaux.Data[:n])

	cadena := fmt.Sprintf(`Bloque%v [label=<
	<table border="2" cellborder="0" cellspacing="1" bgcolor="lightsalmon" color="black">
		<tr> 
			<TD align ="center"><font color="white" >Bloque de Datos</font></TD> 
		</tr>
		<tr>
			<TD align="left"> %s </TD>
		</tr>
	</table>
	>];
	
	`, NoBloque, contenido)

	return cadena

}
