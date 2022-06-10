package win32

import (
	"gitee.com/wengo/go-win32/errco"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	// Library
	libkernel32 *windows.LazyDLL

	// Functions
	readProcessMemory  *windows.LazyProc
	writeProcessMemory *windows.LazyProc
)

func init() {
	// Library
	libkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	readProcessMemory = libkernel32.NewProc("WriteProcessMemory")
	writeProcessMemory = libkernel32.NewProc("WriteProcessMemory")
}

// 📑 https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-readprocessmemory
func ReadProcessMemory(hProcess HPROCESS,
	baseAddress uintptr, buffer []byte) (numBytesRead uint64, e error) {

	ret, _, err := syscall.Syscall6(readProcessMemory.Addr(), 5,
		uintptr(hProcess), baseAddress, uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)), uintptr(unsafe.Pointer(&numBytesRead)), 0)
	if ret == 0 {
		numBytesRead, e = 0, errco.ERROR(err)
	}
	return
}

// 📑 https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-writeprocessmemory
func WriteProcessMemory(hProcess HPROCESS,
	baseAddress uintptr, data []byte) (numBytesWritten uint64, e error) {

	ret, _, err := syscall.Syscall6(writeProcessMemory.Addr(), 5,
		uintptr(hProcess), baseAddress, uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)), uintptr(unsafe.Pointer(&numBytesWritten)), 0)
	if ret == 0 {
		numBytesWritten, e = 0, errco.ERROR(err)
	}
	return
}
