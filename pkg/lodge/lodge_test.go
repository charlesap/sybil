package lodge

import (
	"os"
	"testing"
	"regexp"

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

type LodgeTests struct { Test *testing.T}

func TestRunner(t *testing.T) {

    t.Run("A=init", func(t *testing.T) {
        test:= LodgeTests{Test: t}
        test.TestInitializeStore()
//        test.TestCreateConfirmedUser()
//        test.TestCreateMasterUser()
//        test.TestCreateUserTwice()
    })
    t.Run("A=create", func(t *testing.T) {
        test:= LodgeTests{Test: t}
        test.TestCreateRegularUser()
//        test.TestCreateConfirmedUser()
//        test.TestCreateMasterUser()
//        test.TestCreateUserTwice()
    })
    t.Run("A=login", func(t *testing.T) {
        test:= LodgeTests{Test: t}
        test.TestLoginRegularUser()
//        test.TestLoginConfirmedUser()
//        test.TestLoginMasterUser()
    })
    t.Run("A=cleanup", func(t *testing.T) {
        test:= LodgeTests{Test: t}
        test.TestRemoveStore()
    })
}

func (t *LodgeTests) TestInitializeStore() {
    size := int64(1<<30)
    fd, err := os.Create("test.store")
    assert.Equal(t.Test, nil, err)
    _, err = fd.Seek(size-1, 0)
    assert.Equal(t.Test, nil, err)
    _, err = fd.Write([]byte{0})
    assert.Equal(t.Test, nil, err)
    err = fd.Close()
    assert.Equal(t.Test, nil, err)
}

func (t *LodgeTests) TestRemoveStore() {
    err := os.Remove("test.store")
    assert.Equal(t.Test, nil, err)
}

func (t *LodgeTests) TestCreateRegularUser() {
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

func (t *LodgeTests) TestLoginRegularUser() {
}
