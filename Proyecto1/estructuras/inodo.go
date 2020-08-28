package estructuras

//Inodo struct
type Inodo struct {
	NumeroInodo        int32
	FileSize           int32
	NumeroBloques      int32
	ApuntadoresBloques [4]int32
	ApuntadorIndirecto int32
	Proper             [20]byte
}

//BloqueDatos struct
type BloqueDatos struct {
	Data [25]byte
}
