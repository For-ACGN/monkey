package monkey

import (
	"reflect"
	"syscall"
	"unsafe"
)

func entryAddress(p uintptr, l int) []byte {
	var b []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b)) // #nosec
	sh.Data = p
	sh.Len = l
	sh.Cap = l
	return b
}

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}
