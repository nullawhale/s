package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Wall struct {
	start sdl.FPoint
	end   sdl.FPoint
}

func (a *Wall) draw(renderer *sdl.Renderer) (err error) {
	if err = renderer.SetDrawColor(255, 255, 255, 255); err != nil {
		return err
	}

	if err = renderer.DrawLineF(a.start.X, a.start.Y, a.end.X, a.end.Y); err != nil {
		return err
	}

	return
}
