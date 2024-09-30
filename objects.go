package main

import rl "github.com/gen2brain/raylib-go/raylib"

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

func (p *Polygon) HitBy(ray Ray) (bool, float32) {
	// Stores distances
	collisions := make([]float32, 0, 2)

	// Iterate over edges, taking a point and the next one.
	for currPoint := range p.Points {
		nextPoint := (currPoint + 1) % len(p.Points)

		// TODO implement this manually
		var collisionPoint rl.Vector2
		collide := rl.CheckCollisionLines(
			ray.Origin,
			ray.End(),
			p.Points[currPoint],
			p.Points[nextPoint],
			&collisionPoint,
		)

		if collide {
			distance := rl.Vector2Distance(ray.Origin, collisionPoint)
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
	return rl.Vector2{X: sum.X/pointNum, Y: sum.Y/pointNum}
}
