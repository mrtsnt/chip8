package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type sdlHandle struct {
	surface *sdl.Surface
	window  *sdl.Window
}

func newSdl() sdlHandle {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	wnd, err := sdl.CreateWindow("chip8", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 1280, 640, sdl.WINDOW_SHOWN)
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
	}
}

func (h sdlHandle) cleanup() {
	h.window.Destroy()
	sdl.Quit()
}

func (h sdlHandle) drawWindow(c *chip8) {
	pixelSize := 20
	if c.isHighRes {
		pixelSize = 10
	}
	for p := 0; p < len(c.screen); p++ {
		var color uint32
		if c.screen[p] {
			color = 0xFFFFFFFF
		}
		rect := sdl.Rect{X: int32(p % c.xLen * pixelSize), Y: int32(p / c.xLen * pixelSize), W: int32(pixelSize), H: int32(pixelSize)}
		h.surface.FillRect(&rect, color)
	}
	h.window.UpdateSurface()
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

func (h sdlHandle) updateKeyboard(c *chip8) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		kbe, ok := event.(*sdl.KeyboardEvent)
		if ok {
			toSet := false
			if kbe.State == sdl.PRESSED {
				toSet = true
			}

			switch kbe.Keysym.Scancode {
			case sdl.SCANCODE_1:
				c.keys[0x1] = toSet
			case sdl.SCANCODE_2:
				c.keys[0x2] = toSet
			case sdl.SCANCODE_3:
				c.keys[0x3] = toSet
			case sdl.SCANCODE_4:
				c.keys[0xC] = toSet
			case sdl.SCANCODE_Q:
				c.keys[0x4] = toSet
			case sdl.SCANCODE_W:
				c.keys[0x5] = toSet
			case sdl.SCANCODE_E:
				c.keys[0x6] = toSet
			case sdl.SCANCODE_R:
				c.keys[0xD] = toSet
			case sdl.SCANCODE_A:
				c.keys[0x7] = toSet
			case sdl.SCANCODE_S:
				c.keys[0x8] = toSet
			case sdl.SCANCODE_D:
				c.keys[0x9] = toSet
			case sdl.SCANCODE_F:
				c.keys[0xE] = toSet
			case sdl.SCANCODE_Z:
				c.keys[0xA] = toSet
			case sdl.SCANCODE_X:
				c.keys[0x0] = toSet
			case sdl.SCANCODE_C:
				c.keys[0xB] = toSet
			case sdl.SCANCODE_V:
				c.keys[0xF] = toSet
			}
		}
	}
}
