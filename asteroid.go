package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var asteroidRadii = []float64{48, 36, 28, 20, 12}

type Asteroid struct {
	pos   Vector
	vel   Vector
	size  int
	shape []Vector
}

func (a *Asteroid) radius() float64 {
	if a.size < 0 || a.size >= len(asteroidRadii) {
		return asteroidRadii[len(asteroidRadii)-1]
	}
	return asteroidRadii[a.size]
}

func (a *Asteroid) points() int {
	if a.size < 0 {
		return 1
	}
	if a.size >= len(asteroidRadii) {
		return len(asteroidRadii)
	}
	return a.size + 1
}

func (a *Asteroid) update() {
	a.pos = a.pos.Add(a.vel)
	a.pos = wrapPosition(a.pos, ScreenWidth, ScreenHeight)
}

func (a *Asteroid) draw(screen *ebiten.Image) {
	r := a.radius()
	lineColor := color.White

	if len(a.shape) < 3 {
		return
	}

	for i := range len(a.shape) {
		p1 := a.shape[i]
		p2 := a.shape[(i+1)%len(a.shape)]
		x1 := a.pos.X + p1.X*r
		y1 := a.pos.Y + p1.Y*r
		x2 := a.pos.X + p2.X*r
		y2 := a.pos.Y + p2.Y*r
		vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 1, lineColor, false)
	}
}

func newAsteroid(size int, pos Vector) *Asteroid {
	return &Asteroid{
		pos:   pos,
		vel:   randomAsteroidVel(size),
		size:  size,
		shape: generateAsteroidShape(),
	}
}

func generateAsteroidShape() []Vector {
	const pointsMin = 10
	const pointsMax = 16
	count := pointsMin + rng.Intn(pointsMax-pointsMin+1)
	points := make([]Vector, 0, count)

	for i := 0; i < count; i++ {
		angle := float64(i) * (2 * math.Pi / float64(count))
		jitter := 0.65 + rng.Float64()*0.4
		points = append(points, Vector{
			X: math.Cos(angle) * jitter,
			Y: math.Sin(angle) * jitter,
		})
	}

	return points
}
