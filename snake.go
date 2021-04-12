package main

import "github.com/veandco/go-sdl2/sdl"

type Part struct {
	x, y int32
}

type Snake struct {
	x, y   int32
	dx, dy int32
	size   int32
	body   []Part
}

func (s *Snake) draw(renderer *sdl.Renderer) (err error) {
	for i := 0; i < len(s.body); i++ {
		rect := sdl.Rect{X: s.body[i].x, Y: s.body[i].y, W: s.size, H: s.size}
		if err = renderer.SetDrawColor(255, 255, 255, 255); err != nil {
			return err
		}
		if err = renderer.FillRect(&rect); err != nil {
			return err
		}
	}
	return
}

func (s *Snake) eat(a Apple) bool {
	if s.x == a.x && s.y == a.y {
		s.body = append(s.body, Part{a.x+1, a.y+1})
		return true
	}
	return false
}

func (s *Snake) dead() bool {
	for i := 1; i < len(s.body); i++ {
		if s.body[i].x == s.body[0].x && s.body[i].y == s.body[0].y {
			return true
		}
	}
	return false
}

func (s *Snake) update(d Direction) {
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
		s.dx = -ObjectSize
		s.dy = 0
	case RIGHT:
		s.dx = ObjectSize
		s.dy = 0
	case UP:
		s.dy = -ObjectSize
		s.dx = 0
	case DOWN:
		s.dy = ObjectSize
		s.dx = 0
	case IDLE:
		s.dx = 0
		s.dy = 0
	}

	for i := len(s.body) - 1; i > 0; i-- {
		s.body[i] = s.body[i-1]
	}
	s.body[0].x, s.body[0].y = s.x, s.y
}
