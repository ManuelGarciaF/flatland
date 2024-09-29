package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const WINDOW_WIDTH = 1280
const WINDOW_HEIGHT = 250

var BACKGROUND_COLOR = rl.Black

const VIEW_DISTANCE = 20

type World []WorldObject

var world World = []WorldObject{
	&Polygon{
		Points: []rl.Vector2{
			{X: 2, Y: 8},
			{X: 2, Y: 16},
			{X: -2, Y: 16},
			{X: -2, Y: 8},
		},
		color: rl.Red,
	},

	&Polygon{
		Points: []rl.Vector2{
			{X: 7, Y: 8},
			{X: 7, Y: 16},
			{X: 6, Y: 16},
			{X: 6, Y: 8},
		},
		color: rl.Blue,
	},

	// Front facing triangle
	&Polygon{
		Points: []rl.Vector2{
			{X: 10, Y: 8},
			{X: 15, Y: 20},
			{X: 15, Y: 8},
		},
		color: rl.Green,
	},
}

func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Flatland")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	c := NewCamera(rl.Vector2{X: 0, Y: 0}, rl.Vector2{X: 0, Y: 1}, &world)

	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	rl.EndDrawing()
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		colors := c.CastRays()

		if rl.IsKeyDown(rl.KeyW) {
			c.Position = rl.Vector2Add(c.Position, rl.Vector2{X: 0, Y: 0.1})
		}
		if rl.IsKeyDown(rl.KeyS) {
			c.Position = rl.Vector2Add(c.Position, rl.Vector2{X: 0, Y: -0.1})
		}
		if rl.IsKeyDown(rl.KeyE) {
			c.Direction = rl.Vector2Rotate(c.Direction, -0.05)
		}
		if rl.IsKeyDown(rl.KeyQ) {
			c.Direction = rl.Vector2Rotate(c.Direction, 0.05)
		}

		for pX, color := range colors {
			rl.DrawLine(int32(pX), 0, int32(pX), WINDOW_HEIGHT, color)
		}

		// Print position
		rl.DrawText(
			fmt.Sprintf(
				"Pos: X=%v, Y=%v; Dir: X=%v, Y=%v",
				c.Position.X,
				c.Position.Y,
				c.Direction.X,
				c.Direction.Y,
			),
			10,
			10,
			10,
			rl.RayWhite)

		rl.EndDrawing()
	}
}
