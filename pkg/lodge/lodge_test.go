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

    base, err := ScratchStore("test.store")
    assert.Nil(t.Test,err)
    Dupstore(base,b)

}

func (t *LodgePkgTests) TestRemoveStore() {
    err := os.Remove("test.store")
    assert.Equal(t.Test, nil, err)
}

func (t *LodgePkgTests) TestWorldExists(b * Base) { //TODO: perform recursion

	w85 := "&URU15#@8/)}XLWy?1hsG0w0v.(O76/e6%P"
	wb, _ := z85.Decode(w85)
	var wbh Hash
	for i:=0;i<28;i++{wbh[i]=wb[i]}
	_,found := b.Get0(&wbh)
	assert.True(t.Test,found)
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
