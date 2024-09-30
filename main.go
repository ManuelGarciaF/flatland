package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const WINDOW_WIDTH = 1280

// const WINDOW_HEIGHT = 250
const WINDOW_HEIGHT = 250 + 300
const FIRST_PERSON_VIEW_HEIGHT = 250
const FIRST_PERSON_VIEW_START_Y = WINDOW_HEIGHT - FIRST_PERSON_VIEW_HEIGHT
const TOPDOWN_SCALE = 3

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
			{X: 10, Y: 20},
			{X: 15, Y: 8},
			{X: 20, Y: 20},
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
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(BACKGROUND_COLOR)

		handleMovement(c)

		// Draw 2D world
		drawFirstPersonView(c)

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

func drawFirstPersonView(c *Camera) {
	center := rl.Vector2{X: WINDOW_WIDTH / 2, Y: FIRST_PERSON_VIEW_HEIGHT / 2}
	for _, obj := range world {
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

	rl.DrawCircleV(rl.Vector2Add(center, rl.Vector2Scale(c.Position, TOPDOWN_SCALE)), 5, rl.RayWhite)

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
			rl.Vector2Scale(perpDirection, -0.15),
		)
	}
	if rl.IsKeyDown(rl.KeyD) {
		c.Position = rl.Vector2Add(
			c.Position,
			rl.Vector2Scale(perpDirection, 0.15),
		)
	}
	if rl.IsKeyDown(rl.KeyQ) {
		c.Direction = rl.Vector2Rotate(c.Direction, 0.05)
	}
	if rl.IsKeyDown(rl.KeyE) {
		c.Direction = rl.Vector2Rotate(c.Direction, -0.05)
	}
}
