package main

import (
	"fmt"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)


const TOP_DOWN_VIEW_HEIGHT = 400
const FIRST_PERSON_VIEW_HEIGHT = 250
const FIRST_PERSON_VIEW_START_Y = WINDOW_HEIGHT - FIRST_PERSON_VIEW_HEIGHT
const TOPDOWN_SCALE = 3


const WINDOW_WIDTH = 1280
const WINDOW_HEIGHT = FIRST_PERSON_VIEW_HEIGHT + TOP_DOWN_VIEW_HEIGHT + 5

var BACKGROUND_COLOR = rl.Black

const VIEW_DISTANCE = 50
const RANDOM_OBJECTS = 40

const WORLD_SIZE = 75

type World []WorldObject

func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Flatland")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// Make a random world
	world := makeRandomWorld()

	c := NewCamera(rl.Vector2{X: 0, Y: 0}, rl.Vector2{X: 0, Y: 1}, &world)

	rl.BeginDrawing()
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(BACKGROUND_COLOR)

		handleMovement(c)

		// Draw 2D world
		drawTopDownView(c, &world)

		// Draw divider
		rl.DrawLineEx(
			rl.Vector2{X: 0, Y: FIRST_PERSON_VIEW_START_Y - 3},
			rl.Vector2{X: WINDOW_WIDTH, Y: FIRST_PERSON_VIEW_START_Y - 3},
			3,
			rl.RayWhite,
		)

		// Draw the results of rays
		colors := c.CastRays()
		for pX, color := range colors {
			rl.DrawLine(int32(pX), FIRST_PERSON_VIEW_START_Y, int32(pX), WINDOW_HEIGHT, color)
		}

		printInfo(c)

		rl.EndDrawing()
	}
}

func makeRandomWorld() World {
	world := make(World, 0, RANDOM_OBJECTS)
	possibleColors := []rl.Color{rl.Blue, rl.Yellow, rl.Red, rl.Purple, rl.Green}
	for i := 0; i < RANDOM_OBJECTS; i++ {
		pol := NewPolygon(
			rl.Vector2{
				X: float32(rand.Intn(150)) - 75.0,
				Y: float32(rand.Intn(150)) - 75.0,
			},
			rand.Intn(6)+2, // Make sure they are at least lines
			rand.Float32()*2.0*math.Pi,
			float32(rand.Intn(8))+3.0,
			possibleColors[rand.Intn(len(possibleColors))],
		)
		world = append(world, &pol)
	}
	return world
}

func drawTopDownView(c *Camera, world *World) {
	center := rl.Vector2{X: WINDOW_WIDTH / 2, Y: TOP_DOWN_VIEW_HEIGHT / 2}
	for _, obj := range *world {
		switch obj := obj.(type) {
		case *Polygon:
			for currPoint := range obj.Points {
				nextPoint := (currPoint + 1) % len(obj.Points)
				rl.DrawLineV(
					rl.Vector2Add(center, rl.Vector2Scale(obj.Points[currPoint], TOPDOWN_SCALE)),
					rl.Vector2Add(center, rl.Vector2Scale(obj.Points[nextPoint], TOPDOWN_SCALE)),
					obj.Color(),
				)
			}
		}
	}

	rl.DrawCircleV(rl.Vector2Add(center, rl.Vector2Scale(c.Position, TOPDOWN_SCALE)), 3, rl.RayWhite)
	rl.DrawLineEx(
		rl.Vector2Add(center, rl.Vector2Scale(c.Position, TOPDOWN_SCALE)),
		rl.Vector2Add(
			rl.Vector2Add(center, rl.Vector2Scale(c.Position, TOPDOWN_SCALE)),
			rl.Vector2Scale(c.Direction, TOPDOWN_SCALE*2),
		),
		2,
		rl.RayWhite,
	)

}

func handleMovement(c *Camera) {
	if rl.IsKeyDown(rl.KeyW) {
		c.Position = rl.Vector2Add(
			c.Position,
			rl.Vector2Scale(c.Direction, 0.15),
		)
	}
	if rl.IsKeyDown(rl.KeyS) {
		c.Position = rl.Vector2Add(
			c.Position,
			rl.Vector2Scale(c.Direction, -0.15),
		)
	}
	perpDirection := rl.Vector2{
		X: c.Direction.Y,
		Y: -c.Direction.X,
	}
	if rl.IsKeyDown(rl.KeyA) {
		c.Position = rl.Vector2Add(
			c.Position,
			rl.Vector2Scale(perpDirection, 0.15),
		)
	}
	if rl.IsKeyDown(rl.KeyD) {
		c.Position = rl.Vector2Add(
			c.Position,
			rl.Vector2Scale(perpDirection, -0.15),
		)
	}
	if rl.IsKeyDown(rl.KeyQ) {
		c.Direction = rl.Vector2Rotate(c.Direction, -0.05)
	}
	if rl.IsKeyDown(rl.KeyE) {
		c.Direction = rl.Vector2Rotate(c.Direction, 0.05)
	}

	// Clamp position
	if c.Position.X < -WORLD_SIZE {
		c.Position.X = -WORLD_SIZE
	}
	if c.Position.X > WORLD_SIZE {
		c.Position.X = WORLD_SIZE
	}
	if c.Position.Y < -WORLD_SIZE {
		c.Position.Y = -WORLD_SIZE
	}
	if c.Position.Y > WORLD_SIZE {
		c.Position.Y = WORLD_SIZE
	}
}

func printInfo(c *Camera) {
	rl.DrawFPS(WINDOW_WIDTH-100, 10)

	// Window titles
	rl.DrawText("Top Down View", 10, 10, 20, rl.RayWhite)
	rl.DrawText("First Person View", 10, FIRST_PERSON_VIEW_START_Y+10, 20, rl.RayWhite)

	// Controls
	rl.DrawText(
		"WASD to move, Q and E to rotate camera",
		10,
		FIRST_PERSON_VIEW_START_Y-30,
		20,
		rl.RayWhite,
	)

	// Print position
	rl.DrawText(
		fmt.Sprintf("Pos: X=%.3f, Y=%.3f", c.Position.X, c.Position.Y),
		10,
		FIRST_PERSON_VIEW_START_Y-50,
		10,
		rl.RayWhite,
	)
}
