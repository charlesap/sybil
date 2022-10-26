package lodge

import (
	"fmt"
)

type Lodge struct{
	Fqdn [120] byte
	Subn [120] byte
	Path [120] byte
	Name [120] byte
	Size uint64
	Htop uint32
	Ptop uint32
	Rtop uint32
	X    uint32
	Y    uint64
}

type Hash [28] byte

type Sign [64] byte

type Span struct{ // 64 bytes * 3 spans per message
	Hsh Hash
	Loc uint32
	Lnk [8] uint32
}

type Mdta struct{ // 32 bytes
	Dnt uint64
	Opr uint32
	Acc [5] uint32
}

type Tiny [64] byte
type Smol [128] byte
type Medm [256] byte
type Larg [512] byte
type Huge [1024] byte

func Emit(name string) {

	fmt.Println( name )
}

