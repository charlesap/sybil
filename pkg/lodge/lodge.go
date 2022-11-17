package lodge

import (
//    a85 "encoding/ascii85"
    "fmt"
    "time"
    "github.com/nofeaturesonlybugs/z85"
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

func (h Hash) Archive() string {
	t, _ := z85.Encode(h[:])
	return t
}

type Sign [64] byte

func (s Sign) Archive() string {
	t, _ := z85.Encode(s[:])
return t
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

func (s Span7) Archive() string {
	return s.Hsh.Archive()
}

type Span6 struct{ // 56 bytes 
	Hsh Hash
	Top Kndx
	Lnk Twig6
}

func (s Span6) Archive() string {
	return s.Hsh.Archive()
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

func Op2string(o byte) string {
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

func (k Knod) UnixTime() time.Time {
	return time.Unix( int64(k.Time[0])+
	                  int64(k.Time[1])*256+
	                  int64(k.Time[2])*65536+
	                  int64(k.Date[0])*16777216+
	                  int64(k.Date[1])*4294967296+
	                  int64(k.Date[2])*1099511627776,0)
}

func (k Knod) Archive() string {
	h:= Hash {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	t:=k.UnixTime()

	//255 character string output so the 256th may be a newline 
	return fmt.Sprintf("%v %v%s%s%s%s%s", Op2string(k.Op), time.Unix(t.UnixMilli(),0), h.Archive() , k.H.Archive(), k.P.Archive(), k.R.Archive(), k.S.Archive()) 
}



type Body struct{ // 252 possble bytes of text
	Pad1 byte  // 0 = lookup empty and terminal, 252 = url to content, 253 = text continued, 254 = text start, 255 = lookup empty but bounce
	Pad2 [3] byte  // text length remaining including this
	Text [252] byte
	                // all utf8
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


	fmt.Println( t.Archive() )
	fmt.Println()
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
