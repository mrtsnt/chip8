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
		for i := 0; i < 10; i++ {
			instr := readInstruction(chip.memory[chip.pc : chip.pc+2])
			execute(handle, &chip, instr)
		}

		if chip.delayTimer > 0 {
			chip.delayTimer -= 1
		}

		time.Sleep(time.Millisecond * 5)
	}
}

func execute(handle sdlHandle, chip *chip8, instr instruction) {
	handle.updateKeyboard(chip)
	chip.pc += 2

	switch {
	case instr.u16 == 0x00E0: // clear
		chip.clearScreen()
		handle.drawWindow(chip)

	case instr.u16 == 0x00EE: // exit routine
		chip.sp--
		chip.pc = chip.stack[chip.sp]

	case instr.u16 == 0x00FE: // switch to low res
		chip.setLowRes()

	case instr.u16 == 0x00FF: // switch to high res
		chip.setHighRes()

	case instr.u16 & 0xFFF0 == 0x00C0: // shift screen down
		offset := int(instr.nibbles[3])
		copy(chip.screen[offset * chip.xLen:], chip.screen)
		for i := 0; i < offset * chip.xLen; i++ {
			chip.screen[i] = false
		}

	case instr.u16 == 0x00FB: // shift screen right
	  for r := 0; r < chip.yLen; r++ {
			offset := r * chip.xLen
			copy(chip.screen[offset+4:offset+chip.xLen], chip.screen[offset:offset+chip.xLen])
			for i := 0; i < 4; i++ {
				chip.screen[offset + i] = false
			}
		}

	case instr.u16 == 0x00FC: // shift screen left
	  for r := 0; r < chip.yLen; r++ {
			offset := r * chip.xLen
			copy(chip.screen[offset:offset+chip.xLen], chip.screen[offset+4:offset+chip.xLen])
			for i := 0; i < 4; i++ {
				chip.screen[offset + chip.xLen - i - 1] = false
			}
		}

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
		regX := instr.nibbles[1]
		regY := instr.nibbles[2]

		switch instr.nibbles[3] {
		case 0x0: // set first register to second
			chip.registers[regX] = chip.registers[regY]

		case 0x1: // or
			chip.registers[regX] = chip.registers[regX] | chip.registers[regY]

		case 0x2: // and
			chip.registers[regX] = chip.registers[regX] & chip.registers[regY]

		case 0x3: // xor
			chip.registers[regX] = chip.registers[regX] ^ chip.registers[regY]

		case 0x4: // add
			res := uint16(chip.registers[regX]) + uint16(chip.registers[regY])
			if res > 255 {
				chip.registers[0xF] = 1
			}
			chip.registers[regX] = uint8(res)

		case 0x5: // substract first from second
			if chip.registers[regX] > chip.registers[regY] {
				chip.registers[0xF] = 1
			} else {
				chip.registers[0xF] = 0
			}
			chip.registers[regX] = chip.registers[regX] - chip.registers[regY]

		case 0x6: // shift right
			chip.registers[0xF] = chip.registers[regX] & 0x1
			chip.registers[regX] = chip.registers[regX] >> 1

		case 0x7: // substract second from first
			if chip.registers[regY] > chip.registers[regX] {
				chip.registers[0xF] = 1
			} else {
				chip.registers[0xF] = 0
			}
			chip.registers[regX] = chip.registers[regY] - chip.registers[regX]

		case 0xE: // shift left
			chip.registers[0xF] = (chip.registers[regX] & 0x80) >> 7
			chip.registers[regX] = chip.registers[regX] << 1
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

	case instr.nibbles[0] == 0xD: // draw
	  if instr.nibbles[3] == 0x0 {
			chip.registers[0xF] = 0
			x := chip.registers[instr.nibbles[1]] % uint8(chip.xLen)
			y := chip.registers[instr.nibbles[2]] % uint8(chip.yLen)

			for r := uint8(0); r < 16 && y < uint8(chip.yLen); r++ {
				for bp := 0; bp < 2; bp++ {
					sprite := chip.memory[chip.index+uint16(r)+uint16(bp)]
					for bx := 0; bx+int(x)+bp*8 < chip.xLen && bx < 8; bx++ {
						loc := chip.xLen*int(y) + int(x) + bx
						if chip.screen[loc] && sprite&(1<<(7-bx)) > 0 {
							chip.screen[loc] = false
							chip.registers[0xF] = 1
						} else if sprite&(1<<(7-bx)) > 0 {
							chip.screen[loc] = true
						}
					}
				}
				y++
			}
		} else {
			chip.registers[0xF] = 0
			spriteRows := instr.nibbles[3]
			x := chip.registers[instr.nibbles[1]] % uint8(chip.xLen)
			y := chip.registers[instr.nibbles[2]] % uint8(chip.yLen)

			for r := uint8(0); r < spriteRows && y < uint8(chip.yLen); r++ {
				sprite := chip.memory[chip.index+uint16(r)]
				for bx := 0; bx+int(x) < chip.xLen && bx < 8; bx++ {
					loc := chip.xLen*int(y) + int(x) + bx
					if chip.screen[loc] && sprite&(1<<(7-bx)) > 0 {
						chip.screen[loc] = false
						chip.registers[0xF] = 1
					} else if sprite&(1<<(7-bx)) > 0 {
						chip.screen[loc] = true
					}
				}
				y++
			}
		}
		handle.drawWindow(chip)

	case instr.nibbles[0] == 0xE && instr.nibbles[3] == 0xE: // skip if pressed
		if chip.keys[chip.registers[instr.nibbles[1]]] {
			chip.pc += 2
		}

	case instr.nibbles[0] == 0xE && instr.nibbles[3] == 0x1: // skip if not pressed
		if !chip.keys[chip.registers[instr.nibbles[1]]] {
			chip.pc += 2
		}

	case instr.nibbles[0] == 0xF && instr.value == 0x07: // set register to delay timer
		chip.registers[instr.nibbles[1]] = uint8(chip.delayTimer)

	case instr.nibbles[0] == 0xF && instr.value == 0x15: // set delay timer to register
		chip.delayTimer = uint16(chip.registers[instr.nibbles[1]])

	case instr.nibbles[0] == 0xF && instr.value == 0x18: // set sound timer to register, not implemented

	case instr.nibbles[0] == 0xF && instr.value == 0x1E: // add to index register
		newIndex := chip.index + uint16(chip.registers[instr.nibbles[1]])
		chip.index = newIndex % 4096
		if newIndex >= 4096 {
			chip.registers[0xF] = 1
		}

	case instr.nibbles[0] == 0xF && instr.value == 0x0A: // block for key
		pressed := false
		for k := 0; k < 16; k++ {
			if chip.keys[k] {
				pressed = true
				chip.registers[instr.nibbles[1]] = uint8(k)
			}
		}

		if !pressed {
			chip.pc -= 2
		}

	case instr.nibbles[0] == 0xF && instr.value == 0x29: // set index register to font position
		chip.index = uint16(chip.fontOffset) + uint16(5*chip.registers[instr.nibbles[1]])

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
}
