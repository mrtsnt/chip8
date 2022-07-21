package main

import (
	"github.com/veandco/go-sdl2/sdl"
)


type sdlHandle struct {
	surface *sdl.Surface
	window *sdl.Window
}

func newSdl() *sdlHandle {
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

	return &sdlHandle{surface: srf, window: wnd}
}

func (h *sdlHandle) cleanup() {
	h.window.Destroy()
	sdl.Quit()
}

func (h *sdlHandle) drawWindow(c *chip8) {
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
