//go:build !go1.17
// +build !go1.17

package creflect

import (
	"unsafe"
)

// name is an encoded type name with optional extra data.
type name struct {
	bytes *byte
}

func (n name) name() (s string) {
	if n.bytes == nil {
		return
	}
	b := (*[4]byte)(unsafe.Pointer(n.bytes)) // #nosec

	hdr := (*stringHeader)(unsafe.Pointer(&s)) // #nosec
	hdr.data = unsafe.Pointer(&b[3])           // #nosec
	hdr.len = int(b[1])<<8 | int(b[2])
	return s
}
