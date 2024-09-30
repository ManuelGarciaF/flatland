package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ray struct {
	Origin    rl.Vector2
	Direction rl.Vector2
}

func (r *Ray) End() rl.Vector2 {
	return rl.Vector2Add(r.Origin,
		rl.Vector2Scale(rl.Vector2Subtract(r.Direction, r.Origin), float32(VIEW_DISTANCE)))
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
	// TODO paralellize this
	for pixel := int32(0); pixel < c.PixelCount; pixel++ {
		// Calculate the corresponding coordinates of the viewframe and create a ray
		// using that as a direction.
		ray := Ray{
			Origin:    c.Position,
			Direction: c.getViewPortPixel(pixel),
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
	 *        ▲─► basis (length 1)
	 *        │
	 *  	  │ viewportCenter
	 *  	  │
	 *  	  O
	 *  	Camera
	 */

	// Perpendicular vector to the direction to the viewport acts as a basis for the viewport.
	basis := rl.Vector2Normalize(rl.Vector2{X: c.Direction.Y, Y: -c.Direction.X})

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
		-float32(math.Sinh(float64(dist/VIEW_DISTANCE))),
	)
}
