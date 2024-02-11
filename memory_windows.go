package monkey

import (
	"fmt"
	"syscall"
	"unsafe"
)

const pageReadWrite = 0x40

var virtualProtect = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")

func writeMemory(address uintptr, data []byte) {
	var old uint32
	ret, _, err := virtualProtect.Call(
		address, uintptr(len(data)), pageReadWrite, uintptr(unsafe.Pointer(&old)),
	) // #nosec
	if ret == 0 {
		panic(fmt.Sprintf("failed to call VirtialProtect: %s", err))
	}
	target := readMemory(address, len(data))
	copy(target, data)
	var ignore uint32
	ret, _, err = virtualProtect.Call(
		address, uintptr(len(data)), uintptr(old), uintptr(unsafe.Pointer(&ignore)),
	) // #nosec
	if ret == 0 {
		panic(fmt.Sprintf("failed to call VirtialProtect: %s", err))
	}
}
