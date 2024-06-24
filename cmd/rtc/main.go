package main

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const FieldOfViewDegrees = 30
const IncrementsDegrees = .5

const PI = 3.1415926535
const PI2 = PI * 2.0
const FieldOfViewRads = PI * FieldOfViewDegrees / 180
const IncrementsRads = PI * IncrementsDegrees / 180

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

var colors map[uint8]Color

func init() {
	colors = make(map[uint8]Color)
	colors[1] = Color{R: 200, G: 200, B: 200, A: 255}
	colors[2] = Color{R: 0, G: 200, B: 0, A: 255}
	colors[3] = Color{R: 0, G: 0, B: 200, A: 255}
}

type Environment struct {
	PixelWidth  int32
	PixelHeight int32
	UnitWidth   int32
	UnitHeight  int32
	xSize       int32
	ySize       int32
	Area        [256]int8
}

func NewEnvironment() *Environment {
	unitWidth := 16
	unitHeight := 16
	xSize := 40
	ySize := 40

	env := &Environment{
		PixelWidth:  int32(unitWidth * xSize),
		PixelHeight: int32(unitHeight * ySize),
		UnitWidth:   int32(unitWidth),
		UnitHeight:  int32(unitHeight),
		xSize:       int32(xSize),
		ySize:       int32(ySize),
		Area: [256]int8{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 2, 2, 2, 2, 2, 2, 2, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 2, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 2, 3, 0, 0, 2, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 2, 3, 3, 0, 2, 2, 2, 2, 0, 0, 1,
			1, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

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

func (e *Environment) areaValue(xIndex, yIndex int) int8 {
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

	return &Player{X: 75, Y: 75, Width: 6, Height: 6,
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

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 2*env.PixelWidth, env.PixelHeight, sdl.WINDOW_SHOWN)
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
	renderer.SetDrawColor(100, 100, 100, 255)
	renderer.Clear()

	renderer.SetViewport(&sdl.Rect{X: 0, Y: 0, W: env.PixelWidth, H: env.PixelHeight})
	render2D(renderer, env, pl)

	renderer.SetViewport(&sdl.Rect{X: env.PixelWidth + 1, Y: 0, W: env.PixelHeight, H: env.PixelHeight})
	render3D(renderer, env, pl)

	renderer.Present()
}

func render3D(renderer *sdl.Renderer, env *Environment, pl *Player) {

	// Determine the xBox the player is in
	xBox := env.xIndex(pl.X)
	yBox := env.yIndex(pl.Y)

	// 30 degrees of view
	width := float64(env.PixelWidth) * IncrementsDegrees / FieldOfViewDegrees

	x := 0
	for offset := -FieldOfViewRads / 2; offset < FieldOfViewRads/2; offset += IncrementsRads {
		_, _, distance, color, angle := detectCollision(float64(pl.X), float64(pl.Y), xBox, yBox, pl.Angle+offset, env)

		renderer.SetDrawColor(uint8(float64(color.R)*angle), uint8(float64(color.G)*angle), uint8(float64(color.B)*angle), color.A)

		height := float64(env.PixelHeight) / math.Log10(distance)

		y := (float64(env.PixelHeight) - height) / 2.0
		renderer.FillRect(&sdl.Rect{X: int32(x), Y: int32(y), W: int32(width), H: int32(height)})

		x += int(width)
	}
}

func render2D(renderer *sdl.Renderer, env *Environment, pl *Player) {
	// Draw the map
	xUnit := env.PixelWidth / env.UnitWidth
	yUnit := env.PixelHeight / env.UnitHeight
	for yIndex := 0; yIndex < int(env.UnitHeight); yIndex++ {
		for xIndex := 0; xIndex < int(env.UnitWidth); xIndex++ {
			value := env.Area[yIndex*int(env.UnitHeight)+xIndex]
			color := colors[uint8(value)]

			if value == 0 {
				renderer.SetDrawColor(200, 200, 200, 255)
				rect := sdl.Rect{X: int32(xIndex * int(xUnit)), Y: int32(yIndex * int(yUnit)), W: xUnit, H: yUnit}
				renderer.DrawRect(&rect)
			} else {
				renderer.SetDrawColor(color.R, color.G, color.B, color.A)
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

	// Determine the xBox the player is in
	xBox := env.xIndex(pl.X)
	yBox := env.yIndex(pl.Y)

	renderer.SetDrawColor(255, 0, 0, 255)

	for offset := -FieldOfViewRads / 2; offset < FieldOfViewRads/2; offset += IncrementsRads {
		xCollision, yCollision, _, _, _ := detectCollision(float64(pl.X), float64(pl.Y), xBox, yBox, pl.Angle+offset, env)
		renderer.DrawLine(int32(pl.X), int32(pl.Y), xCollision, yCollision)
	}
}

func detectCollision(x, y float64, xBox, yBox int, angle float64, env *Environment) (collisionX, collisionY int32, collisionDistance float64, color Color, collisionAngle float64) {

	var xSave, ySave, angleSave, distance float64
	var colorSave Color
	distance = 100_000_000

	yMax := int(env.PixelHeight / env.UnitHeight)
	xMax := int(env.PixelWidth / env.UnitWidth)

	tan := math.Tan(float64(angle))

	// Normalize the angle
	for angle > PI2 {
		angle -= PI2
	}
	for angle < 0 {
		angle += PI2
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
			if value != 0 {
				// Compute distance (no need to take sqrt ... yet)
				dist := xDelta*xDelta + yDelta*yDelta

				// HIT, record the xDist
				if dist < distance {
					xSave = x + xDelta
					ySave = y + yDelta
					angleSave = math.Sin(angle)
					colorSave = colors[uint8(value)]
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
			if value != 0 {
				// Compute distance (no need to take sqrt ... yet)
				dist := xDelta*xDelta + yDelta*yDelta

				// HIT, record the xDist
				if dist < distance {
					xSave = x + xDelta
					ySave = y + yDelta
					angleSave = -1 * math.Sin(angle)
					colorSave = colors[uint8(value)]
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
			if value != 0 {
				// Compute distance (no need to take sqrt ... yet)
				dist := xDelta*xDelta + yDelta*yDelta

				// HIT, record the xDist
				if dist < distance {
					xSave = x + xDelta
					ySave = y + yDelta
					angleSave = -1 * math.Cos(angle)
					colorSave = colors[uint8(value)]
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
			if value != 0 {
				// Compute distance (no need to take sqrt ... yet)
				dist := xDelta*xDelta + yDelta*yDelta

				// HIT, record the xDist
				if dist < distance {
					xSave = x + xDelta
					ySave = y + yDelta
					angleSave = math.Cos(angle)
					colorSave = colors[uint8(value)]
					distance = dist
				}
				break
			}
		}
	}

	return int32(xSave), int32(ySave), math.Sqrt(distance), colorSave, angleSave
}
