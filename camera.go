package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ray struct {
	Origin    rl.Vector2
	Direction rl.Vector2 // The ray's direction, normalized
}

func NewRay(origin, direction rl.Vector2) Ray {
	return Ray{
		Origin:    origin,
		Direction: rl.Vector2Normalize(direction),
	}
}

// Returns if there is a collision, and the distance to it.
func (r *Ray) CollidesWithLine(p1, p2 rl.Vector2) (bool, float32) {
	// Direction vectors
	rO := r.Origin
	rDir := rl.Vector2Normalize(r.Direction)
	lDir := rl.Vector2Subtract(p2, p1)

	/* We find line intersections by Cramer's rule:
	 *  ┌               ┐ ┌ ┐   ┌         ┐
	 *  │rDir.X  -lDir.X│ │t│   │p1.X-rO.x│
	 *  │rDir.Y  -lDir.Y│ │u│ = │p1.Y-rO.Y│
	 *  └               ┘ └ ┘   └         ┘
	 */

	det := -rDir.X*lDir.Y + lDir.X*rDir.Y

	// If the determinant is zero, the equations are linearly dependant, meaning they are parallel.
	parallel := math.Abs(float64(det)) <= 1e-6
	if !parallel {

		// Solve the system of equations
		t := (-(p1.X-rO.X)*lDir.Y + lDir.X*(p1.Y-rO.Y)) / det
		u := (rDir.X*(p1.Y-rO.Y) - (p1.X-rO.X)*rDir.Y) / det

		// The collision is valid only if t is positive and u is between 0 and 1
		if t >= 0 && u >= 0 && u <= 1 {
			// The distance is just t (the ray's parameter), since direction is normalized
			return true, t
		}
	}

	return false, -1
}

// Currently, camera only looks to the +y direction
type Camera struct {
	Position                       rl.Vector2
	Direction                      rl.Vector2 // Direction of the camera, normalized.
	PixelCount                     int32
	ViewPortDistance, ViewPortSize float32
	World                          *World
}

func NewCamera(position, direction rl.Vector2, world *World) *Camera {
	return &Camera{
		Position:         position,
		Direction:        rl.Vector2Normalize(direction), // Ensure its normalized
		PixelCount:       WINDOW_WIDTH,
		ViewPortDistance: 1,
		ViewPortSize:     2, // Gives a 90 degree FOV (FOV = tg-1(viewportsize/(2*distance)))
		World:            world,
	}
}

func (c *Camera) CastRays() []rl.Color {
	colors := make([]rl.Color, c.PixelCount)
	// TODO may be a good place to use goroutines
	for pixel := int32(0); pixel < c.PixelCount; pixel++ {
		// Calculate the corresponding coordinates of the viewframe and create a ray
		// using that as a direction.
		// Substract the camera position so its a direction vector.
		dir := rl.Vector2Subtract(c.getViewPortPixel(pixel), c.Position)
		ray := Ray{
			Origin:    c.Position,
			Direction: dir,
		}

		obj, dist := checkCollisions(c.World, ray)

		if obj != nil {
			colors[pixel] = calculateColor(obj, dist)
		} else {
			colors[pixel] = BACKGROUND_COLOR
		}
	}

	return colors
}

// Returns the point contained in the viewport that corresponds to the pixel in the screen.
func (c *Camera) getViewPortPixel(pixel int32) rl.Vector2 {
	/* We can find a point in the viewport by offsetting from the center by a perpendicular vector.
	 *
	 *  Viewport
	 *  ├───────────┤
	 *        ^ ─> basis (length 1)
	 *        │
	 *  	  │ viewportCenter
	 *  	  │
	 *  	  O
	 *  	Camera
	 */

	// Perpendicular vector to the direction to the viewport acts as a basis for the viewport.
	basis := rl.Vector2Normalize(rl.Vector2{X: -c.Direction.Y, Y: c.Direction.X})

	// The offset from the center of the viewport to the pixel.
	centeredPixel := pixel - c.PixelCount/2
	normalizedOffset := float32(centeredPixel) / float32(c.PixelCount)
	offset := rl.Vector2Scale(
		basis,
		normalizedOffset*c.ViewPortSize,
	)

	viewportCenter := rl.Vector2Add(c.Position, rl.Vector2Scale(c.Direction, c.ViewPortDistance))

	// The point we are looking for is obtained by adding the offset to the center
	return rl.Vector2Add(viewportCenter, offset)
}

func checkCollisions(w *World, r Ray) (WorldObject, float32) {

	collisions := make(map[WorldObject]float32)

	// Check collisions for every object.
	for _, o := range *w {
		hit, dist := o.HitBy(r)
		if hit {
			collisions[o] = dist
		}
	}

	// Use the closest collision
	minDist := float32(VIEW_DISTANCE)
	var closestObj WorldObject
	for o, dist := range collisions {
		if dist < minDist {
			minDist = dist
			closestObj = o
		}
	}

	// closestObj may be nil if there are no collisions.
	return closestObj, minDist
}

// Darken the color by a function of the distance
func calculateColor(obj WorldObject, dist float32) rl.Color {
	return rl.ColorBrightness(
		obj.Color(),
		-dist/VIEW_DISTANCE,
	)
}
