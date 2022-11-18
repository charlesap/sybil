package lodge

import (
	"fmt"
	"os"
	"time"
	"crypto/rand"
	"crypto/ed25519"
        "path/filepath"

	"github.com/ncw/directio"
	"github.com/nofeaturesonlybugs/z85"
)

const(
	wpub85 = "=GghkI+R0ESP@Yp/a-!ug<u6!B6=RBg1n=.anj*("
	wpriv85 = "NaVE{n8nqx=f+GIOL!wWVCa}D)C+wLu#2*6fi]L.=GghkI+R0ESP@Yp/a-!ug<u6!B6=RBg1n=.anj*("

	UNINITIALIZED int = iota // Lodge status
	UNPREPARED
	AVAILABLE
	LOADING
	STANDBY
)

type Base struct{
	Status int
	StoreA *os.File
	StoreB *os.File
	Fqdn string
	Subn string
	StoreNameA string
	StoreNameB string
}

type Lodge interface{
	Init(a,b string) error
	Prepare() error
}

func (b Base) Init (fna, fnb string) (e error) {
	b.Status = UNINITIALIZED 
	b.StoreNameA = fna
	b.StoreNameB = fnb
	fmt.Println("\nInitializing Lodge\n")

	baseName := filepath.Base(os.Args[0])

	b.StoreA, e = directio.OpenFile(fna, os.O_RDONLY, 0666)
	if e != nil {
		return e
	}
	b.StoreB, e = directio.OpenFile(fnb, os.O_RDONLY, 0666)
	if e != nil {
		return e
	}
	b.Status = UNPREPARED

	fmt.Println("\nPreparing Lodge\n")

	b.Status = AVAILABLE
	fmt.Println("\nLodge available\n")

	Emit(baseName)

	return e
}

func (b Base) Prepare () error {

	return nil
}

type Kdate [3]byte

type Btlen [3]byte

type Btext [252]byte

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
		kop = "PRIVT"
	case 4:
		kop = "TEMPO"
	case 5:
		kop = "DFOLD"
	case 6:
		kop = "NMACC"
	case 7:
		kop = "ITEXT"
	case 8:
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
	Mode byte  // 0 = lookup empty and terminal, 252 = url to content, 253 = text continued, 254 = text start, 255 = lookup empty but bounce
	Len Btlen  // text length remaining including this
	Text Btext // utf8
}

func (b Body) Archive() string {
	lentag:="---"
	if ((b.Len[2] == 0) && (b.Len[1] == 0) && (b.Len[0] < 253)) {
		lentag = fmt.Sprintf("%03d",b.Len[0])
	}
	return fmt.Sprintf("%s%s",lentag,string(b.Text[:]))
}

func Hash0 (k Knod) Hash {

	h:= Hash {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}

	return h
}

func Hash1 (k Knod, b Body ) Hash {

	h:= Hash {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}

	return h
}

// sign the concatenated binary value of the op, date, tag hash, parent hash, ref hash, text content.

func Sign0 (ks,k Knod) Sign {

	s:= Sign {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	return s
}

func Sign1 (ks,k Knod, b Body ) Sign {

	s:= Sign {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	return s
}

func ZeroBody () Body {
	return Body { 0, Btlen { 0, 0, 0 }, Btext {
		32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,
		32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,
		32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,
		32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,
		32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,
		32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,
		32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,
		32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32}}
}

func ZeroKnod () Knod {

	return Knod { 
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

}

func MintLabel(s string) (Knod, Body) { // mints a universal label from a string up to 252 characters in length

	k:= ZeroKnod()
	b:= ZeroBody()
	k.Op = 1
	b.Mode = 254
	l := len(s)
	b.Len[0]=byte(l%256) // bytes 1 and 2 will be zero
	for i:=0;i<l && i < 252; i++ {
		b.Text[i]=s[i]
	}
	return k,b
}

func MintPrincipal() (Knod, Knod, Body, Knod, Body) { // mints a principal entity from an ed25519 private key and an ed25519 public key

	k0:= ZeroKnod()
	k0.Op = 4

	k1:= ZeroKnod()
	b1:= ZeroBody()
	k1.Op = 2
	b1.Mode = 254

	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	pub85, _ := z85.Encode(pub[:])
	priv85, _ := z85.Encode(priv[:])

	l := len(pub85)
	b1.Len[0]=byte(l%256) // bytes 1 and 2 will be zero
	for i:=0;i<l && i < 252; i++ {
		b1.Text[i]=pub85[i]
	}

	k2:= ZeroKnod()
	b2:= ZeroBody()
	k2.Op = 3
	b2.Mode = 254

	l = len(priv85)
	b2.Len[0]=byte(l%256) // bytes 1 and 2 will be zero
	for i:=0;i<l && i < 252; i++ {
		b2.Text[i]=priv85[i]
	}


	return k0,k1,b1,k2,b2
}

func Emit(name string) {

	t,b := MintLabel("World")
	fmt.Println( t.Archive() )
	fmt.Println( b.Archive() )
	t,b = MintLabel("Day")
	fmt.Println( t.Archive() )
	fmt.Println( b.Archive() )
	t,b = MintLabel("Lodge")
	fmt.Println( t.Archive() )
	fmt.Println( b.Archive() )
	t,b = MintLabel("Keychain")
	fmt.Println( t.Archive() )
	fmt.Println( b.Archive() )
	t,b = MintLabel("Privatekey")
	fmt.Println( t.Archive() )
	fmt.Println( b.Archive() )
	t,b = MintLabel("Instance")  //FQDN/subset
	fmt.Println( t.Archive() )
	fmt.Println( b.Archive() )
	t,b = MintLabel("Principals")
	fmt.Println( t.Archive() )
	fmt.Println( b.Archive() )
	t,b = MintLabel("Sessions")
	fmt.Println( t.Archive() )
	fmt.Println( b.Archive() )
	t,b = MintLabel("Temporacle")
	fmt.Println( t.Archive() )
	fmt.Println( b.Archive() )
	t,b = MintLabel("Timestamp")
	fmt.Println( t.Archive() )
	fmt.Println( b.Archive() )
	t0,t1,b1,t2,b2 := MintPrincipal()
	fmt.Println( t0.Archive() )
	fmt.Println( t1.Archive() )
	fmt.Println( b1.Archive() )
	fmt.Println( t2.Archive() )
	fmt.Println( b2.Archive() )
	fmt.Println()
	fmt.Println()
}

func Format() { // write zeros to empty all blocks

}

func Retrieve() { // hash and bounce (hash again, look again) until match or zeros
	// a (very unlikely) hash resulting in all zeros is simply hashed again.
	// a match with Op being zero or > 252 is not a match

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
