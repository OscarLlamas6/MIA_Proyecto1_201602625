package estructuras

//Inodo struct
type Inodo struct {
	NumeroInodo        int32
	FileSize           int32
	NumeroBloques      int32
	ApuntadoresBloques [4]int64
	ApuntadorIndirecto int64
	Proper             [20]byte
}

//BloqueDatos struct
type BloqueDatos struct {
	Data [25]byte
}
