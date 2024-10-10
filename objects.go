package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type WorldObject interface {
	HitBy(ray Ray) (bool, float32) // Returns if it hits the object, then the distance to it.
	Color() rl.Color
	Center() rl.Vector2
}

// Implements WorldObject
type Polygon struct {
	Points []rl.Vector2
	color  rl.Color
}

func NewPolygon(center rl.Vector2, sides int, angle, radius float32, color rl.Color) Polygon {
	points := make([]rl.Vector2, 0, sides)
	for i := 0; i < sides; i++ {
		pointAngle := (float32(i)*2.0*math.Pi)/float32(sides) + angle
		offset := rl.Vector2Scale(
			rl.Vector2Rotate(rl.Vector2{X: 1, Y: 0}, pointAngle),
			radius,
		)
		point := rl.Vector2Add(center, offset)
		points = append(points, point)
	}
	return Polygon{
		Points: points,
		color:  color,
	}
}

func (p *Polygon) HitBy(ray Ray) (bool, float32) {
	// Stores distances
	collisions := make([]float32, 0, 2)

	// Iterate over edges, taking a point and the next one.
	for currPoint := range p.Points {
		nextPoint := (currPoint + 1) % len(p.Points)

		collide, distance := ray.CollidesWithLine(p.Points[currPoint], p.Points[nextPoint])

		if collide {
			collisions = append(collisions, distance)
		}
	}

	// Select the nearest collision
	if len(collisions) == 0 {
		return false, 0
	}

	min := collisions[0]
	for _, d := range collisions {
		if d < min {
			min = d
		}
	}
	return true, min
}

func (p *Polygon) Color() rl.Color { return p.color }

func (p *Polygon) Center() rl.Vector2 {
	sum := rl.Vector2{X: 0, Y: 0}
	for _, p := range p.Points {
		sum = rl.Vector2Add(sum, p)
	}
	pointNum := float32(len(p.Points))
	return rl.Vector2{X: sum.X / pointNum, Y: sum.Y / pointNum}
}
