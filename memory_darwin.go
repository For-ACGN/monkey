package monkey

import (
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

func writeMemory(address uintptr, data []byte) {
	tp := pageStart(address)
	ret := write(
		address, ptrOf(data), len(data), tp, syscall.Getpagesize(),
		syscall.PROT_READ|syscall.PROT_EXEC,
	)
	if ret != 0 {
		panic(fmt.Sprintf("failed to write memory, code: %v", ret))
	}
}

func ptrOf(val []byte) uintptr {
	return (*reflect.SliceHeader)(unsafe.Pointer(&val)).Data // #nosec
}

//go:cgo_import_dynamic mach_task_self mach_task_self "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic mach_vm_protect mach_vm_protect "/usr/lib/libSystem.B.dylib"
func write(target, data uintptr, len int, page uintptr, pageSize, oriProt int) int
