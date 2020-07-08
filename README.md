# keymgr (Credentials Manager)
## Golang / Windows Credential Manager Wrapper


A small library that interacts with Go and Windows Credential Manager.

It incorporates Windows Credential GUI. It also uses and is integrated with the status 
The Credentials will be stored with the Following format.

    - Windows Credential Manager - Generic Credentials
    - format:
    - *['program name']~['Server name or Address']*
    
### Example
*This is a dummy connect and returns an int value*.

```go
package main
import (
	"fmt"
	"gopkg.in/ro-ag/keymgr.v1"
)

// This example prepares the keymgr.Cred type
// Login function needs a CallBack type function with interface return
// The Flow Actions in Login function ( this example ):
//  1.- Check if Credentials with name *[keymgr]~[test]* exits
//      1.1.- If exists go to step 5
//  2.- Run a GUI to prompt user and password
//  3.- Store Credentials in Windows Credential Manager
//  4.- Load Credentials ( check if works )
//  5.- Run connect function
//      5.1 - If status is not 0 remove credentials
//      5.2 - Go to Step 2
//      5.3 - if Attempts is >= Limit return nil
//  6.- Return connect interface

func main() {
	Parameters := keymgr.Cred{
		Program: "keymgr",   // Program Name Identifier
		Target:  "Test",       // Target, usually a site or server
		Limit:   3,            // Number of Attempts before return nil
		Debug:   true,         // Internal Log messages
	}
	out := Parameters.Login(connect).(int) // CallBack function
	fmt.Println(out)
}


// This function Simulate any type of connection passing keymgr.CallArgs structure
//  .User : User From Windows Credential Manager
//  .Pass : Password
//  .Target  : Normally Site or Server ( this comes from keymgr.Cred { .Target } type )
//  .Attempt : Current number of attempts to run connect function
// The function will pass when status is set to 0

func connect( Credentials keymgr.CallArgs, status *int )interface{}{
	if Credentials.User == "pepe" && Credentials.Pass == "pepepass" {
		fmt.Println("User  : ", Credentials.User)
		fmt.Println("Pass  : ", Credentials.Pass)
		*status = 0
	} else {
        fmt.Println("  Target: ", Credentials.Target)
		fmt.Println("  User  : ", Credentials.User)
		fmt.Println("  Pass  : ", Credentials.Pass)
		fmt.Println("  Attempt : ", Credentials.Attempts)
		*status = 1
	}
	// Return int
	return *status
}
```
