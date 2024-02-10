package monkey

func buildJMPDirective(addr uintptr) []byte {
	d0 := byte(addr)
	d1 := byte(addr >> 8)
	d2 := byte(addr >> 16)
	d3 := byte(addr >> 24)
	jmp := []byte{
		0xBA, d0, d1, d2, d3, // mov edx, addr
		0xFF, 0x22, // jmp [edx]
	}
	return jmp
}
