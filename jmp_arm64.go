package monkey

import "unsafe"

func buildJMPDirective(double uintptr) []byte {
	d0d1 := double & 0xFFFF
	d2d3 := double >> 16 & 0xFFFF
	d4d5 := double >> 32 & 0xFFFF
	d6d7 := double >> 48 & 0xFFFF
	jmp := make([]byte, 0, 24)
	jmp = append(jmp, movImm(0x02, 0, d0d1)...)          // MOVZ x26, double[16:0]
	jmp = append(jmp, movImm(0x03, 1, d2d3)...)          // MOVK x26, double[32:16]
	jmp = append(jmp, movImm(0x03, 2, d4d5)...)          // MOVK x26, double[48:32]
	jmp = append(jmp, movImm(0x03, 3, d6d7)...)          // MOVK x26, double[64:48]
	jmp = append(jmp, []byte{0x4A, 0x03, 0x40, 0xF9}...) // LDR x10, [x26]
	jmp = append(jmp, []byte{0x40, 0x01, 0x1F, 0xD6}...) // BR x10
	return jmp
}

func movImm(opc, shift int, val uintptr) []byte {
	var m uint32 = 26           // rd
	m |= uint32(val) << 5       // imm16
	m |= uint32(shift&3) << 21  // hw
	m |= 0x25 << 23             // const
	m |= uint32(opc&0x03) << 29 // opc
	m |= 0x01 << 31             // sf
	b := make([]byte, 4)
	*(*uint32)(unsafe.Pointer(&b[0])) = m // #nosec
	return b
}
