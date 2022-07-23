package main

type instruction struct {
	nibbles []uint8
	jump    uint16
	value   uint8
	u16     uint16
}

func readInstruction(bts []byte) instruction {
	u16 := uint16(bts[0])<<8 | uint16(bts[1])
	return instruction{
		nibbles: []uint8{bts[0] >> 4, bts[0] & 0x0F, bts[1] >> 4, bts[1] & 0x0F},
		jump:    u16 & 0x0FFF,
		value:   uint8(u16 & 0x00FF),
		u16:     u16,
	}
}
