package w32

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Jabba-Team/jabba/cfg"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	modshell32 = syscall.NewLazyDLL("shell32.dll")
	// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762154(v=vs.85).aspx
	procShellExecuteEx = modshell32.NewProc("ShellExecuteExW")
)

// some of the code below was borrowed from
// https://github.com/AllenDang/w32/blob/65507298e138d537445133ed145a1f2685782b34/shell32.go

func ShellExecuteAndWait(hwnd HWND, lpOperation, lpFile, lpParameters, lpDirectory string, nShowCmd int) error {
	var lpctstrVerb, lpctstrParameters, lpctstrDirectory LPCTSTR
	if len(lpOperation) != 0 {
		lpctstrVerb = LPCTSTR(unsafe.Pointer(syscall.StringToUTF16Ptr(lpOperation)))
	}
	if len(lpParameters) != 0 {
		lpctstrParameters = LPCTSTR(unsafe.Pointer(syscall.StringToUTF16Ptr(lpParameters)))
	}
	if len(lpDirectory) != 0 {
		lpctstrDirectory = LPCTSTR(unsafe.Pointer(syscall.StringToUTF16Ptr(lpDirectory)))
	}
	i := &SHELLEXECUTEINFO{
		fMask:        SEE_MASK_NOCLOSEPROCESS,
		hwnd:         hwnd,
		lpVerb:       lpctstrVerb,
		lpFile:       LPCTSTR(unsafe.Pointer(syscall.StringToUTF16Ptr(lpFile))),
		lpParameters: lpctstrParameters,
		lpDirectory:  lpctstrDirectory,
		nShow:        nShowCmd,
	}
	i.cbSize = DWORD(unsafe.Sizeof(*i))
	return ShellExecuteEx(i)
}

func ShellExecuteEx(pExecInfo *SHELLEXECUTEINFO) error {
	ret, _, _ := procShellExecuteEx.Call(uintptr(unsafe.Pointer(pExecInfo)))
	if ret == 1 && pExecInfo.fMask&SEE_MASK_NOCLOSEPROCESS != 0 {
		s, e := syscall.WaitForSingleObject(syscall.Handle(pExecInfo.hProcess), syscall.INFINITE)
		switch s {
		case syscall.WAIT_OBJECT_0:
			break
		case syscall.WAIT_FAILED:
			return os.NewSyscallError("WaitForSingleObject", e)
		default:
			return errors.New("Unexpected result from WaitForSingleObject")
		}
	}
	errorMsg := ""
	if pExecInfo.hInstApp != 0 && pExecInfo.hInstApp <= 32 {
		switch int(pExecInfo.hInstApp) {
		case SE_ERR_FNF:
			errorMsg = "The specified file was not found"
		case SE_ERR_PNF:
			errorMsg = "The specified path was not found"
		case ERROR_BAD_FORMAT:
			errorMsg = "The .exe file is invalid (non-Win32 .exe or error in .exe image)"
		case SE_ERR_ACCESSDENIED:
			errorMsg = "The operating system denied access to the specified file"
		case SE_ERR_ASSOCINCOMPLETE:
			errorMsg = "The file name association is incomplete or invalid"
		case SE_ERR_DDEBUSY:
			errorMsg = "The DDE transaction could not be completed because other DDE transactions were being processed"
		case SE_ERR_DDEFAIL:
			errorMsg = "The DDE transaction failed"
		case SE_ERR_DDETIMEOUT:
			errorMsg = "The DDE transaction could not be completed because the request timed out"
		case SE_ERR_DLLNOTFOUND:
			errorMsg = "The specified DLL was not found"
		case SE_ERR_NOASSOC:
			errorMsg = "There is no application associated with the given file name extension"
		case SE_ERR_OOM:
			errorMsg = "There was not enough memory to complete the operation"
		case SE_ERR_SHARE:
			errorMsg = "A sharing violation occurred"
		default:
			errorMsg = fmt.Sprintf("Unknown error occurred with error code %v", pExecInfo.hInstApp)
		}
	} else {
		return nil
	}
	return errors.New(errorMsg)
}

func ElevatedRun(name string, arg ...string) (bool, error) {
	ok, err := run("cmd", nil, append([]string{"/C", name}, arg...)...)
	if err != nil {
		rootDir := filepath.Join(cfg.Dir(), "windows")
		ok, err = run("elevate.cmd", &rootDir, append([]string{"cmd", "/C", name}, arg...)...)
	}
	return ok, err
}
func run(name string, dir *string, arg ...string) (bool, error) {
	c := exec.Command(name, arg...)
	if dir != nil {
		c.Dir = *dir
	}
	var stderr bytes.Buffer
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		return false, errors.New(fmt.Sprint(err) + ": " + stderr.String())
	}

	return true, nil
}
func IsAccessDenied(err error) bool {
	fmt.Println(fmt.Sprintf("%v", err))

	if strings.Contains(strings.ToLower(err.Error()), "access is denied") {
		fmt.Println("See https://bit.ly/nvm4w-help")
		return true
	}

	return false
}
