package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var asteroidRadii = []float64{48, 36, 28, 20, 12}

type Asteroid struct {
	pos  Vector
	vel  Vector
	size int
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
	const segments = 10
	r := a.radius()
	lineColor := color.White

	for i := range segments {
		angle1 := float64(i) * (2 * math.Pi / segments)
		angle2 := float64(i+1) * (2 * math.Pi / segments)
		x1 := a.pos.X + math.Cos(angle1)*r
		y1 := a.pos.Y + math.Sin(angle1)*r
		x2 := a.pos.X + math.Cos(angle2)*r
		y2 := a.pos.Y + math.Sin(angle2)*r
		vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 1, lineColor, false)
	}
}
