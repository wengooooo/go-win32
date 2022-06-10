package win32

import "golang.org/x/sys/windows"

var (
	// Library
	librasapi32 *windows.LazyDLL
)

func init() {
	//is64bit := unsafe.Sizeof(uintptr(0)) == 8

	// Library
	librasapi32 = windows.NewLazySystemDLL("rasapi32.dll")

	// Functions

}
