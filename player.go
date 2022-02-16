package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Player struct {
	center sdl.FPoint
	dx, dy float32
	size   float32
	angle  float64
	a      float32
}

func (s *Player) draw(renderer *sdl.Renderer) (err error) {

	if err = renderer.SetDrawColor(255, 255, 255, 255); err != nil {
		return err
	}

	renderer.DrawLinesF([]sdl.FPoint{
		rotate(s.center, sdl.FPoint{X: s.center.X - s.size, Y: s.center.Y - s.size}, s.angle),
		rotate(s.center, sdl.FPoint{X: s.center.X + s.size, Y: s.center.Y - s.size}, s.angle),
		rotate(s.center, sdl.FPoint{X: s.center.X, Y: s.center.Y + s.size*2}, s.angle),
		rotate(s.center, sdl.FPoint{X: s.center.X - s.size, Y: s.center.Y - s.size}, s.angle),
	})

	//if err = renderer.SetDrawColor(255, 0, 0, 255); err != nil {
	//	return err
	//}
	//
	//renderer.DrawPointF(s.center.X, s.center.Y)

	return
}

func rotate(orig sdl.FPoint, p sdl.FPoint, a float64) sdl.FPoint {
	sin := float32(math.Sin(a))
	cos := float32(math.Cos(a))

	newX := cos*(p.X-orig.X) - sin*(p.Y-orig.Y) + orig.X
	newY := sin*(p.X-orig.X) + cos*(p.Y-orig.Y) + orig.Y

	return sdl.FPoint{X: newX, Y: newY}
}

func (s *Player) eat(a Apple) bool {
	//if s.center.X == a.X && s.center.Y == a.Y {
	//	s.body = append(s.body, Part{a.X+1, a.Y+1})
	//	return true
	//}
	return false
}

func (s *Player) dead() bool {
	//if s.body.X == s.body.X && s.body.Y == s.body.Y {
	//	return true
	//}
	return false
}

func (s *Player) update(d Direction) {
	s.center.X += s.dx
	s.center.Y += s.dy

	if s.center.X >= ScreenWidth {
		s.center.X = 0
	}
	if s.center.X <= -ObjectSize {
		s.center.X = ScreenWidth
	}
	if s.center.Y >= ScreenHeight {
		s.center.Y = 0
	}
	if s.center.Y <= -ObjectSize {
		s.center.Y = ScreenHeight
	}

	switch d {
	case LEFT:
		s.angle -= RotationSpeed
		break
	case RIGHT:
		s.angle += RotationSpeed
		break
	case UP:
		if s.a < 1 {
			s.a += 0.05
		}
		s.dx = -s.a * float32(math.Sin(s.angle))
		s.dy = s.a * float32(math.Cos(s.angle))
	case DOWN:
		if s.a > 0 {
			s.a -= 0.05
		}
		s.dx = -s.a * float32(math.Sin(s.angle))
		s.dy = s.a * float32(math.Cos(s.angle))
	case IDLE:
		s.a = 0
		s.dx = 0
		s.dy = 0
	}

	if s.angle > math.Pi {
		s.angle -= 2 * math.Pi
	}
	if s.angle < -math.Pi {
		s.angle += 2 * math.Pi
	}
}
