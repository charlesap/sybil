package lodge

import (
	"fmt"
)

type Span struct{
	Hsh [28] byte
	Loc uint32
	Lnk [8] uint32
}

type Sign [64] byte

type Mdta struct{
	Dnt uint64
	Opr uint32
	Acc [5] uint32
}

func Emit(name string) {

	fmt.Println( name )
}

