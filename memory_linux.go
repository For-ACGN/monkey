package monkey

import (
	"fmt"
	"syscall"
)

func writeMemory(address uintptr, data []byte) {
	protect := syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
	err := mProtectCrossPage(address, len(data), protect)
	if err != nil {
		panic(fmt.Sprintf("failed to call Mprotect: %s", err))
	}
	target := readMemory(address, len(data))
	copy(target, data)
	err = mProtectCrossPage(address, len(data), syscall.PROT_READ|syscall.PROT_EXEC)
	if err != nil {
		panic(fmt.Sprintf("failed to call Mprotect: %s", err))
	}
}

func mProtectCrossPage(address uintptr, length int, protect int) error {
	pageSize := syscall.Getpagesize()
	for p := pageStart(address); p < address+uintptr(length); p += uintptr(pageSize) {
		page := readMemory(p, pageSize)
		if err := syscall.Mprotect(page, protect); err != nil {
			return err
		}
	}
	return nil
}
