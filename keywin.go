package keymgr

//#cgo CFLAGS: -I./C -funsigned-char
//#cgo LDFLAGS: -lgdi32 -lCredui
//#include <windows.h>
//#include <stdio.h>
//#include <wincred.h>
//#include <errno.h>
//#include <wchar.h>
import "C"
import "unsafe"

// LoadCredentials :
// Load Credentials from Windows Credential Manager
func loadCredentials(Target string) (UserName string, Password string, rc int) {
	var pCred C.PCREDENTIALA
	TargetName := C.CString(Target)
	defer C.free(unsafe.Pointer(TargetName))
	if C.CredReadA(TargetName, C.CRED_TYPE_GENERIC, 0, &pCred) != 0 {
		UserName = C.GoString(pCred.UserName)
		Password = C.GoString((*C.char)(unsafe.Pointer(pCred.CredentialBlob)))
	} else {
		rc = 1
	}
	return UserName, Password, rc
}

// RemoveCredentials :
// Delete Credentials from Windows Credential Manager
func removeCredentials(Target string) int {
	TargetName := C.CString(Target)
	code := C.CredDeleteA(TargetName, C.CRED_TYPE_GENERIC, 0)
	return int(code)
}
// SaveCredentials :
// Save Credentials in Windows Credential Manager
func saveCredentials(Target string, UserName string, Password string) int {
	TargetName := C.CString(Target)
	UserC := C.CString(UserName)
	PassC := C.CString(Password)
	defer C.free(unsafe.Pointer(TargetName))
	defer C.free(unsafe.Pointer(PassC))
	defer C.free(unsafe.Pointer(UserC))

	var cred C.CREDENTIALA
	C.memset(unsafe.Pointer(&cred), 0, C.sizeof_CREDENTIALA)
	cred.Type = C.CRED_TYPE_GENERIC
	cred.TargetName = TargetName
	cred.CredentialBlobSize = C.ulong(C.strlen(PassC) + 1)
	cred.CredentialBlob = (C.LPBYTE)(unsafe.Pointer(PassC))
	cred.Persist = C.CRED_PERSIST_LOCAL_MACHINE
	cred.UserName = UserC
	rc := C.CredWriteA(&cred, 0)
	return int(rc)
}

// GuiCredentials :
// Open User and Password login dialog
func guiCredentials(Title, Target, TargetName string, Attempt int) (UserName string, PassWord string, rc int) {

	TargetNameC := C.CString(TargetName)
	defer C.free(unsafe.Pointer(TargetNameC))

	var gui C.CREDUI_INFOA

	var szPswd [C.CREDUI_MAX_PASSWORD_LENGTH + 1]C.char
	var szName [C.CREDUI_MAX_USERNAME_LENGTH + 1]C.char

	pszName := &szName[0]
	pszPswd := &szPswd[0]

	pszMess := C.CString("Enter Credentials for: " + Target)
	pszCapt := C.CString("TitleC " + Title)
	defer C.free(unsafe.Pointer(pszMess))
	defer C.free(unsafe.Pointer(pszCapt))

	var fSave C.BOOL

	gui.cbSize = C.sizeof_CREDUI_INFOA
	gui.hwndParent = nil
	gui.pszMessageText = pszMess
	gui.pszCaptionText = pszCapt
	gui.hbmBanner = nil
	fSave = C.TRUE

	C.memset(unsafe.Pointer(pszPswd), 0, C.CREDUI_MAX_PASSWORD_LENGTH+1)
	C.memset(unsafe.Pointer(pszName), 0, C.CREDUI_MAX_USERNAME_LENGTH+1)

	var Expression C.ulong

	if Attempt == 0 {
		Expression = 0
	} else {
		Expression = C.CREDUI_FLAGS_INCORRECT_PASSWORD
	}

	dwErr := C.CredUIPromptForCredentialsA(
		&gui,                           // CREDUI_INFOA structure
		TargetNameC,                    // (name of the Credential Stored)
		nil,                            // Reserved
		0,                              // Reason
		pszName,                        // User name
		C.CREDUI_MAX_USERNAME_LENGTH+1, // Max number of char for user name
		pszPswd,                        // Password
		C.CREDUI_MAX_PASSWORD_LENGTH+1, // Max number of char for password
		&fSave,                         // State of save check box
		C.CREDUI_FLAGS_GENERIC_CREDENTIALS| // flags
			C.CREDUI_FLAGS_ALWAYS_SHOW_UI|
			C.CREDUI_FLAGS_DO_NOT_PERSIST|
			Expression)

	rc = int(dwErr)

	if rc == 0 {
		UserName = C.GoStringN(pszName, C.CRED_MAX_CREDENTIAL_BLOB_SIZE)
		PassWord = C.GoStringN(pszPswd, C.CRED_MAX_USERNAME_LENGTH)
	}
	return UserName, PassWord, rc
}