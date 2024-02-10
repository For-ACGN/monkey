package monkey

func buildJMPDirective(addr uintptr) []byte {
	d0 := byte(addr)
	d1 := byte(addr >> 8)
	d2 := byte(addr >> 16)
	d3 := byte(addr >> 24)
	d4 := byte(addr >> 32)
	d5 := byte(addr >> 40)
	d6 := byte(addr >> 48)
	d7 := byte(addr >> 56)
	jmp := []byte{
		0x48, 0xBA, d0, d1, d2, d3, d4, d5, d6, d7, // mov rdx, addr
		0xFF, 0x22, // jmp [rdx]
	}
	return jmp
}
