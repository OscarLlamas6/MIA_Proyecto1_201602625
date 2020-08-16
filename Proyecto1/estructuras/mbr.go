package estructuras

//MBR struct
type MBR struct {
	msize       uint32
	mdate       [20]byte
	msignature  uint32
	mpartitions [4]Partition
}

//Partition struct
type Partition struct {
	pstart  uint32
	psize   uint32
	pname   [16]byte
	pstatus byte
	ptype   byte
	pfit    byte
}
