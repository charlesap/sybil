package lodge

import (
	"fmt"
//	"bytes"
	"errors"
	"os"
	"time"
	"crypto/rand"
	"crypto/ed25519"
	"crypto/sha256"
        "path/filepath"
//	"encoding/binary"

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

	HASH = 1
	HASHSIGN = 2
	SIGN = 3
	VERIFY = 4

)

var (
	wpub   ed25519.PublicKey    // for content universally signed identically by anyone
	wpriv  ed25519.PrivateKey
	preUL = [][]string {{"en---:World",      "es---:Mundo",            "fr---:Monde",          "cn---:世界",    "jp---:世界",            "de---:Welt" },
			    {"en---:Day",        "es---:Día",              "fr---:Jour",           "cn---:天",      "jp---:日",              "de---:Tag" },
			    {"en---:Lodge",      "es---:Alojarse",         "fr---:Hôtel",          "cn---:小屋",    "jp---:ロッジ",          "de---:Hütte" },
			    {"en---:Keychain",   "es---:Llavero",          "fr---:Porte-clés",     "cn---:钥匙链",  "jp---:キーホルダー",    "de---:Schlüsselbund" },
			    {"en---:Secret",     "es---:Secreto",          "fr---:Secret",         "cn---:秘密",    "jp---:ひみつ",          "de---:Geheimnis" },
			    {"en---:Instance",   "es---:Instancia",        "fr---:Exemple",        "cn---:实例",    "jp---:実例",            "de---:Beispiel" },
			    {"en---:Principal",  "es---:principal",        "fr---:Directeur",      "cn---:主要的",  "jp---:主要",            "de---:Rektor" },
			    {"en---:Session",    "es---:Sesión",           "fr---:Session",        "cn---:会议",    "jp---:セッション",      "de---:Sitzung" },
			    {"en---:Timekeeper", "es---:Cronometrador",    "fr---:Chronométreur",  "cn---:计时员",  "jp---:タイムキーパー",  "de---:Zeitnehmer" },
			    {"en---:Timestamp",  "es---:Marca de Tiempo",  "fr---:Horodatage",     "cn---:时间戳",  "jp---:タイムスタンプ",  "de---:Zeitstempel" }}
)

type Base struct{
	Status int
	Store *os.File
	Limit uint64
	Fqdn string
	Subn string
	StoreName string
}

type Lodge interface{
	Init(a,b string) error
	Prepare() error
}

func (b * Base) Init (fn string, reinit bool) (br * Base, e error) {

	buf, _ := z85.Decode(wpriv85)
	wpriv = ed25519.PrivateKey(buf)
	buf, _ =  z85.Decode(wpub85)
	wpub =  ed25519.PublicKey(buf)

//	buf = []byte("Hello There")
//
//	signature := ed25519.Sign(wpriv, buf)
//	sig85,_ := z85.Encode(signature)
//	fmt.Printf("Signature! %s\n",sig85)
//
//	binsig,_ :=z85.Decode(sig85)
//
//	if (bytes.Compare(signature,binsig)==0) {fmt.Println("enc85/dec85 works...")}
//
//	ok := ed25519.Verify(wpub, buf, binsig)
//
//	if ok {fmt.Println("signature check successful")}

	b.Status = UNINITIALIZED 
	b.StoreName = fn
	fmt.Println("\nPreparing Lodge")

	baseName := filepath.Base(os.Args[0])

	b.Store, e = os.OpenFile(fn, os.O_RDWR, 0777)
	if e != nil {
		return nil, e
	}

	fi, err := b.Store.Stat()
	if err != nil {
		return nil, err
	}

	b.Limit = uint64(fi.Size() / 256)

	b.Status = UNPREPARED

	var k *Knod
	k,e = b.ReadKnodBlock(0)
	if e != nil {
		return nil, e
	}

	if (k.Op != KROOT) || (reinit == true) {
		if (k.Op != KROOT) && (reinit == false) {
			return nil, errors.New("datastore not initialized and reinitialize not requested") 
		}else{
			if k.Op == KROOT {
				fmt.Println("\nRe-initializing Lodge")
			}else{
				fmt.Println("\nInitializing Lodge")
			}
			k.Op=KROOT
			e = b.WriteKnodBlock(k,0)
			if e != nil {return nil, e}

			e = b.mintPreULs()
			if e != nil {return nil, e}

			fmt.Println("\nPrepared Lodge")
		}
	}else{
		fmt.Println("\nServing from Existing Lodge")
	}

	b.Status = AVAILABLE
	fmt.Println("\nLodge available")

	Emit(baseName)


	return b, e
}

func (b * Base) place2(kt *Knod, kb *Body) (e error) {
	e = nil

	tloc,e:=hash2block(&kt.Hk,0,b.Limit)
	if e==nil{
		if b.isfree(tloc,2) {
			e = b.WriteKnodBlock(kt,tloc)
			if e !=nil { return e }
			e = b.WriteBodyBlock(kb,tloc+1)
			if e !=nil { return e }
		}else{
			st, e := b.ReadKnodBlock(tloc)
			if e != nil { return e }
			if st.Hk == kt.Hk {
				fmt.Println("knod already present")
			}else{
				return errors.New("something else here, gotta bounce.")
			}
		}
	}
	return e
}

func (b * Base) mintPreULs() (e error){
	e = nil

	for _,v:=range preUL {
		kt,kb := MintLabel(v[0])
		kt.HashSignVerify ( HASHSIGN, nil,nil,nil,&kb,nil,nil )

		e = b.place2(&kt,&kb)
		if e != nil {return e}

//		fmt.Println( kt.Archive() )
//		fmt.Println( kb.Archive() )
	}
	return e
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

//func (h Hash) String() string {
//	return fmt.Sprintf("%s",h)
//}

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


// hash, sign, or verify the concatenated binary value of the op, date, tag hash, parent hash, ref hash, text content.
func (k *Knod) HashSignVerify (svo int, ks,kp,kt *Knod, b,bs,bv *Body ) (ok bool) {

	// knod signature form without body text in bytes:
	// 1  @ 0 Op, universal LABEL is 1
	// 1  @ 1 Depth, zero if universal label
	// 6  @ 2 Date
	// 28 @ 8 Self hash
	// 28 @ 36 Reference hash, may be zeros if universal label or not-referring, may be dangling (not have a native knod to refer to)
	// 28 @ 64 Hash of signer, may be zeros if universal label, in which case signed by wpriv
	// 28 @ 92 Parent hash, may be zeros if universal label
	// 28 @ 120 Tag hsh, may be zeros if universal label or tagless
	// 148 total bytes + length of body text if any

	blen:=0
	slen:=0
	vlen:=0
	if (b!=nil){
		blen = int(b.Len[0])+(int(b.Len[1])*256)+(int(b.Len[2])*65536)
	}
	if (bs!=nil){
		slen =int(bs.Len[0])+(int(bs.Len[1])*256)+(int(bs.Len[2])*65536)
	}
	if (bv!=nil){
		vlen =int(bv.Len[0])+(int(bv.Len[1])*256)+(int(bv.Len[2])*65536)
	}

	buf := make([]byte,148+blen) // relying on being initialized to zero

	buf[0] = k.Op
	buf[1] = k.Idpt
	for i:=0;i<6;i++ { buf [2+i]= k.Date[i] }
	if k.Op != LABEL { for i:=0;i<28;i++ { buf [36+i]= k.Hr[i] } }
	if k.Op != LABEL { for i:=0;i<28;i++ { buf [64+i]= ks.Hk[i] } }
	if k.Op != LABEL { for i:=0;i<28;i++ { buf [92+i]= kp.Hk[i] } }
	if k.Op != LABEL { for i:=0;i<28;i++ { buf [120+i]= kt.Hk[i] } }

	if blen > 0 { for i:=0;i<blen;i++ { buf [148+i]= b.Text[i] } }

	skey := wpriv
	vkey := wpub
	if (k.Op != LABEL) && (slen > 0) {
		skey, _ = z85.Decode(string(bs.Text[:slen]))
	}
	if (k.Op != LABEL) && (vlen > 0) {
		vkey, _ = z85.Decode(string(bs.Text[:vlen]))
	}

	if (svo == HASH) || (svo == HASHSIGN) {
		h := sha256.New()
		h.Write(buf)
		hash := h.Sum(nil)
		for i:=0;i<28;i++{ k.Hk[i]=hash[i] }
		for i:=0;i<28;i++{ buf [8+i]= k.Hk[i] }
		ok = true
	}else{ //buffer needs self-hash after hashing to be included in signing
		for i:=0;i<28;i++ { buf [8+i]= k.Hk[i] }
	}

	if (svo == SIGN) || (svo == HASHSIGN) {
		signature := ed25519.Sign(skey, buf)
		for i:=0;i<64;i++{ k.S[i]=signature[i] }
		ok = true
	}

	if svo == VERIFY {
		ok = ed25519.Verify(vkey, buf, k.S[:] )

	}


	return ok
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

// temporacles have null Hr, all other principals have an Hr pointing to a timestamp generated by a temporacle

func MintPrincipal(pub, priv []byte) (Knod, Knod, Body, Knod, Body) { // mints a principal entity from an ed25519 private key and an ed25519 public key

	k0:= ZeroKnod()
	k0.Op = 4

	k1:= ZeroKnod()
	b1:= ZeroBody()
	k1.Op = 2
	b1.Mode = 254

//	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
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

	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	_ ,_ ,_ ,_ ,_  = MintPrincipal(pub,priv)
//	t0,t1,b1,t2,b2 := MintPrincipal(pub,priv)

//	fmt.Println( t0.Archive() )
//	fmt.Println( t1.Archive() )
//	fmt.Println( b1.Archive() )
//	fmt.Println( t2.Archive() )
//	fmt.Println( b2.Archive() )
//	fmt.Println()
//	fmt.Println()
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

