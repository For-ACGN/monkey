//go:build !windows
// +build !windows

package monkey

import (
	"syscall"
)

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}
