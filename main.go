package main

import (
	"os"
	"fmt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: chip8 [filename]")
		os.Exit(0)
	}

	chip := newChip(os.Args[1])
	handle := newSdl()
	defer handle.cleanup()

	for {
		instruction := uint16(chip.memory[chip.pc]) << 8 | uint16(chip.memory[chip.pc + 1])
		firstByte := instruction >> 12
		chip.pc += 2

		switch firstByte {
		case 0x0:
			chip.clearScreen()
			handle.drawWindow(chip)
		case 0x1:
			chip.pc = instruction & 0x0FFF
		case 0x6:
			register := (instruction & 0x0F00) >> 8
			chip.registers[register] = uint8(instruction & 0x00FF)
		case 0x7:
			register := (instruction & 0x0F00) >> 8
			chip.registers[register] += uint8(instruction & 0x00FF)
		case 0xA:
			chip.i = instruction & 0x0FFF
		case 0xD:
			chip.registers[0xF] = 0
			rows := instruction & 0x000F
			y := chip.registers[(instruction & 0x00FF) >> 4] % 32
			for r := uint16(0); r < rows && y < 32; r++ {
				x := chip.registers[(instruction >> 8) & 0x0F] % 64
				sprite := chip.memory[chip.i + r]
				for bytePos := uint8(0); bytePos < 8 && x < 64; bytePos++ {
					bitSet := sprite & (1 << (7 - bytePos))
					if bitSet > 0 && chip.screen[y][x] {
						chip.screen[y][x] = false
						chip.registers[0xF] = 1
					} else if bitSet > 0 {
						chip.screen[y][x] = true
					}
					x++
				}
				y++
			}
			handle.drawWindow(chip)
		default:
			panic("unknown instruction")
		}
	}
}
