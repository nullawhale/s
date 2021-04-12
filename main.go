package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"runtime"
	"time"
)

type Direction int32

const (
	IDLE  Direction = -1
	LEFT  Direction = 0
	RIGHT Direction = 1
	UP    Direction = 2
	DOWN  Direction = 3
)

const ScreenWidth = 200
const ScreenHeight = 200
const ObjectSize = 10

var window *sdl.Window
var renderer *sdl.Renderer
var snake Snake
var apple Apple
var direction Direction

func run(events <-chan sdl.Event) <-chan error {
	tick := time.Tick(65 * time.Millisecond)
	//ticks := sdl.GetTicks()
	//done := make(chan bool)
	errors := make(chan error)
	running := true

	for running {
		select {
		case event := <-events:
			switch e := event.(type) {
			case *sdl.QuitEvent:
				fmt.Print("123")
				running = false
			case *sdl.KeyboardEvent:
				switch e.Keysym.Sym {
				case sdl.K_LEFT:
					if direction != RIGHT {
						direction = LEFT
					}
				case sdl.K_RIGHT:
					if direction != LEFT {
						direction = RIGHT
					}
				case sdl.K_UP:
					if direction != DOWN {
						direction = UP
					}
				case sdl.K_DOWN:
					if direction != UP {
						direction = DOWN
					}
				}
			}
		case <-tick:
			if err := clearScreen(renderer); err != nil {
				errors <- err
			}

			/*for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch e := event.(type) {
				case *sdl.QuitEvent:
					fmt.Print("123")
					running = false
				case *sdl.KeyboardEvent:
					switch e.Keysym.Sym {
					case sdl.K_LEFT:
						if direction != RIGHT {
							direction = LEFT
						}
					case sdl.K_RIGHT:
						if direction != LEFT {
							direction = RIGHT
						}
					case sdl.K_UP:
						if direction != DOWN {
							direction = UP
						}
					case sdl.K_DOWN:
						if direction != UP {
							direction = DOWN
						}
					}
				}
			}*/

			snake.update(direction)
			if err := snake.draw(renderer); err != nil {
				errors <- err
			}
			if err := apple.draw(renderer); err != nil {
				errors <- err
			}
			if snake.eat(apple) {
				apple = Apple{
					rand.Int31n(ScreenWidth/ObjectSize) * ObjectSize,
					rand.Int31n(ScreenHeight/ObjectSize) * ObjectSize, ObjectSize,
				}
				if err := apple.draw(renderer); err != nil {
					errors <- err
				}
			}

			if snake.dead() {
				running = false
			}
			renderer.Present()
		}
	}

	return errors
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

func runWorld() (errors error) {
	events := make(chan sdl.Event)
	errorC := run(events)

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errorC:
			return err
		}
	}
}

func main() {
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer sdl.Quit()

	window, err = sdl.CreateWindow(
		"Input", 150, 400,
		ScreenWidth, ScreenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	snake = Snake{
		ObjectSize, ObjectSize,
		0, 0,
		ObjectSize,
		[]Part{{0, 0}},
	}
	apple = Apple{
		rand.Int31n(ScreenWidth/ObjectSize) * ObjectSize,
		rand.Int31n(ScreenHeight/ObjectSize) * ObjectSize,
		ObjectSize,
	}
	direction = IDLE

	if err = runWorld(); err != nil {
		_, _ = fmt.Printf("Error: %s\n", err)
	}
}
