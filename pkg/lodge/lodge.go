package lodge

import (
    b64 "encoding/base64"
    "fmt"
)

type Lodge struct{
	Fqdn [126] byte
	Subn [126] byte
	Bdev [256] byte
}

type Kdate [3]byte

func (k Kdate) String() string {
	return fmt.Sprintf("YYMMDD")
}

type Ktime [3]byte

func (k Ktime) String() string {
	return fmt.Sprintf("hhmmss")
}

type Hash [28] byte

func (h Hash) String() string {
	return fmt.Sprintf(b64.StdEncoding.EncodeToString(h[:]))
}

type Sign [64] byte

func (s Sign) String() string {
	return fmt.Sprintf(b64.StdEncoding.EncodeToString(s[:]))
}

type Kndx  uint32 // index of a Knod in block device of Knods

type Cksm  uint32 // checksum of checksums of nodes held by this node or
                  // checksum of checksums of nodes referring to this node

type Twig7 [7] Kndx

func (t Twig7) String() string {
	return fmt.Sprintf("%0.8x",t[:])
}

type Twig6 [6] Kndx

func (t Twig6) String() string {
	return fmt.Sprintf("%0.8x",t[:])
}

type Span7 struct{ // 56 bytes 
	Hsh Hash
	Lnk Twig7
}

func (s Span7) String() string {
	return fmt.Sprintf("%v",s.Hsh)
}

type Span6 struct{ // 56 bytes 
	Hsh Hash
	Top Kndx
	Lnk Twig6
}

func (s Span6) String() string {
	return fmt.Sprintf("%v",s.Hsh)
}

type Knod struct{ // 256 bytes // knowledge node, Tnod is a text representation
	Op   byte  // op except 0 = hash slot free, 255 = hash slot available due to allocation bounce on content size
	Date Kdate
	Tag  Kndx  // index to a well-known universal label for un-tracked tagging, expanded to hash on transmission or zero, also may be used as salt
	H    Span7
	Idpt byte  // principal signer depth, 1 = self
	Time Ktime
	Iloc Kndx  // principal signer quick link
	P    Span6
	Pchk Cksm  // checksum of Pchk of nodes held... if two lodges differ on this, synchronization is indicated
	Rchk Cksm  // checksum of Rchk of nodes held... if two lodges differ on this, synchronization is indicated
	R    Span6
	S    Sign
}

func op2string(o byte) string {
	kop := "UNKWN"
	switch o {
	case 0:
		kop = "KFREE"
	case 1:
		kop = "LABEL"
	case 2:
		kop = "INTRO"
	case 3:
		kop = "TEMPO"
	case 4:
		kop = "DFOLD"
	case 5:
		kop = "NMACC"
	case 6:
		kop = "ITEXT"
	case 7:
		kop = "REFR"
	}
	return kop
}

func (k Knod) String() string {
	return fmt.Sprintf("%s %v %v %0.2x %0.8x %0.8x %v %v %v %v", op2string(k.Op), k.Date, k.Time, k.Idpt, k.Iloc, k.Tag, k.H, k.P, k.R, k.S) // swap k.Tag for hash
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

	t := Knod { 
	0,
	Kdate {0,0,0},
	0,
	Span7 {Hash {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}, Twig7 {0,0,0,0,0,0,0}},
	0,
	Ktime {0,0,0},
	0,
	Span6 {Hash {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},0, Twig6 {0,0,0,0,0,0}},
	0,
	0,
	Span6 {Hash {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},0, Twig6 {0,0,0,0,0,0}},
	Sign {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}}


	fmt.Println( t )
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
