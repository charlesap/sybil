package lodge

import (
	"fmt"
)

type Lodge struct{
	Htop uint32
	Ptop uint32
	Rtop uint32
	X    uint32
	Fqdn [126] byte
	Subn [126] byte
	Bdev [256] byte
}

type Hash [28] byte

type Sign [64] byte

type Span struct{ // 56 bytes 
	Hsh Hash
	Lnk [7] uint32
}

type Mesg struct{ // 256 bytes
	H    Span
	Time [7] byte
	Op   byte
	P    Span
	Acc1 uint32
	Acc2 uint32
	R    Span
	Acc3 uint32
	Acc4 uint32
	S    Sign
}

type Body struct{
	Ones Hash  // all zeros = empty all ones = text otherwise Mesg
	Text [228] byte
}

func Emit(name string) {

	fmt.Println( name )
}

func Format() { // write 28 zeros at every 256th location to empty all blocks

}

func Retrieve() { // hash and bounce (hash again, look again) until match or zeros

}

func Place() { // hash and bounce until match (e.g. exists) or zeros, in which write.

}

func Readin() { // iteratively load from text file and place in the lodge

}

func Writeout() { // iteratively retrieve from the lodge and write to text file

}

