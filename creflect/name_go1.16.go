//go:build go1.11 && !go1.17
// +build go1.11,!go1.17

package creflect

import (
	"unsafe"
)

// name is an encoded type name with optional extra data.
type name struct {
	bytes *byte
}

// #nosec
func (n name) name() (s string) {
	if n.bytes == nil {
		return
	}
	b := (*[4]byte)(unsafe.Pointer(n.bytes))
	hdr := (*stringHeader)(unsafe.Pointer(&s))
	hdr.data = unsafe.Pointer(&b[3])
	hdr.len = int(b[1])<<8 | int(b[2])
	return s
}
