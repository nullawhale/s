package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	pos          Vector
	vel          Vector
	angle        float64
	radius       float64
	invulnTicks  int
	fireCooldown int
}

func (p *Player) draw(screen *ebiten.Image) {
	var lineColor color.Color = color.White
	if p.invulnTicks > 0 && (p.invulnTicks/6)%2 == 0 {
		lineColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
	}

	p1 := rotate(p.pos, Vector{X: p.pos.X, Y: p.pos.Y - p.radius}, p.angle)
	p2 := rotate(p.pos, Vector{X: p.pos.X + p.radius*0.8, Y: p.pos.Y + p.radius}, p.angle)
	p3 := rotate(p.pos, Vector{X: p.pos.X - p.radius*0.8, Y: p.pos.Y + p.radius}, p.angle)

	vector.StrokeLine(screen, float32(p1.X), float32(p1.Y), float32(p2.X), float32(p2.Y), 1, lineColor, false)
	vector.StrokeLine(screen, float32(p2.X), float32(p2.Y), float32(p3.X), float32(p3.Y), 1, lineColor, false)
	vector.StrokeLine(screen, float32(p3.X), float32(p3.Y), float32(p1.X), float32(p1.Y), 1, lineColor, false)
}

func (p *Player) update(input PlayerInput) {
	const maxSpeed = 4.2
	const accel = 0.12
	const drag = 0.99
	const brake = 0.96

	if p.invulnTicks > 0 {
		p.invulnTicks--
	}
	if p.fireCooldown > 0 {
		p.fireCooldown--
	}

	if input.Left {
		p.angle -= RotationSpeed
	}
	if input.Right {
		p.angle += RotationSpeed
	}

	fwd := Vector{X: math.Sin(p.angle), Y: -math.Cos(p.angle)}
	if input.Thrust {
		p.vel = p.vel.Add(fwd.Scale(accel))
	}
	if input.Brake {
		p.vel = p.vel.Add(fwd.Scale(-accel * 0.6))
	}

	if input.Brake {
		p.vel = p.vel.Scale(brake)
	} else {
		p.vel = p.vel.Scale(drag)
	}

	if p.vel.Len() > maxSpeed {
		p.vel = p.vel.Normalize().Scale(maxSpeed)
	}

	p.pos = p.pos.Add(p.vel)
	p.pos = wrapPosition(p.pos, ScreenWidth, ScreenHeight)
}

func rotate(orig Vector, p Vector, a float64) Vector {
	sin := math.Sin(a)
	cos := math.Cos(a)

	newX := cos*(p.X-orig.X) - sin*(p.Y-orig.Y) + orig.X
	newY := sin*(p.X-orig.X) + cos*(p.Y-orig.Y) + orig.Y

	return Vector{X: newX, Y: newY}
}
