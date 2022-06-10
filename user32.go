package win32

import (
	"github.com/wengooooo/go-win32/errco"
	"golang.org/x/sys/windows"
	"sync"
	"syscall"
	"unsafe"
)

var (
	// Library
	libuser32       *windows.LazyDLL
	postQuitMessage *windows.LazyProc
	postMessage     *windows.LazyProc
	messageBox      *windows.LazyProc
	sendMessage     *windows.LazyProc
	enumWindows     *windows.LazyProc
)

func init() {
	//is64bit := unsafe.Sizeof(uintptr(0)) == 8

	// Library
	libuser32 = windows.NewLazySystemDLL("user32.dll")

	// Functions
	postQuitMessage = libuser32.NewProc("PostQuitMessage")
	postMessage = libuser32.NewProc("PostMessage")
	sendMessage = libuser32.NewProc("SendMessage")
	messageBox = libuser32.NewProc("MessageBox")
	enumWindows = libuser32.NewProc("EnumWindows")
}

// PostQuitMessage
//
// 📑 https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-postquitmessage
func PostQuitMessage(exitCode int32) {
	syscall.Syscall(postQuitMessage.Addr(), 1,
		uintptr(exitCode),
		0,
		0)
}

// PostMessage
//
// 📑 https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-postmessage
func PostMessage(hWnd HWND, msg Msg, wParam WPARAM, lParam LPARAM) uintptr {
	ret, _, _ := syscall.Syscall6(postMessage.Addr(), 4,
		uintptr(hWnd),
		uintptr(msg),
		uintptr(wParam),
		uintptr(lParam),
		0,
		0)

	return ret
}

// SendMessage
//
// 📑 https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage
func SendMessage(hWnd HWND, msg Msg, wParam WPARAM, lParam LPARAM) uintptr {
	ret, _, _ := syscall.Syscall6(postMessage.Addr(), 4,
		uintptr(hWnd),
		uintptr(msg),
		uintptr(wParam),
		uintptr(lParam),
		0,
		0)

	return ret
}

func MessageBox(hWnd HWND, text, caption string, uType MB) ID {
	ret, _, err := syscall.Syscall6(messageBox.Addr(), 4,
		uintptr(hWnd), uintptr(unsafe.Pointer(Str.ToNativePtr(text))),
		uintptr(unsafe.Pointer(Str.ToNativePtr(caption))), uintptr(uType),
		0, 0)
	if ret == 0 {
		panic(errco.ERROR(err))
	}

	return ID(ret)
}

var (
	_globalEnumWindowsCallback uintptr = syscall.NewCallback(_EnumWindowsProc)
	_globalEnumWindowsFuncs    map[*_EnumWindowsPack]struct{}
	_globalEnumWindowsMutex    = sync.Mutex{}
)

// 📑 https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-enumwindows
func EnumWindows(callback func(hWnd HWND) bool) {
	pPack := &_EnumWindowsPack{f: callback}
	_globalEnumWindowsMutex.Lock()
	if _globalEnumWindowsFuncs == nil { // the set was not initialized yet?
		_globalEnumWindowsFuncs = make(map[*_EnumWindowsPack]struct{}, 1)
	}
	_globalEnumWindowsFuncs[pPack] = struct{}{} // store pointer in the set
	_globalEnumWindowsMutex.Unlock()

	ret, _, err := syscall.Syscall(enumWindows.Addr(), 2,
		_globalEnumWindowsCallback, uintptr(unsafe.Pointer(pPack)), 0)
	if ret == 0 {
		panic(errco.ERROR(err))
	}
}

type _EnumWindowsPack struct{ f func(hWnd HWND) bool }

func _EnumWindowsProc(hWnd HWND, lParam LPARAM) uintptr {
	pPack := (*_EnumWindowsPack)(unsafe.Pointer(lParam))
	retVal := uintptr(0)

	_globalEnumWindowsMutex.Lock()
	_, isStored := _globalEnumWindowsFuncs[pPack]
	_globalEnumWindowsMutex.Unlock()

	if isStored {
		retVal = BoolToUintptr(pPack.f(hWnd))
		if retVal == 0 {
			_globalEnumWindowsMutex.Lock()
			delete(_globalEnumWindowsFuncs, pPack) // remove from the set
			_globalEnumWindowsMutex.Unlock()
		}
	}
	return retVal
}
