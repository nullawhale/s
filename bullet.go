package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Bullet struct {
	pos    sdl.FPoint
	dx, dy float32
	angle  float32
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
	}

	if b.active && (b.pos.X >= ScreenWidth || b.pos.X <= 0 || b.pos.Y >= ScreenHeight || b.pos.Y <= 0) {
		b.active = false
	}
}
