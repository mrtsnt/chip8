package main

import (
	"fmt"
	"log"
	"os"
)

type instruction struct {
	nibbles []uint8
	jumpPosition uint16
	value uint8
}

func readInstruction(bts []byte) instruction {
	u16 := uint16(bts[0]) << 8 | uint16(bts[1])
	return instruction{
		nibbles: []uint8{ bts[0] >> 4, bts[0] & 0x0F, bts[1] >> 4, bts[1] & 0x0F },
		jumpPosition: u16 & 0x0FFF,
		value: uint8(u16 & 0x00FF),
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: chip8 [filename]")
		os.Exit(0)
	}

	chip := newChip(os.Args[1])
	handle := newSdl()
	defer handle.cleanup()

	runEmulator(chip, handle)
}

func runEmulator(chip *chip8, handle *sdlHandle) {
	for {
		instr := readInstruction(chip.memory[chip.pc : chip.pc+2])
		chip.pc += 2

		switch instr.nibbles[0] {
		case 0x0:
			chip.clearScreen()
			handle.drawWindow(chip)

		case 0x1:
			chip.pc = instr.jumpPosition

		case 0x6:
			register := instr.nibbles[1]
			chip.registers[register] = instr.value

		case 0x7:
			register := instr.nibbles[1]
			chip.registers[register] += instr.value

		case 0xA:
			chip.i = instr.jumpPosition

		case 0xD:
			chip.registers[0xF] = 0
			rows := instr.nibbles[3]
			y := chip.registers[instr.nibbles[2]] % 32
			for r := uint8(0); r < rows && y < 32; r++ {
				x := chip.registers[instr.nibbles[1]] % 64
				sprite := chip.memory[chip.i+uint16(r)]
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
			log.Fatal("unknown instruction", instr)
		}
	}
}
