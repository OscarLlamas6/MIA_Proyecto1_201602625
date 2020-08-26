package estructuras

//DD struct
type DD struct {
	DDFiles     [5]DDFile
	ApuntadorDD int64
}

//DDFile struct
type DDFile struct {
	Name              [20]byte
	ApuntadorInodo    int64
	FechaCreacion     [20]byte
	FechaModificacion [20]byte
}
