package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

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

func runEmulator(chip chip8, handle sdlHandle) {
	for {
		instr := readInstruction(chip.memory[chip.pc : chip.pc+2])
		chip.pc += 2

		for i := 0; i < 4; i++ {
			fmt.Printf("%x", instr.nibbles[i])
		}
		fmt.Println()
		fmt.Println(chip.delayTimer)

		switch {
		case instr.u16 == 0x00E0: // clear
			chip.clearScreen()
			handle.drawWindow(chip)

		case instr.u16 == 0x00EE: // exit routine
			chip.sp--
			chip.pc = chip.stack[chip.sp]

		case instr.nibbles[0] == 0x1: // jump
			chip.pc = instr.jump

		case instr.nibbles[0] == 0x2: // call routine
			chip.stack[chip.sp] = chip.pc
			chip.pc = instr.jump
			chip.sp++

		case instr.nibbles[0] == 0x3: // skip if equal value
			if chip.registers[instr.nibbles[1]] == instr.value {
				chip.pc += 2
			}

		case instr.nibbles[0] == 0x4: // skip if not equal value
			if chip.registers[instr.nibbles[1]] != instr.value {
				chip.pc += 2
			}

		case instr.nibbles[0] == 0x5: // skip if equal registers
			if chip.registers[instr.nibbles[1]] == chip.registers[instr.nibbles[2]] {
				chip.pc += 2
			}

		case instr.nibbles[0] == 0x6: // set register to value
			chip.registers[instr.nibbles[1]] = instr.value

		case instr.nibbles[0] == 0x7: // add value to register
			chip.registers[instr.nibbles[1]] += instr.value

		case instr.nibbles[0] == 0x8: // arithmetic
			regOne := instr.nibbles[1]
			regTwo := instr.nibbles[2]

			switch instr.nibbles[3] {
			case 0x0: // set first register to second
				chip.registers[regOne] = chip.registers[regTwo]

			case 0x1: // or
				chip.registers[regOne] = chip.registers[regOne] | chip.registers[regTwo]

			case 0x2: // and
				chip.registers[regOne] = chip.registers[regOne] & chip.registers[regTwo]

			case 0x3: // xor
				chip.registers[regOne] = chip.registers[regOne] ^ chip.registers[regTwo]

			case 0x4: // add
				res := uint16(chip.registers[regOne]) + uint16(chip.registers[regTwo])
				if res > 255 {
					chip.registers[0xF] = 1
				}
				chip.registers[regOne] = uint8(res)

			case 0x5: // substract first from second
				if chip.registers[regOne] > chip.registers[regTwo] {
					chip.registers[0xF] = 1
				} else {
					chip.registers[0xF] = 0
				}
				chip.registers[regOne] = chip.registers[regOne] - chip.registers[regTwo]

			case 0x6: // shift right
				chip.registers[0xF] = chip.registers[regOne] & 0x1
				chip.registers[regOne] = chip.registers[regOne] >> 1

			case 0x7: // substract second from first
				if chip.registers[regTwo] > chip.registers[regOne] {
					chip.registers[0xF] = 1
				} else {
					chip.registers[0xF] = 0
				}
				chip.registers[regOne] = chip.registers[regTwo] - chip.registers[regOne]

			case 0xE: // shift left
				chip.registers[0xF] = (chip.registers[regOne] & 0x80) >> 7
				chip.registers[regOne] = chip.registers[regOne] << 1
			}

		case instr.nibbles[0] == 0x9: // skip if not equal registers
			if chip.registers[instr.nibbles[1]] != chip.registers[instr.nibbles[2]] {
				chip.pc += 2
			}

		case instr.nibbles[0] == 0xA: // set index register
			chip.index = instr.jump

		case instr.nibbles[0] == 0xB: // jump with offset
			chip.pc = uint16(instr.jump) + uint16(chip.registers[0x0])

		case instr.nibbles[0] == 0xC: // random
			chip.registers[instr.nibbles[1]] = uint8(rand.Uint32()) & instr.value

		// TODO: cleanup mess
		case instr.nibbles[0] == 0xD: // draw
			chip.registers[0xF] = 0
			rows := instr.nibbles[3]
			y := chip.registers[instr.nibbles[2]] % 32
			for r := uint8(0); r < rows && y < 32; r++ {
				x := chip.registers[instr.nibbles[1]] % 64
				sprite := chip.memory[chip.index+uint16(r)]
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

		case instr.nibbles[0] == 0xE && instr.nibbles[3] == 0xE: // skip if pressed
			key, _ := handle.getKeyPressed()
			if instr.nibbles[1] == key {
				chip.pc += 2
			}

		case instr.nibbles[0] == 0xE && instr.nibbles[3] == 0x1: // skip if not pressed
			key, _ := handle.getKeyPressed()
			if instr.nibbles[1] != key {
				chip.pc += 2
			}

		case instr.nibbles[0] == 0xF && instr.value == 0x07: // set register to delay timer
			chip.registers[instr.nibbles[1]] = uint8(chip.delayTimer)

		case instr.nibbles[0] == 0xF && instr.value == 0x15: // set delay timer to register
			chip.delayTimer = uint16(chip.registers[instr.nibbles[1]])

		case instr.nibbles[0] == 0xF && instr.value == 0x18: // set sound timer to register
			chip.soundTimer = chip.registers[instr.nibbles[1]]

		case instr.nibbles[0] == 0xF && instr.value == 0x1E: // add to index register
			newIndex := chip.index + uint16(chip.registers[instr.nibbles[1]])
			chip.index = newIndex % 4096
			if newIndex >= 4096 {
				chip.registers[0xF] = 1
			}

		case instr.nibbles[0] == 0xF && instr.value == 0x0A: // block for key
			key, ok := handle.getKeyPressed()
			if ok {
				chip.registers[instr.nibbles[1]] = key
			} else {
				chip.pc -= 2
			}

		case instr.nibbles[0] == 0xF && instr.value == 0x29: // set index register to font position
			chip.setFontPosition(chip.registers[instr.nibbles[1]])

		case instr.nibbles[0] == 0xF && instr.value == 0x33: // binary coded decimal conversion
			tvx := chip.registers[instr.nibbles[1]]
			for i := 2; i >= 0; i-- {
				remainder := tvx % 10
				chip.memory[chip.index+uint16(i)] = remainder
				tvx /= 10
			}

		case instr.nibbles[0] == 0xF && instr.value == 0x55: // write registers to memory
			for r := uint8(0); r <= instr.nibbles[1]; r++ {
				chip.memory[chip.index+uint16(r)] = chip.registers[r]
			}

		case instr.nibbles[0] == 0xF && instr.value == 0x65: // write memory to registers
			for r := uint8(0); r <= instr.nibbles[1]; r++ {
				chip.registers[r] = chip.memory[chip.index+uint16(r)]
			}

		default:
			log.Fatal("unknown instruction", instr)
		}

		if chip.delayTimer > 0 {
			chip.delayTimer -= 1
		}

		time.Sleep(time.Millisecond * 5)
	}
}
