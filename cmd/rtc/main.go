package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	//  SDL_Renderer * renderer = SDL_CreateRenderer(window, -1, 0);

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Failed to create renderer: %s\n", err)
		return
	}
	defer renderer.Destroy()

	rect := sdl.Rect{X: 10, Y: 10, W: 5, H: 5}

	render(renderer, &rect)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym

				if t.State == sdl.RELEASED {
					// Ignore
					break
				}

				switch keyCode {
				case sdl.K_RIGHT:
					// fmt.Println("Right")
					rect.X += 5
				case sdl.K_LEFT:
					// fmt.Println("Left")
					rect.X -= 5
				case sdl.K_UP:
					// fmt.Println("Forward")
					rect.Y -= 5
				case sdl.K_DOWN:
					// fmt.Println("Backward")
					rect.Y += 5
				}

				render(renderer, &rect)

			}
		}

		sdl.Delay(33)
	}
}

func render(renderer *sdl.Renderer, rect *sdl.Rect) {

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	renderer.SetDrawColor(0, 255, 0, 255)
	renderer.FillRect(rect)

	// renderer.SetDrawColor(255, 255, 255, 255)
	// renderer.DrawPoint(150, 300)

	// renderer.SetDrawColor(0, 0, 255, 255)
	// renderer.DrawLine(0, 0, 200, 200)

	// points := []sdl.Point{{0, 0}, {100, 300}, {100, 300}, {200, 0}}
	// renderer.SetDrawColor(255, 255, 0, 255)
	// renderer.DrawLines(points)

	// rect := sdl.Rect{300, 0, 200, 200}
	// renderer.SetDrawColor(255, 0, 0, 255)
	// renderer.DrawRect(rect)

	// rects := []sdl.Rect{{400, 400, 100, 100}, {550, 350, 200, 200}}
	// renderer.SetDrawColor(0, 255, 255, 255)
	// renderer.DrawRects(rects)

	// rect = sdl.Rect{250, 250, 200, 200}

	// rects = []sdl.Rect{{500, 300, 100, 100}, {200, 300, 200, 200}}
	// renderer.SetDrawColor(255, 0, 255, 255)
	// renderer.FillRects(rects)

	renderer.Present()
	// sdl.Delay(16)
}
