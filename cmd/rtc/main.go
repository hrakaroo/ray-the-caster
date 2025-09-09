package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

// Raycasting configuration constants
const FieldOfViewDegrees = 30    // Field of view in degrees
const IncrementDegrees = 0.5     // Angular resolution for ray casting

// Mathematical constants
const PI = 3.1415926535
const PI2 = PI * 2.0
const FieldOfViewRadians = PI * FieldOfViewDegrees / 180
const IncrementRadians = PI * IncrementDegrees / 180

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
	MapWidth    int32  // Number of map units horizontally
	MapHeight   int32  // Number of map units vertically
	UnitWidth   int32  // Pixel width of each map unit
	UnitHeight  int32  // Pixel height of each map unit
	Area        []int8 // Map data (dynamically sized)
}

func NewEnvironment() *Environment {
	unitPixelWidth := 16
	unitPixelHeight := 16
	mapWidth := 16
	mapHeight := 16

	env := &Environment{
		PixelWidth:  int32(unitPixelWidth * mapWidth),
		PixelHeight: int32(unitPixelHeight * mapHeight),
		UnitWidth:   int32(unitPixelWidth),
		UnitHeight:  int32(unitPixelHeight),
		MapWidth:    int32(mapWidth),
		MapHeight:   int32(mapHeight),
		Area: []int8{
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

// LoadEnvironmentFromFile loads a map from a text file
// File format: first line contains "width height", followed by rows of space-separated integers
func LoadEnvironmentFromFile(filename string) (*Environment, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open map file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	
	// Read first line to get dimensions
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty map file")
	}
	
	dimensions := strings.Fields(scanner.Text())
	if len(dimensions) != 2 {
		return nil, fmt.Errorf("invalid header format, expected 'width height'")
	}
	
	mapWidth, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, fmt.Errorf("invalid width: %w", err)
	}
	
	mapHeight, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, fmt.Errorf("invalid height: %w", err)
	}
	
	unitPixelWidth := 40  // Configurable unit size
	unitPixelHeight := 40
	
	// Create environment
	env := &Environment{
		PixelWidth:  int32(unitPixelWidth * mapWidth),
		PixelHeight: int32(unitPixelHeight * mapHeight),
		UnitWidth:   int32(unitPixelWidth),
		UnitHeight:  int32(unitPixelHeight),
		MapWidth:    int32(mapWidth),
		MapHeight:   int32(mapHeight),
		Area:        make([]int8, mapWidth*mapHeight),
	}
	
	// Read map data
	row := 0
	for scanner.Scan() && row < mapHeight {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		values := strings.Fields(line)
		if len(values) != mapWidth {
			return nil, fmt.Errorf("row %d has %d values, expected %d", row, len(values), mapWidth)
		}
		
		for col, valueStr := range values {
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return nil, fmt.Errorf("invalid value at row %d, col %d: %w", row, col, err)
			}
			
			env.Area[env.getMapIndex(col, row)] = int8(value)
		}
		
		row++
	}
	
	if row != mapHeight {
		return nil, fmt.Errorf("expected %d rows, got %d", mapHeight, row)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	
	return env, nil
}

// getYPixelPosition returns the y screen location given a y map index
func (e *Environment) getYPixelPosition(yIndex int) int32 {
	return int32(yIndex) * e.UnitHeight
}

// getXPixelPosition returns the x screen location given an x map index
func (e *Environment) getXPixelPosition(xIndex int) int32 {
	return int32(xIndex) * e.UnitWidth
}

// getMapValue returns the map value at given coordinates, or -1 if out of bounds
func (e *Environment) getMapValue(xIndex, yIndex int) int8 {
	if xIndex < 0 || yIndex < 0 || xIndex >= int(e.MapWidth) || yIndex >= int(e.MapHeight) {
		return -1
	}
	return e.Area[e.getMapIndex(xIndex, yIndex)]
}

// getYMapIndex converts y pixel coordinate to map index
func (e *Environment) getYMapIndex(yPixel float64) int {
	return int(int32(yPixel) / e.UnitHeight)
}

// getXMapIndex converts x pixel coordinate to map index
func (e *Environment) getXMapIndex(xPixel float64) int {
	return int(int32(xPixel) / e.UnitWidth)
}

// getMapIndex converts 2D map coordinates to 1D array index
func (e *Environment) getMapIndex(xIndex, yIndex int) int {
	return yIndex*int(e.MapWidth) + xIndex
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

	// Try to load map from file, fallback to default if not found
	var env *Environment
	var err error
	
	if len(os.Args) > 1 {
		// Map file specified as command line argument
		env, err = LoadEnvironmentFromFile(os.Args[1])
		if err != nil {
			fmt.Printf("Failed to load map file '%s': %v\n", os.Args[1], err)
			fmt.Println("Using default map instead.")
			env = NewEnvironment()
		}
	} else {
		// Try to load default map file
		env, err = LoadEnvironmentFromFile("map.txt")
		if err != nil {
			// Use hardcoded default map
			env = NewEnvironment()
		}
	}

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

	// Determine the map coordinates the player is in
	xBox := env.getXMapIndex(pl.X)
	yBox := env.getYMapIndex(pl.Y)

	renderer.SetViewport(&sdl.Rect{X: 0, Y: 0, W: env.PixelWidth, H: env.PixelHeight})
	render2D(renderer, env, pl, xBox, yBox)

	renderer.SetViewport(&sdl.Rect{X: env.PixelWidth + 1, Y: 0, W: env.PixelHeight, H: env.PixelHeight})
	render3D(renderer, env, pl, xBox, yBox)

	renderer.Present()
}

func render3D(renderer *sdl.Renderer, env *Environment, pl *Player, xBox, yBox int) {

	// Calculate the width of each ray column based on field of view
	width := float64(env.PixelWidth) * IncrementDegrees / FieldOfViewDegrees

	x := 0
	// Cast rays across the field of view from left to right
	for offset := -FieldOfViewRadians / 2; offset < FieldOfViewRadians/2; offset += IncrementRadians {
		_, _, distance, color, wallAngle := detectCollision(float64(pl.X), float64(pl.Y), xBox, yBox, pl.Angle+offset, env)

		// Apply angle-based shading (walls at sharper angles appear darker)
		wallAngle = math.Sqrt(wallAngle)

		// Modulate color intensity based on wall angle
		color.R = uint8(float64(color.R) * wallAngle)
		color.G = uint8(float64(color.G) * wallAngle)
		color.B = uint8(float64(color.B) * wallAngle)

		// Apply distance-based shading (further objects appear darker)
		delta := float64(env.PixelWidth) - distance/2
		if delta < 0 {
			// Object is too far away, skip rendering
			continue
		}

		// Calculate brightness ratio based on distance
		brightness := delta / float64(env.PixelWidth)

		// Apply distance-based darkening
		color.R = uint8(float64(color.R) * brightness)
		color.G = uint8(float64(color.G) * brightness)
		color.B = uint8(float64(color.B) * brightness)

		renderer.SetDrawColor(color.R, color.G, color.B, color.A)

		// Fix fish-eye distortion by correcting for viewing angle
		distance *= math.Cos(offset)

		// Calculate wall height on screen (closer walls appear taller)
		height := float64(env.PixelHeight) / math.Log10(distance)

		y := (float64(env.PixelHeight) - height) / 2.0
		renderer.FillRect(&sdl.Rect{X: int32(x), Y: int32(y), W: int32(width), H: int32(height)})

		x += int(width)
	}
}

func render2D(renderer *sdl.Renderer, env *Environment, pl *Player, xBox, yBox int) {
	// Draw the map
	xUnit := env.PixelWidth / env.MapWidth
	yUnit := env.PixelHeight / env.MapHeight
	for yIndex := 0; yIndex < int(env.MapHeight); yIndex++ {
		for xIndex := 0; xIndex < int(env.MapWidth); xIndex++ {
			value := env.getMapValue(xIndex, yIndex)
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

	renderer.SetDrawColor(255, 0, 0, 255)

	// Cast rays for 2D visualization
	for offset := -FieldOfViewRadians / 2; offset < FieldOfViewRadians/2; offset += IncrementRadians {
		xCollision, yCollision, _, _, _ := detectCollision(float64(pl.X), float64(pl.Y), xBox, yBox, pl.Angle+offset, env)
		renderer.DrawLine(int32(pl.X), int32(pl.Y), xCollision, yCollision)
	}
}

// detectCollision performs raycasting to find the nearest wall intersection
// Returns collision coordinates, distance, wall color, and angle for shading
func detectCollision(x, y float64, xBox, yBox int, angle float64, env *Environment) (collisionX, collisionY int32, collisionDistance float64, color Color, collisionAngle float64) {

	var xSave, ySave, angleSave, distance float64
	var colorSave Color
	distance = 100_000_000 // Initialize to very large distance

	yMax := int(env.MapHeight)
	xMax := int(env.MapWidth)

	tan := math.Tan(float64(angle))

	// Normalize the angle to [0, 2π] range
	for angle > PI2 {
		angle -= PI2
	}
	for angle < 0 {
		angle += PI2
	}

	// Cast ray in positive Y direction (downward)
	if angle > 0 && angle < PI {
		// Check horizontal grid lines below the player
		for i := yBox + 1; i < yMax; i++ {
			yDelta := float64(env.getYPixelPosition(i)) - y
			xDelta := yDelta / tan

			xIndex := env.getXMapIndex(x + xDelta)
			yIndex := env.getYMapIndex(y + yDelta)

			value := env.getMapValue(xIndex, yIndex)
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
	// Cast ray in negative Y direction (upward)
	} else if angle > PI && angle < PI2 {
		// Check horizontal grid lines above the player
		for i := yBox; i >= 0; i-- {
			yDelta := float64(env.getYPixelPosition(i)) - y
			xDelta := yDelta / tan

			xIndex := env.getXMapIndex(x + xDelta)
			yIndex := env.getYMapIndex(y + yDelta)

			// Since yIndex is going to be on the line, the lines in front of us belong to the cubes in front of them.
			yIndex--

			value := env.getMapValue(xIndex, yIndex)
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

	// Cast ray in negative X direction (leftward)
	if angle > PI/2 && angle < 6*PI/4 {
		// Check vertical grid lines to the left of the player
		for i := xBox; i >= 0; i-- {
			xDelta := float64(env.getXPixelPosition(i)) - x

			// Compute the yDelta
			yDelta := xDelta * tan

			xIndex := env.getXMapIndex(x + xDelta)
			yIndex := env.getYMapIndex(y + yDelta)
			xIndex--

			value := env.getMapValue(xIndex, yIndex)
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
	// Cast ray in positive X direction (rightward)
	} else if angle < PI/2 || angle > 6*PI/4 {
		// Check vertical grid lines to the right of the player
		for i := xBox + 1; i < xMax; i++ {
			xDelta := float64(env.getXPixelPosition(i)) - x
			yDelta := xDelta * tan

			xIndex := env.getXMapIndex(x + xDelta)
			yIndex := env.getYMapIndex(y + yDelta)

			value := env.getMapValue(xIndex, yIndex)
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
