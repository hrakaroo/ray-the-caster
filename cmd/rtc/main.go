package main

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const PI = 3.14159265
const PI2 = PI * 2

type Environment struct {
	width  int32
	height int32
	area   [100]int16
}

type Player struct {
	x          int32
	y          int32
	width      int32
	height     int32
	direction  float32
	directionX float32
	directionY float32
}

func NewPlayer() *Player {
	direction := 0.0
	directionX := float32(math.Cos(direction))
	directionY := float32(math.Sin(direction))

	return &Player{x: 150, y: 150, width: 6, height: 6,
		direction: float32(direction), directionX: directionX, directionY: directionY}
}

func (p *Player) Left() {
	p.direction -= 0.15
	if p.direction < 0 {
		p.direction = PI2
	}
	p.directionX = float32(math.Cos(float64(p.direction)))
	p.directionY = float32(math.Sin(float64(p.direction)))
}

func (p *Player) Right() {
	p.direction += 0.15
	if p.direction > PI2 {
		p.direction = 0
	}
	p.directionX = float32(math.Cos(float64(p.direction)))
	p.directionY = float32(math.Sin(float64(p.direction)))
}

func (p *Player) Forward() {
	p.x += int32(p.directionX * 4)
	p.y += int32(p.directionY * 4)
}

func (p *Player) Backward() {
	p.x -= int32(p.directionX * 4)
	p.y -= int32(p.directionY * 4)
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	env := Environment{
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

	pl := NewPlayer()

	render(renderer, &env, pl)

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
					pl.Right()
				case sdl.K_LEFT:
					pl.Left()
				case sdl.K_UP:
					pl.Forward()
				case sdl.K_DOWN:
					pl.Backward()
				}

				render(renderer, &env, pl)

			}
		}

		sdl.Delay(33)
	}
}

func render(renderer *sdl.Renderer, env *Environment, pl *Player) {

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

	// Draw our player
	renderer.SetDrawColor(0, 255, 0, 255)
	rect := sdl.Rect{X: pl.x - pl.width/2, Y: pl.y - pl.height/2, W: pl.width, H: pl.height}
	renderer.FillRect(&rect)

	// Draw the direction vector
	renderer.DrawLine(pl.x, pl.y, pl.x+int32(int32(pl.directionX*15.0)), pl.y+int32(pl.directionY*15.0))

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
