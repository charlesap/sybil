package main

import (
	"os"
	"testing"
	"regexp"
//	"fmt"
//	"encoding/binary"

//	"github.com/nofeaturesonlybugs/z85"
	"github.com/stretchr/testify/assert"
        "github.com/charlesap/sybil/pkg/lodge"
)

//
// function tests
//

// TestOp2string0 calls lodge.Op2string with 3, checking
// for a valid return value.
func TestOp2string3(t *testing.T) {
    var op byte = 3
    want := regexp.MustCompile(`PRIVT`)
    msg  := lodge.Op2string(op)
    if !want.MatchString(msg) {
        t.Fatalf(`Op2string(3) = %q, want match for %#q`, msg, want)
    }
}

//
// procedural tests
//

type LodgeCmdTests struct {
	Test *testing.T
}

func TestRunner(t *testing.T) {

    var base lodge.Base

    cmdtest:= LodgeCmdTests{Test: t}

    t.Run("A=init", func(t *testing.T) {
        cmdtest.TestInitializeStore(&base)
//        test.TestCreateMasterUser()
//        test.TestCreateUserTwice()
    })
    t.Run("A=create", func(t *testing.T) {
        cmdtest.TestCreateRegularUser()
//        test.TestCreateConfirmedUser()
//        test.TestCreateMasterUser()
//        test.TestCreateUserTwice()
    })
    t.Run("A=login", func(t *testing.T) {
        cmdtest.TestLoginRegularUser()
//        test.TestLoginConfirmedUser()
//        test.TestLoginMasterUser()
    })
    t.Run("A=cleanup", func(t *testing.T) {
        cmdtest.TestRemoveStore()
    })
}

func (t *LodgeCmdTests) TestInitializeStore(b * lodge.Base) {

	base, err := lodge.ScratchStore("test.store")
	assert.Nil(t.Test,err)
	lodge.Dupstore(base,b)
}

func (t *LodgeCmdTests) TestRemoveStore() {
    err := os.Remove("test.store")
    assert.Equal(t.Test, nil, err)
}


func (t *LodgeCmdTests) TestCreateRegularUser() {
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

func (t *LodgeCmdTests) TestLoginRegularUser() {
}

