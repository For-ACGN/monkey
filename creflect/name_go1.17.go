//go:build go1.17
// +build go1.17

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
	i, l := n.readVarInt(1)
	hdr := (*stringHeader)(unsafe.Pointer(&s))
	hdr.data = unsafe.Pointer(n.data(1 + i)) // "non-empty string"
	hdr.len = l
	return
}

// #nosec
func (n name) readVarInt(off int) (int, int) {
	v := 0
	for i := 0; ; i++ {
		x := *n.data(off + i) // "read var int"
		v += int(x&0x7F) << uint(7*i)
		if x&0x80 == 0 {
			return i + 1, v
		}
	}
}

func (n name) data(off int) *byte {
	return (*byte)(add(unsafe.Pointer(n.bytes), uintptr(off))) // #nosec
}
