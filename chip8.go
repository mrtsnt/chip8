package main

import (
	"os"
	"time"
)

type chip8 struct {
	memory     [4096]byte
	pc         uint16
	index      uint16
	registers  [16]uint8
	delayTimer uint16
	soundTimer uint16
	sp         uint16
	stack      [128]uint16
	screen     [32][64]bool
}

func newChip(file string) chip8 {
	chip := chip8{pc: 0x200}
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

	for i, v := range font {
		chip.memory[0x50+i] = v
	}

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	_, err = f.Read(chip.memory[512:])
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Millisecond * 17) // 60hz
			if chip.delayTimer > 0 {
				chip.delayTimer -= 1
			}

			if chip.soundTimer > 0 {
				chip.soundTimer -= 1
			}
		}
	}()

	return chip
}

func (c *chip8) clearScreen() {
	for row := 0; row < 32; row++ {
		for col := 0; col < 64; col++ {
			c.screen[row][col] = false
		}
	}
}
