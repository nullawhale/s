package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Part struct {
	x, y float32
}

type Player struct {
	p      Part
	dx, dy float32
	size   float32
	angle  float32
}

func (s *Player) draw(renderer *sdl.Renderer) (err error) {
	//for i := 0; i < len(s.body); i++ {}

	if err = renderer.SetDrawColor(255, 255, 255, 255); err != nil {
		return err
	}

	x, y := rotate(s.p, Part{s.p.x, s.p.y - s.size}, s.angle)
	//x2, y2 := rotate(s.p, Part{s.p.x - s.size, s.p.y}, s.angle)
	//x2, y2 := rotate(s.p, Part{s.p.x - s.size, s.p.y + s.size}, s.angle)
	//x3, y3 := rotate(s.p, Part{s.p.x + s.size, s.p.y + s.size}, s.angle)

	renderer.DrawLinesF([]sdl.FPoint{
		{x, y},
		{s.p.x, s.p.y},
		//{x2, y2},
		//{x3, y3},
		//{x, y},
	})

	if err = renderer.SetDrawColor(255, 0, 0, 255); err != nil {
		return err
	}

	renderer.DrawLinesF([]sdl.FPoint{
		{s.p.x, s.p.y},
		{s.p.x, s.p.y},
	})

	return
}

func rotate(orig Part, p Part, a float32) (float32, float32) {
	sin := math.Sin(rad(float64(a)))
	cos := math.Cos(rad(float64(a)))

	newX := orig.x + float32(cos)*(p.x-orig.x) - float32(sin)*(p.y-orig.y)
	newY := orig.y + float32(sin)*(p.x-orig.x) - float32(cos)*(p.y-orig.y)

	return newX, newY
}

func (s *Player) eat(a Apple) bool {
	if s.p.x == a.x && s.p.y == a.y {
		//s.body = append(s.body, Part{a.x+1, a.y+1})
		return true
	}
	return false
}

func (s *Player) dead() bool {
	//if s.body.x == s.body.x && s.body.y == s.body.y {
	//	return true
	//}
	return false
}

func (s *Player) update(d Direction) {
	s.p.x += s.dx
	s.p.y += s.dy

	if s.p.x >= ScreenWidth {
		s.p.x = 0
	}
	if s.p.x <= -ObjectSize {
		s.p.x = ScreenWidth
	}
	if s.p.y >= ScreenHeight {
		s.p.y = 0
	}
	if s.p.y <= -ObjectSize {
		s.p.y = ScreenHeight
	}

	switch d {
	case LEFT:
		s.angle += RotationSpeed
		//fmt.Printf("%f ", s.angle)
		break
	case RIGHT:
		s.angle -= RotationSpeed
		//fmt.Printf("%f ", s.angle)
		break
	case UP:
		//angleTmp := s.angle

		s.dx = float32(math.Sin(rad(float64(s.angle))))
		s.dy = float32(math.Cos(rad(float64(s.angle))))

		fmt.Printf("%f ", s.dx)
	case DOWN:
		s.dy = 0
		s.dx = 0
	case IDLE:
		s.dx = 0
		s.dy = 0
	}

	if s.angle >= 360 || s.angle <= -360 {
		s.angle = 0
	}
}

func rad(degree float64) float64 {
	return degree * math.Pi / 180
}
