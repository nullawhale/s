package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"time"
)

var window *sdl.Window
var renderer *sdl.Renderer

type Direction int32

const (
	IDLE  Direction = -1
	LEFT  Direction = 0
	RIGHT Direction = 1
	UP    Direction = 2
	DOWN  Direction = 3
)

type Part struct {
	x, y int32
}

type Snake struct {
	x, y   int32
	dx, dy int32
	size   int32
	body   []Part
}

type Apple struct {
	x, y int32
	size int32
}

const ScreenWidth = 200
const ScreenHeight = 200
const ObjectSize = 10

func run() (err error) {
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return
	}
	defer sdl.Quit()

	if window, err = sdl.CreateWindow("Input",
		150, 400,
		ScreenWidth, ScreenHeight, sdl.WINDOW_SHOWN); err != nil {
		return err
	}
	defer window.Destroy()

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return err
	}

	var snake = Snake{
		ObjectSize, ObjectSize,
		0, 0,
		ObjectSize,
		[]Part{{0, 0}},
	}
	var apple = Apple{
		rand.Int31n(ScreenWidth/ObjectSize) * ObjectSize,
		rand.Int31n(ScreenHeight/ObjectSize) * ObjectSize,
		ObjectSize,
	}
	var d = IDLE
	tick := time.Tick(65 * time.Millisecond)
	done := make(chan bool)
	running := true

	for running {
		select {
		case <-done:
			return
		case <-tick:
			if err = clearScreen(renderer); err != nil {
				return err
			}

			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch e := event.(type) {
				case *sdl.QuitEvent:
					running = false
					fmt.Print(e)
				case *sdl.KeyboardEvent:
					switch e.Keysym.Sym {
					case sdl.K_LEFT:
						if d != RIGHT {
							d = LEFT
						}
					case sdl.K_RIGHT:
						if d != LEFT {
							d = RIGHT
						}
					case sdl.K_UP:
						if d != DOWN {
							d = UP
						}
					case sdl.K_DOWN:
						if d != UP {
							d = DOWN
						}
					}
				}
			}

			snake.update(d)
			if err = snake.draw(renderer); err != nil {
				return err
			}
			if err = apple.draw(renderer); err != nil {
				return err
			}
			if snake.eat(apple) {
				apple = Apple{
					rand.Int31n(ScreenWidth/ObjectSize) * ObjectSize,
					rand.Int31n(ScreenHeight/ObjectSize) * ObjectSize, ObjectSize,
				}
				if err = apple.draw(renderer); err != nil {
					return err
				}
			}

			if snake.dead() {
				running = false
			}
			renderer.Present()
		}
	}

	return
}

func clearScreen(renderer *sdl.Renderer) (err error) {
	if err = renderer.SetDrawColor(0, 0, 0, 255); err != nil {
		return err
	}
	if err = renderer.Clear(); err != nil {
		return err
	}

	return
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

func (a *Apple) draw(renderer *sdl.Renderer) (err error) {
	var outlineRect = sdl.Rect{X: a.x, Y: a.y, W: a.size, H: a.size}
	if err = renderer.SetDrawColor(255, 0, 0, 255); err != nil {
		return err
	}
	if err = renderer.FillRect(&outlineRect); err != nil {
		return err
	}

	return
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
