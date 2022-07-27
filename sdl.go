package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type sdlHandle struct {
	surface *sdl.Surface
	window  *sdl.Window
	keyMap  map[int]uint8
}

/*
	map
	1 2 3 4
	Q W E R
	A S D F
	Z X C V

	chip8 kb
	1	2 3 C
	4 5 6 D
	7 8 9 E
	A 0 B F
*/

func newSdl() sdlHandle {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	wnd, err := sdl.CreateWindow("chip8", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	srf, err := wnd.GetSurface()
	if err != nil {
		panic(err)
	}

	return sdlHandle{
		surface: srf,
		window:  wnd,
		keyMap: map[int]uint8{
			sdl.SCANCODE_1: 0x1,
			sdl.SCANCODE_2: 0x2,
			sdl.SCANCODE_3: 0x3,
			sdl.SCANCODE_4: 0xC,
			sdl.SCANCODE_Q: 0x4,
			sdl.SCANCODE_W: 0x5,
			sdl.SCANCODE_E: 0x6,
			sdl.SCANCODE_R: 0xD,
			sdl.SCANCODE_A: 0x7,
			sdl.SCANCODE_S: 0x8,
			sdl.SCANCODE_D: 0x9,
			sdl.SCANCODE_F: 0xE,
			sdl.SCANCODE_Z: 0xA,
			sdl.SCANCODE_X: 0x0,
			sdl.SCANCODE_C: 0xB,
			sdl.SCANCODE_V: 0xF,
		},
	}
}

func (h sdlHandle) cleanup() {
	h.window.Destroy()
	sdl.Quit()
}

func (h sdlHandle) drawWindow(c chip8) {
	for row := 0; row < 32; row++ {
		for col := 0; col < 64; col++ {
			var color uint32
			if c.screen[row][col] {
				color = 0xFFFFFFFF
			}
			rect := sdl.Rect{X: int32(col * 10), Y: int32(row * 10), W: 10, H: 10}
			h.surface.FillRect(&rect, color)
		}
	}
	h.window.UpdateSurface()
}

func (h sdlHandle) getKeyPressed() (uint8, bool) {
	kb := sdl.GetKeyboardState()
	sdl.PumpEvents()

	for k, v := range h.keyMap {
		if kb[k] > 0 {
			return v, true
		}
	}
	return 0xFF, false
}
