package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type environment struct {
	width  int32
	height int32
	area   [100]int16
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	env := environment{
		width:  600,
		height: 600,
		area: [100]int16{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 0, 0, 1, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 1, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 1, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 1, 1, 0, 0, 1,
			1, 0, 0, 0, 0, 1, 1, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, env.width, env.height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Failed to create renderer: %s\n", err)
		return
	}
	defer renderer.Destroy()

	player := sdl.Rect{X: 150, Y: 150, W: 5, H: 5}

	render(renderer, &env, &player)

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
					player.X += 5
				case sdl.K_LEFT:
					// fmt.Println("Left")
					player.X -= 5
				case sdl.K_UP:
					// fmt.Println("Forward")
					player.Y -= 5
				case sdl.K_DOWN:
					// fmt.Println("Backward")
					player.Y += 5
				}

				render(renderer, &env, &player)

			}
		}

		sdl.Delay(33)
	}
}

func render(renderer *sdl.Renderer, env *environment, player *sdl.Rect) {

	// I'm nearly positive this is not the most efficient way to do this, but I'm not super
	//  interested in hyper optimizing the actual drawing. SDL is just a tool here

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	// Draw the map
	xUnit := env.width / 10
	yUnit := env.height / 10
	renderer.SetDrawColor(255, 255, 255, 255)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			value := env.area[y*10+x]
			if value == 1 {
				// Draw a wall box
				rect := sdl.Rect{X: int32(x*int(xUnit)) + 1, Y: int32(y*int(yUnit)) + 1, W: xUnit - 2, H: yUnit - 2}
				renderer.FillRect(&rect)
			}
		}
	}

	renderer.SetDrawColor(0, 255, 0, 255)
	renderer.FillRect(player)

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
