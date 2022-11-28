package lodge

import (
	"os"
	"testing"
	"regexp"
//	"fmt"
//	"encoding/binary"

	"github.com/nofeaturesonlybugs/z85"
	"github.com/stretchr/testify/assert"
)

//
// function tests
//

// TestOp2string0 calls lodge.Op2string with 3, checking
// for a valid return value.
func TestOp2string3(t *testing.T) {
    var op byte = 3
    want := regexp.MustCompile(`PRIVT`)
    msg  := Op2string(op)
    if !want.MatchString(msg) {
        t.Fatalf(`Op2string(3) = %q, want match for %#q`, msg, want)
    }
}


//
// procedural tests
//

type LodgePkgTests struct { 
	Test *testing.T
}

func TestRunner(t *testing.T) {

    var base Base

    test:= LodgePkgTests{Test: t}

    t.Run("A=init", func(t *testing.T) {
        test.TestInitializeStore(&base)
        test.TestWorldExists(&base)
//        test.TestCreateMasterUser()
//        test.TestCreateUserTwice()
    })
    t.Run("A=create", func(t *testing.T) {
        test.TestCreateRegularUser()
//        test.TestCreateConfirmedUser()
//        test.TestCreateMasterUser()
//        test.TestCreateUserTwice()
    })
    t.Run("A=login", func(t *testing.T) {
        test.TestLoginRegularUser()
//        test.TestLoginConfirmedUser()
//        test.TestLoginMasterUser()
    })
    t.Run("A=cleanup", func(t *testing.T) {
        test.TestRemoveStore()
    })
}

func (t *LodgePkgTests) TestInitializeStore(b * Base) {

    var base * Base

    size := int64(1<<30)
    fd, err := os.Create("test.store")
    assert.Nil(t.Test,err)
    _, err = fd.Seek(size-1, 0)
    assert.Nil(t.Test,err)
    _, err = fd.Write([]byte{0})
    assert.Nil(t.Test,err)
    err = fd.Close()
    assert.Nil(t.Test,err)

    base, err = b.Init("test.store",true)
    assert.Nil(t.Test,err)

    b.Status = base.Status
    b.Store = base.Store
    b.Limit = base.Limit
    b.Fqdn = base.Fqdn
    b.Subn = base.Subn
    b.StoreName = base.StoreName
}

func (t *LodgePkgTests) TestRemoveStore() {
    err := os.Remove("test.store")
    assert.Equal(t.Test, nil, err)
}

func (t *LodgePkgTests) TestWorldExists(b * Base) { //TODO: perform recursion

	var wbinhash []byte
	var ok error
	var wbh Hash
	var tloc uint64
	var st *Knod

//	fmt.Print(b)
	e:= assert.NotEqual(t.Test,b.Limit,0) 
	if e {
		w85hash := "&URU15#@8/)}XLWy?1hsG0w0v.(O76/e6%P"
		wbinhash, ok = z85.Decode(w85hash)
		e = assert.Nil(t.Test,ok)
	}
	if e {
		for i:=0;i<28;i++{wbh[i]=wbinhash[i]}
//		fmt.Println(" : ",b.Limit)
		tloc, ok = Hash2block(&wbh,0,b.Limit)
		e = assert.Nil(t.Test,ok)
	}
//	fmt.Println("tloc: ",tloc)
	if e {
		st, ok = b.ReadKnodBlock(tloc)
		e = assert.Nil(t.Test,ok)
	}
	if e {
		m := true
		for i:=0;i<28;i++{if st.Hk[i]!=wbinhash[i] {m=false;}}
		e = assert.Equal(t.Test,m,true)
	}

}

func (t *LodgePkgTests) TestCreateRegularUser() {
//    registerRegularUser := util.TableTest{
//        Method:      "POST",
//        Path:        "/iot/users",
//        Status:      http.StatusOK,
//        Name:        "registerRegularUser",
//        Description: "register Regular User has to return 200",
//        Body: SerializeUser(RegularUser),
//    }
//    response := util.SpinSingleTableTests(t.Test, registerRegularUser)
//    util.LogIfVerbose(color.BgCyan, "IOT/USERS/TEST", response)
}

func (t *LodgePkgTests) TestLoginRegularUser() {
}
