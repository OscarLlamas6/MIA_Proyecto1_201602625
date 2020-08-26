package estructuras

//AVD struct
type AVD struct {
	FechaCreacion [20]byte
	NombreDir     [20]byte
	ApuntadorSubs [6]int64
	ApuntadorAVD  int64
	ApuntadorDD   int64
	Proper        [20]byte
}
