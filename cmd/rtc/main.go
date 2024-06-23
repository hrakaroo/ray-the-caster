package main

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const PI = 3.1415926535
const PI2 = PI * 2.0

type Environment struct {
	PixelWidth  int32
	PixelHeight int32
	UnitWidth   int32
	UnitHeight  int32
	xSize       int32
	ySize       int32
	Area        [100]int16
}

func NewEnvironment() *Environment {
	env := &Environment{
		PixelWidth:  600,
		PixelHeight: 600,
		UnitWidth:   10,
		UnitHeight:  10,
		Area: [100]int16{
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

	env.xSize = env.PixelWidth / env.UnitWidth
	env.ySize = env.PixelHeight / env.UnitHeight

	return env
}

// yLoc returns the y screen location given an y index
func (e *Environment) yLoc(yIndex int) int32 {
	return int32(yIndex) * e.ySize
}

// xLoc returns the x screen location given an x index
func (e *Environment) xLoc(xIndex int) int32 {
	return int32(xIndex) * e.xSize
}

func (e *Environment) areaValue(xIndex, yIndex int) int16 {
	if xIndex < 0 || yIndex < 0 || xIndex >= int(e.UnitWidth) || yIndex >= int(e.UnitHeight) {
		return -1
	}
	return e.Area[yIndex*int(e.UnitHeight)+xIndex]
}

func (e *Environment) yIndex(yLoc float64) int {
	return int(int32(yLoc) / e.ySize)
}

func (e *Environment) xIndex(xLoc float64) int {
	return int(int32(xLoc) / e.xSize)
}

type Player struct {
	X      float64
	Y      float64
	Width  int32
	Height int32
	Angle  float64
	DeltaX float64
	DeltaY float64
}

func NewPlayer() *Player {
	direction := 0.0
	directionX := math.Cos(direction)
	directionY := math.Sin(direction)

	return &Player{X: 150, Y: 150, Width: 6, Height: 6,
		Angle: direction, DeltaX: directionX, DeltaY: directionY}
}

func (p *Player) Left() {
	p.Angle -= 0.15
	if p.Angle < 0 {
		p.Angle = PI2
	}
	p.DeltaX = math.Cos(p.Angle)
	p.DeltaY = math.Sin(p.Angle)
}

func (p *Player) Right() {
	p.Angle += 0.15
	if p.Angle > PI2 {
		p.Angle = 0
	}
	p.DeltaX = math.Cos(p.Angle)
	p.DeltaY = math.Sin(p.Angle)
}

func (p *Player) Forward() {
	p.X += p.DeltaX * 4.0
	p.Y += p.DeltaY * 4.0
}

func (p *Player) Backward() {
	p.X -= p.DeltaX * 4.0
	p.Y -= p.DeltaY * 4.0
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	env := NewEnvironment()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, env.PixelWidth, env.PixelHeight, sdl.WINDOW_SHOWN)
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

	render(renderer, env, pl)

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

				render(renderer, env, pl)

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
	xUnit := env.PixelWidth / env.UnitWidth
	yUnit := env.PixelHeight / env.UnitHeight
	renderer.SetDrawColor(200, 200, 200, 255)
	for yIndex := 0; yIndex < int(env.UnitHeight); yIndex++ {
		for xIndex := 0; xIndex < int(env.UnitWidth); xIndex++ {
			value := env.Area[yIndex*int(env.UnitHeight)+xIndex]
			if value == 0 {
				rect := sdl.Rect{X: int32(xIndex * int(xUnit)), Y: int32(yIndex * int(yUnit)), W: xUnit, H: yUnit}
				renderer.DrawRect(&rect)
			}
			if value == 1 {
				// Draw a wall box
				rect := sdl.Rect{X: int32(xIndex*int(xUnit)) + 1, Y: int32(yIndex*int(yUnit)) + 1, W: xUnit - 2, H: yUnit - 2}
				renderer.FillRect(&rect)
			}
		}
	}

	// Draw our player
	renderer.SetDrawColor(0, 255, 0, 255)
	rect := sdl.Rect{X: int32(pl.X) - pl.Width/2, Y: int32(pl.Y) - pl.Height/2, W: pl.Width, H: pl.Height}
	renderer.FillRect(&rect)

	// Draw the direction vector
	// It's a unit vector so multiply it by 15 (selected at random)
	renderer.DrawLine(int32(pl.X), int32(pl.Y), int32(pl.X+pl.DeltaX*15.0), int32(pl.Y+pl.DeltaY*15.0))

	// Determine the xBox the player is in
	xBox := env.xIndex(pl.X)
	yBox := env.yIndex(pl.Y)

	renderer.SetDrawColor(255, 0, 0, 255)
	for i := -0.5; i < 0.5; i += 0.01 {
		xCollision, yCollision, _, _ := detectCollision(float64(pl.X), float64(pl.Y), xBox, yBox, pl.Angle+i, env)

		renderer.DrawLine(int32(pl.X), int32(pl.Y), xCollision, yCollision)
	}

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

func detectCollision(x, y float64, xBox, yBox int, angle float64, env *Environment) (collisionX int32, collisionY int32, collisionDistance int32, collisionAngle float64) {

	var xSave, ySave, distance float64
	distance = 100_000_000

	yMax := int(env.PixelHeight / env.UnitHeight)
	xMax := int(env.PixelWidth / env.UnitWidth)

	tan := math.Tan(float64(angle))

	// Normalize the angle
	for angle > PI2 {
		angle -= PI2
	}

	if angle > 0 && angle < PI {
		// fmt.Println("Positive Y")

		// We want to find the intersection with the line below us which is owned by yBox+1
		for i := yBox + 1; i < yMax; i++ {
			yDelta := float64(env.yLoc(i)) - y
			xDelta := yDelta / tan

			xIndex := env.xIndex(x + xDelta)
			yIndex := env.yIndex(y + yDelta)

			value := env.areaValue(xIndex, yIndex)
			if value == 1 {
				// Compute distance (no need to take sqrt ... yet)
				dist := xDelta*xDelta + yDelta*yDelta

				// HIT, record the xDist
				if dist < distance {
					xSave = x + xDelta
					ySave = y + yDelta
					distance = dist
				}
				break
			}
		}
	} else if angle > PI && angle < PI2 {
		// fmt.Println("Negative Y")

		// We want to find the intersection with the line above us which is owned by yBox
		for i := yBox; i >= 0; i-- {
			yDelta := float64(env.yLoc(i)) - y
			xDelta := yDelta / tan

			xIndex := env.xIndex(x + xDelta)
			yIndex := env.yIndex(y + yDelta)

			// Since yIndex is going to be on the line, the lines in front of us belong to the cubes in front of them.
			yIndex--

			value := env.areaValue(xIndex, yIndex)
			if value == 1 {
				// Compute distance (no need to take sqrt ... yet)
				dist := xDelta*xDelta + yDelta*yDelta

				// HIT, record the xDist
				if dist < distance {
					xSave = x + xDelta
					ySave = y + yDelta
					distance = dist
				}
				break
			}
		}
	}

	if angle > PI/2 && angle < 6*PI/4 {
		// fmt.Println("Negative X")
		// We want to find the intersection with the line before us which is owned by xBox
		for i := xBox; i >= 0; i-- {
			xDelta := float64(env.xLoc(i)) - x

			// Compute the yDelta, subtract 0.0001 to counter roundoff
			yDelta := xDelta * tan

			xIndex := env.xIndex(x + xDelta)
			yIndex := env.yIndex(y + yDelta)
			xIndex--

			value := env.areaValue(xIndex, yIndex)
			if value == 1 {
				// Compute distance (no need to take sqrt ... yet)
				dist := xDelta*xDelta + yDelta*yDelta

				// HIT, record the xDist
				if dist < distance {
					xSave = x + xDelta
					ySave = y + yDelta
					distance = dist
				}
				break
			}
		}
	} else if angle < PI/2 || angle > 6*PI/4 {
		// fmt.Println("Positive X")
		// We want to find the intersection with the line below us which is owned by yBox+1
		for i := xBox + 1; i < xMax; i++ {
			xDelta := float64(env.yLoc(i)) - x
			yDelta := xDelta * tan

			xIndex := env.xIndex(x + xDelta)
			yIndex := env.yIndex(y + yDelta)

			value := env.areaValue(xIndex, yIndex)
			if value == 1 {
				// Compute distance (no need to take sqrt ... yet)
				dist := xDelta*xDelta + yDelta*yDelta

				// HIT, record the xDist
				if dist < distance {
					xSave = x + xDelta
					ySave = y + yDelta
					distance = dist
				}
				break
			}
		}
	}
	return int32(xSave), int32(ySave), int32(math.Sqrt(float64(distance))), 0
}
