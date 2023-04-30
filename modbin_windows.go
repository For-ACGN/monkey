package monkey

import (
	"syscall"
	"unsafe"
)

func modifyBinary(target uintptr, bytes []byte) {
	const pageReadWrite = 0x40
	vp := syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")
	var old uint32
	result, _, err := vp.Call(
		target, uintptr(len(bytes)), pageReadWrite, uintptr(unsafe.Pointer(&old)),
	) // #nosec
	if result == 0 {
		panic(err)
	}
	function := entryAddress(target, len(bytes))
	copy(function, bytes)
	var ignore uint32
	result, _, err = vp.Call(
		target, uintptr(len(bytes)), uintptr(old), uintptr(unsafe.Pointer(&ignore)),
	) // #nosec
	if result == 0 {
		panic(err)
	}
}
