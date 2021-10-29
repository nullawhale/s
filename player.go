package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Part struct {
	x, y int32
}

type Player struct {
	x, y   int32
	dx, dy int32
	size   int32
	angle  int32
}

func (s *Player) draw(renderer *sdl.Renderer) (err error) {
	//for i := 0; i < len(s.body); i++ {}
	if err = renderer.SetDrawColor(255, 0, 0, 255); err != nil {
		return err
	}

	renderer.DrawLines([]sdl.Point{
		{s.x + 1, s.y},
		{s.x - 1, s.y - 1},
		{s.x + 1, s.y + 1},
		{s.x + 1, s.y},
	})

	if err = renderer.SetDrawColor(255, 255, 255, 255); err != nil {
		return err
	}

	renderer.DrawLines([]sdl.Point{
		{s.x, s.y - s.size},
		{s.x - s.size, s.y + s.size},
		{s.x + s.size, s.y + s.size},
		{s.x, s.y - s.size},
	})

	return
}

func (s *Player) eat(a Apple) bool {
	if s.x == a.x && s.y == a.y {
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
	s.x += s.dx
	s.y += s.dy

	if s.x >= ScreenWidth {
		s.x = 0
	}
	if s.x <= -ObjectSize {
		s.x = ScreenWidth
	}
	if s.y >= ScreenHeight {
		s.y = 0
	}
	if s.y <= -ObjectSize {
		s.y = ScreenHeight
	}

	switch d {
	case LEFT:
		s.angle = -RotationSpeed
	case RIGHT:
		s.angle = RotationSpeed
	case UP:
		//angleTmp := s.angle

		s.dx = -int32(math.Sin(float64(s.angle) * math.Pi / 180))
		s.dy = int32(math.Cos(float64(s.angle) * math.Pi / 180))
	case DOWN:
		s.dy = ObjectSpeed
		s.dx = 0
	case IDLE:
		s.dx = 0
		s.dy = 0
	}
}
