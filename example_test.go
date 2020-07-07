package keymgr_test

import (
	"fmt"
	"gopkg.in/ro-ag/keymgr.v1"
	"testing"
)

// ConnectExample:
// Simulate any type of connection it will give user and password
// to pass the return code needs to be 0
// This will be called the X number of times set in limit
func ConnectExample( c keymgr.CallArgs, status *int )interface{}{
	if c.User == "pepe" && c.Pass == "pepepass" {
		fmt.Println("User  : ", c.User)
		fmt.Println("Pass  : ", c.Pass)
		*status = 0
	} else {
		fmt.Println("  User  : ", c.User)
		fmt.Println("  Pass  : ", c.Pass)
		fmt.Println("  Attempt : ", c.Attempts)
		*status = 1
	}
	// Return int
	return *status
}
// Examplekeymgr:
// This example prepares the keymgr.Cred type
// Login function needs a CallBack type function
// Returns empty interface
func ExampleKeymgr() {
	Parameters := keymgr.Cred{
		Program: "keymgr",
		Target:  "Test",
		Limit:   3,
		Debug:   true,
	}
	out := Parameters.Login(ConnectExample).(int)
	fmt.Println(out)
	//Output:
	//User  :  pepe
	//Pass  :  pepepass
	//0
}

func BenchmarkKeymgr(b *testing.B) {
	Parameters := keymgr.Cred{
		Program: "keymgr:",
		Target:  "Test",
		Limit:   3,
	}
	for i := 0; i < b.N; i++ {
		Parameters.Login(ConnectExample)
	}
}
