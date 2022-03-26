package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

type Asteroid struct {
	x, y   float32
	radius float32
}

func (a *Asteroid) new() {

}

func (a *Asteroid) draw(renderer *sdl.Renderer) (err error) {
	var offsetx float32
	var offsety float32
	var d float32

	offsetx = 0.0
	offsety = a.radius
	d = a.radius - 1

	if err = renderer.SetDrawColor(255, 0, 0, 255); err != nil {
		return err
	}

	for offsety >= offsetx {
		err = renderer.DrawPointsF([]sdl.FPoint{
			{a.x + offsetx, a.y + offsety}, {a.x + offsety, a.y + offsetx},
			{a.x - offsetx, a.y + offsety}, {a.x - offsety, a.y + offsetx},
			{a.x + offsetx, a.y - offsety}, {a.x + offsety, a.y - offsetx},
			{a.x - offsetx, a.y - offsety}, {a.x - offsety, a.y - offsetx},
		})

		if err != nil {
			log.Print(err)
		}

		if d >= 2*offsetx {
			d -= 2*offsetx + 1
			offsetx += 1
		} else if d < 2*(a.radius-offsety) {
			d += 2*offsety - 1
			offsety -= 1
		} else {
			d += 2 * (offsety - offsetx - 1)
			offsety -= 1
			offsetx += 1
		}
	}

	//if err = renderer.FillRectF(&outlineRect); err != nil {
	//	return err
	//}

	//var outlineRect = sdl.FRect{X: a.x, Y: a.y, W: a.radius, H: a.radius}
	//if err = renderer.SetDrawColor(255, 0, 0, 255); err != nil {
	//	return err
	//}
	//if err = renderer.FillRectF(&outlineRect); err != nil {
	//	return err
	//}

	return
}
