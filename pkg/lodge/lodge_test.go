package lodge

import (
	"testing"
	"regexp"
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
