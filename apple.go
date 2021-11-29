package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Apple struct {
	x, y float32
	size float32
}

func (a *Apple) new() {

}

func (a *Apple) draw(renderer *sdl.Renderer) (err error) {
	var outlineRect = sdl.FRect{X: a.x, Y: a.y, W: a.size, H: a.size}
	if err = renderer.SetDrawColor(255, 0, 0, 255); err != nil {
		return err
	}
	if err = renderer.FillRectF(&outlineRect); err != nil {
		return err
	}

	return
}
