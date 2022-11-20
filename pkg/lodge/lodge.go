package lodge

import (
	"fmt"
	"os"
	"time"
	"crypto/rand"
	"crypto/ed25519"
        "path/filepath"
	"encoding/binary"
//	"reflect"

//	"github.com/ncw/directio"
	"github.com/nofeaturesonlybugs/z85"
)

const(
	wpub85 = "=GghkI+R0ESP@Yp/a-!ug<u6!B6=RBg1n=.anj*("
	wpriv85 = "NaVE{n8nqx=f+GIOL!wWVCa}D)C+wLu#2*6fi]L.=GghkI+R0ESP@Yp/a-!ug<u6!B6=RBg1n=.anj*("

	UNINITIALIZED int = 0 // Lodge status
	UNPREPARED int = 1
	AVAILABLE int = 2
	LOADING int = 3
	STANDBY int = 4

	KFREE byte = 0
	LABEL byte = 1
	INTRO byte = 2
	PRIVT byte = 3
	TEMPO byte = 4
	DFOLD byte = 5
	NMACC byte = 6
	ITEXT byte = 7
	REFER byte = 8
	KROOT byte = 251
	BDURL byte = 252
	BODYN byte = 253
	BODY0 byte = 254
	KBNCE byte = 255



//	knodSize = int(unsafe.Sizeof(Knod{}))
//	bodySize = int(unsafe.Sizeof(Body{}))
)

type Base struct{
	Status int
	Store *os.File
	Fqdn string
	Subn string
	StoreName string
}

type Lodge interface{
	Init(a,b string) error
	Prepare() error
}

func (b Base) Init (fn string) (e error) {
	b.Status = UNINITIALIZED 
	b.StoreName = fn
	fmt.Println("\nInitializing Lodge\n")

	baseName := filepath.Base(os.Args[0])

	b.Store, e = os.Open(fn)
	if e != nil {
		return e
	}
	b.Status = UNPREPARED

	var k Knod
	k,e = b.KnodByIndex(Kndx{0,0,0,0,0,0})
	if e != nil {
		return e
	}

	fmt.Println("\nPreparing Lodge %v\n",k)

	b.Status = AVAILABLE
	fmt.Println("\nLodge available\n")

	Emit(baseName)

	return e
}

func (b Base) KnodByIndex (i Kndx ) (Knod, error) {

//	var k *Knod

//	block := directio.AlignedBlock(directio.BlockSize)
//	_, err := io.ReadFull(b.StoreA, 0)

	k := Knod{}
	e := binary.Read(b.Store, binary.LittleEndian, &k)

	fmt.Println("result...",k.Op,Op2string(k.Op))

	return k, e
}

func (b Base) BodyByIndex (i Kndx ) (*Body, error) {

	return nil, nil
}

type Kdate [6]byte

type Btlen [3]byte

type Btext [252]byte

func (k Kdate) String() string {
	return fmt.Sprintf("YYMMDD")
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

type Slst [48] byte

type Kndx  [6] byte // index of a Knod in block device of Knods

type Cksm  uint32 // checksum of checksums of nodes held by this node or
                  // checksum of checksums of nodes referring to this node


type Knod struct{ // 256 bytes // knowledge node, Tnod is a text representation
	Op   byte  // op except 0 = hash slot free, 255 = hash slot available due to allocation bounce on content size
	Idpt byte  // principal signer depth, 1 = self
	Date Kdate
	Hk   Hash
	Hr   Hash
	Pchk Cksm  // checksum of Rchk of nodes held... if two lodges differ on this, synchronization is indicated
	Itag Kndx  // index to a well-known universal label for un-tracked tagging, expanded to hash on transmission or zero, also may be used as salt
	Ptag Kndx
	Tp   Slst
	Rchk Cksm  // checksum of Pchk of nodes held... if two lodges differ on this, synchronization is indicated
	Ttag Kndx 
	Rtag Kndx
	Tr   Slst
	S    Sign
}

func Op2string(o byte) string {

	switch o {
	case KFREE:
		return "KFREE"
	case LABEL:
		return "LABEL"
	case INTRO:
		return "INTRO"
	case PRIVT:
		return "PRIVT"
	case TEMPO:
		return "TEMPO"
	case DFOLD:
		return "DFOLD"
	case NMACC:
		return "NMACC"
	case ITEXT:
		return "ITEXT"
	case REFER:
		return "REFR"
	}
	if o == KROOT { return "KROOT" }
	if o == BDURL { return "BDURL" }
	if o == BODYN { return "BODYN" }
	if o == BODY0 { return "BODY0" }
	if o == KBNCE { return "KBNCE" }
	return "UNKWN"
}

func (k Knod) UnixTime() time.Time {
	return time.Unix( int64(k.Date[0])+
	                  int64(k.Date[1])*256+
	                  int64(k.Date[2])*65536+
	                  int64(k.Date[3])*16777216+
	                  int64(k.Date[4])*4294967296+
	                  int64(k.Date[5])*1099511627776,0)
}

func (k Knod) Archive() string {
	hs:= Hash {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	hp:= Hash {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	t:=k.UnixTime()

	//255 character string output so the 256th may be a newline 
	return fmt.Sprintf("%v %v%s%s%s%s%s", Op2string(k.Op), time.Unix(t.UnixMilli(),0), k.Hk.Archive() , k.Hr.Archive(), hs.Archive(), hp.Archive(), k.S.Archive()) 
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
	0,
	Kdate {0,0,0,0,0,0},
	Hash {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
	Hash {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
	0,
	Kndx {0,0,0,0,0,0},
	Kndx {0,0,0,0,0,0},
	Slst {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
	0,
	Kndx {0,0,0,0,0,0},
	Kndx {0,0,0,0,0,0},
	Slst {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
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
