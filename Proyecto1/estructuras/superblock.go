package estructuras

//Superblock struct
type Superblock struct {
	Name                   [16]byte
	TotalAVDS              int32
	TotalDDS               int32
	TotalInodos            int32
	TotalBloques           int32
	FreeAVD                int32
	FreeDD                 int32
	FreeInodos             int32
	FreeBloques            int32
	DateCreacion           [20]byte
	DateLastMount          [20]byte
	MontajesCount          int64
	APuntadorBitmapAVD     int64
	ApuntadorAVD           int64
	ApuntadorBitMapDds     int64
	ApuntadorDDS           int64
	ApuntadorBitmapInodos  int64
	ApuntadorInodos        int64
	ApuntadorBitmapBloques int64
	ApuntadorBloques       int64
	ApuntadorBitacora      int64
	SizeAVD                int32
	SizeDD                 int32
	sizeInodo              int32
	SizeBloque             int32
	FirstFreeAVD           int64
	FirstFreeDDS           int64
	FirstFreeInodo         int64
	FirstFreeBloque        int64
	MagicNum               int32
}
