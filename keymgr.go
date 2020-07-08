// Package keymgr adds a simple wrapper to Windows Credential Manager for the Go language.
// This package requires MinGW-64
//
// Source code and other details for the project are available at GitHub:
//
//   https://github.com/ro-ag/keymgr
//
package keymgr

import (
	"fmt"
	"github.com/kr/pretty"
	"gopkg.in/go-validator/validator.v2"
	"io/ioutil"
	"log"
	"os"
)

// CallArgs :
// Needed for CallBack function, it contains the User, Password, Target ( usually a server or site )
// Attempts is the current time the credentials are prompted
type CallArgs struct {
	Target   string // Target normally is a Site or Server ( already passed in Cred Structure )
	User     string // User from Credential Manager
	Pass     string // PassWord from Credential Manager
	Attempts int    // Current number of attempts
}

// Cred :
// Initial Structure
type Cred struct {
	Program string `validate:"nonzero"` // Program identifier
	Target  string `validate:"nonzero"` // Target normally is a Site or Server
	Limit   int    `validate:"min=1"`   // Limit of login intents before fail
	Debug   bool   // Print keymgr log information
}

func startLog(f bool) {
	if f == false {
		log.SetOutput(ioutil.Discard)
	}
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Lshortfile)
	log.Println("Log started")
}

func (P *Cred) validation() {
	if errs := validator.Validate(*P); errs != nil {
		ErrorMap := errs.(validator.ErrorMap)
		for key, value := range ErrorMap {
			os.Stderr.WriteString(fmt.Sprintf("-> Key(%s) is %s\n", key, value))
		}
		log.Fatal("Error in Cred types")
	}
}

// Login :
// This function needs a CallBack function with interface return
// The Flow Actions in Login function:
//  1.- Check if Credentials with name "*[Program]~[Target]*" exits
//      1.1.- If exists go to step 5
//  2.- Run a GUI to prompt user and password
//  3.- Store Credentials in Windows Credential Manager
//  4.- Load Credentials ( check if works )
//  5.- Run "CallBack" function
//      5.1 - If status is not 0 remove credentials
//      5.2 - Go to Step 2
//      5.3 - if Attempts is >= Limit return nil
//  6.- Return "CallBack" interface
func (P *Cred) Login(CallBack func(CallArgs, *int) interface{}) (OutLogin interface{}) {
	P.validation()

	startLog(P.Debug)

	log.Printf("%# v\n", pretty.Formatter(P))

	targetName := fmt.Sprintf("*[%s]~[%s]*", P.Program, P.Target)

	log.Printf("%# v\n", pretty.Formatter(targetName))

	argsBack := CallArgs{
		Target:   P.Target,
		User:     "",
		Pass:     "",
		Attempts: 0,
	}
	status := -1

	log.Printf("%# v\n", pretty.Formatter(argsBack))

	var rc int
	for status != 0 && argsBack.Attempts < P.Limit {
		log.Println("Start Loop: Attempt ", argsBack.Attempts, "Status ", status)
		argsBack.User, argsBack.Pass, rc = loadCredentials(targetName)
		log.Printf("Load Credentials : %# v\n", pretty.Formatter(argsBack))
		if rc != 0 {
			argsBack.User, argsBack.Pass, rc = guiCredentials(P.Program, P.Target, targetName, argsBack.Attempts)
			if rc != 0 {
				fmt.Println("Cancel Operation")
				os.Exit(1)
			}
			log.Printf("After Gui : %# v\n", pretty.Formatter(argsBack))
			if saveCredentials(targetName, argsBack.User, argsBack.Pass) != 1 {
				panic("Error to Save Credentials")
			}
			log.Println("Save Credentials ", targetName, argsBack.User, argsBack.Pass)
		} else {
			// CallBack Back
			OutLogin = CallBack(argsBack, &status)
			if status != 0 {
				argsBack.Attempts++
				removeCredentials(targetName)
				log.Println("Remove Credentials ", targetName)
			}
		}
	}
	// Clear Credentials for Security
	argsBack.User = ""
	argsBack.Pass = ""
	return OutLogin
}
