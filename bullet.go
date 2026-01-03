package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Bullet struct {
	pos    Vector
	vel    Vector
	angle  float64
	active bool
	life   int
}

func (b *Bullet) draw(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		float32(b.pos.X),
		float32(b.pos.Y),
		2,
		2,
		color.RGBA{R: 255, A: 255},
		false,
	)
}

func (b *Bullet) update() {
	if b.active {
		b.pos = b.pos.Add(b.vel)
		b.pos = wrapPosition(b.pos, ScreenWidth, ScreenHeight)
		b.life--
	}

	if b.active && b.life <= 0 {
		b.active = false
	}
}

func rmBullet(slice []*Bullet, i int) []*Bullet {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
