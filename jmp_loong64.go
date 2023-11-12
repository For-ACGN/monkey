package monkey

import "unsafe"

const (
	regR0  uint32 = 0
	regR29        = 29
	regR30        = 30
)

const (
	opORI    uint32 = 0x00E << 22
	opLU12IW        = 0x00A << 25
	opLU32ID        = 0x00B << 25
	opLU52ID        = 0x00C << 22
	opLDD           = 0x0A3 << 22
	opJIRL          = 0x013 << 26
)

func buildJMPDirective(double uintptr) []byte {
	bit110 := (double >> 0) & 0xFFF
	bit3112 := (double >> 12) & 0xFFFFF
	bit5132 := (double >> 32) & 0xFFFFF
	bit6352 := (double >> 52) & 0xFFF
	res := make([]byte, 0, 24)
	res = append(res, wireupOP(opLU12IW, regR29, 0, bit3112)...)      // lu12i.w r29, bit3112
	res = append(res, wireupOP(opORI, regR29, regR29, bit110)...)     // ori     r29, r29, bit110
	res = append(res, wireupOP(opLU32ID, regR29, 0, bit5132)...)      // lu32i.d r29, bit5132
	res = append(res, wireupOP(opLU52ID, regR29, regR29, bit6352)...) // lu52i.d r29, bit6352
	res = append(res, wireupOP(opLDD, regR30, regR29, 0)...)          // ld.d,   r30, r29, 0
	res = append(res, wireupOP(opJIRL, regR0, regR30, 0)...)          // jirl    r0,  r30, 0
	return res
}

func wireupOP(opc uint32, rd, rj uint32, val uintptr) []byte {
	var m uint32 = 0
	switch opc {
	case opORI, opLU52ID, opLDD:
		m |= opc
		m |= (rd & 0x1F) << 0            // rd
		m |= (rj & 0x1F) << 5            // rj
		m |= (uint32(val) & 0xFFF) << 10 // si12
	case opLU12IW, opLU32ID:
		m |= opc
		m |= (rd & 0x1F) << 0             // rd
		m |= (uint32(val) & 0xFFFFF) << 5 // si20
	case opJIRL:
		m |= opc
		m |= (rd & 0x1F) << 0             // rd
		m |= (rj & 0x1F) << 5             // rj
		m |= (uint32(val) & 0xFFFF) << 10 // si16
	}
	res := make([]byte, 4)
	*(*uint32)(unsafe.Pointer(&res[0])) = m
	return res
}
