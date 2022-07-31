package main

import "os"

type chip8 struct {
	memory           [4096]byte
	pc               uint16
	index            uint16
	registers        [16]uint8
	delayTimer       uint16
	instructionCount int
	sp               uint16
	stack            [128]uint16
	screen           []bool
	keys             [16]bool
	fontOffset       int
	xLen             int
	yLen             int
	isHighRes        bool
}

func newChip(file string) chip8 {
	chip := chip8{pc: 0x200}
	chip.setLowRes()
	font := []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	chip.fontOffset = 0x50
	for i, v := range font {
		chip.memory[chip.fontOffset+i] = v
	}

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	_, err = f.Read(chip.memory[512:])
	if err != nil {
		panic(err)
	}

	return chip
}

func (c *chip8) setLowRes() {
	c.screen = make([]bool, 32 * 64)
	c.xLen = 64
	c.yLen = 32
	c.isHighRes = false
}

func (c *chip8) setHighRes() {
	c.screen = make([]bool, 64 * 128)
	c.xLen = 128
	c.yLen = 64
	c.isHighRes = true
}

func (c *chip8) clearScreen() {
	for i := 0; i < len(c.screen); i++ {
		c.screen[i] = false
	}
}
