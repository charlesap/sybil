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

type Kndx  uint32 // index of a Knod in block device of Knods

type Cksm  uint32 // checksum of checksums of nodes held by this node or
                  // checksum of checksums of nodes referring to this node

type Span7 struct{ // 56 bytes 
	Hsh Hash
	Lnk [7] Kndx
}

type Span6 struct{ // 56 bytes 
	Hsh Hash
	Top Kndx
	Lnk [6] Kndx
}

type Knod struct{ // 256 bytes // knowledge node, Tnod is a text representation
	Op   byte  // op except 0 = hash slot free, 255 = hash slot available due to allocation bounce on content size
	Time [7] byte
	H    Span7
	Idpt byte  // principal signer depth, 1 = self
	Salt [3] byte
	Iloc Kndx  // principal signer quick link
	P    Span6
	Pchk Cksm  // checksum of Pchk of nodes held... if two lodges differ on this, synchronization is indicated
	Rchk Cksm  // checksum of Rchk of nodes held... if two lodges differ on this, synchronization is indicated
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
