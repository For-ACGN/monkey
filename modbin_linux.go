package monkey

import (
	"fmt"
	"syscall"
)

func modifyBinary(target uintptr, bytes []byte) {
	protect := syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
	err := mProtectCrossPage(target, len(bytes), protect)
	if err != nil {
		panic(fmt.Sprintf("failed to call Mprotect: %s", err))
	}
	function := entryAddress(target, len(bytes))
	copy(function, bytes)
	err = mProtectCrossPage(target, len(bytes), syscall.PROT_READ|syscall.PROT_EXEC)
	if err != nil {
		panic(fmt.Sprintf("failed to call Mprotect: %s", err))
	}
}

func mProtectCrossPage(address uintptr, length int, protect int) error {
	pageSize := syscall.Getpagesize()
	for p := pageStart(address); p < address+uintptr(length); p += uintptr(pageSize) {
		page := entryAddress(p, pageSize)
		if err := syscall.Mprotect(page, protect); err != nil {
			return err
		}
	}
	return nil
}
