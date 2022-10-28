package lodge

import (
	"fmt"
)

type Lodge struct{
	Fqdn [126] byte
	Subn [126] byte
	Bdev [256] byte
}

type Hash [28] byte

type Sign [64] byte

type Span7 struct{ // 56 bytes 
	Hsh Hash
	Lnk [7] uint32
}

type Span6 struct{ // 56 bytes 
	Hsh Hash
	Top uint32
	Lnk [6] uint32
}

type Mesg struct{ // 256 bytes
	Op   byte
	Time [7] byte
	H    Span7
	Pd   byte  // principal signer depth, 1 = self
	Salt [3] byte
	Fld1 byte
	Acc1 [3] byte
	P    Span6
	Fld2 byte
	Acc2 [3] byte
	Fld3 byte
	Acc3 [3] byte
	R    Span6
	S    Sign
}

type Body struct{ // 252 possble bytes of text
	Pad1 byte  // all zeros = empty all ones = text otherwise Mesg
	Text1 [63] byte
	Pad2 byte  // all zeros = empty all ones = text otherwise Mesg
	Text2 [63] byte
	Pad3 byte  // all zeros = empty all ones = text otherwise Mesg
	Text3 [63] byte
	Len  byte  // of this text segment not including null, 255 = continue next block
	Text4 [63] byte
	                // all utf8
	                // if Len is 254 then text is null terminated special reference (image? movie?)
}

func Emit(name string) {

	fmt.Println( name )
}

func Format() { // write zeros to empty all blocks

}

func Retrieve() { // hash and bounce (hash again, look again) until match or zeros
	// a (very unlikely) hash resulting in all zeros is simply hashed again.
	// a match with Op or Fld1 or Fld3  being zero or 255 is not a match

}

func Place() { // hash and bounce until match (e.g. exists) or zeros, in which write.

}

func Validate() { // use ed25516 public key to recalculate signature of canonically structured message

}

func Readone() { // bring canonically structured message into memory from file

}

func Readin() { // iteratively load from text file and place in the lodge

}

func Writeone() {

}

func Writeout() { // iteratively retrieve from the lodge and write to text file, identities first to help Validate on Readin

}

func Init() { // since nothing will hash to zero, the zero block is the root block for H, P, and R.

}
