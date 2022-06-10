package win32

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	// Library
	libshell32 *windows.LazyDLL

	shellExecute *windows.LazyProc
)

func init() {
	// Library
	libshell32 = windows.NewLazySystemDLL("shell32.dll")
	shellExecute = libshell32.NewProc("ShellExecuteW")
}

func ShellExecute(hWnd HWND, verb *uint16, file *uint16, args *uint16, cwd *uint16, showCmd int) bool {
	ret, _, _ := syscall.Syscall6(shellExecute.Addr(), 6,
		uintptr(hWnd),
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		uintptr(unsafe.Pointer(args)),
		uintptr(unsafe.Pointer(cwd)),
		uintptr(showCmd),
	)

	return ret != 0
}
