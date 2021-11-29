package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Bullet struct {
	pos    sdl.FPoint
	dx, dy float32
	angle  float64
	active bool
}

func (b *Bullet) draw(renderer *sdl.Renderer) (err error) {
	if err = renderer.SetDrawColor(255, 0, 0, 255); err != nil {
		return err
	}

	renderer.DrawPointF(b.pos.X, b.pos.Y)

	return
}

func (b *Bullet) update() {
	if b.active {
		b.pos.X += b.dx
		b.pos.Y += b.dy

		b.dx = -BulletSpeed * float32(math.Sin(b.angle))
		b.dy = BulletSpeed * float32(math.Cos(b.angle))
	}

	if b.active && (b.pos.X >= ScreenWidth || b.pos.X <= 0 || b.pos.Y >= ScreenHeight || b.pos.Y <= 0) {
		b.active = false
	}
}

func rmBullet(slice []*Bullet, i int) []*Bullet {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
