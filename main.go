package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"runtime"
	"time"
	//"time"
)

type Direction int32

const (
	IDLE  Direction = -1
	LEFT  Direction = 0
	RIGHT Direction = 1
	UP    Direction = 2
	DOWN  Direction = 3
)

const ScreenWidth = 300
const ScreenHeight = 300
const ObjectSize = 10
const ObjectSpeed = 3
const BulletSpeed = 2
const RotationSpeed = 0.05

var window *sdl.Window
var renderer *sdl.Renderer
var player Player

var bullets []*Bullet
var apple Apple
var direction Direction

func run() <-chan error {
	ticker := time.NewTicker(time.Second / 60)
	errors := make(chan error)
	running := true

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				fmt.Print("123")
				running = false
			case *sdl.KeyboardEvent:
				switch e.Keysym.Sym {
				case sdl.K_LEFT:
					direction = LEFT
				case sdl.K_RIGHT:
					direction = RIGHT
				case sdl.K_UP:
					direction = UP
				case sdl.K_DOWN:
					direction = DOWN
				case sdl.K_SPACE:
					bullets = append(bullets, &Bullet{
						pos:    player.center,
						angle:  player.angle,
						active: true,
					})
				}
			}
		}

		<-ticker.C

		if err := clearScreen(renderer); err != nil {
			errors <- err
		}
		player.update(direction)
		if err := player.draw(renderer); err != nil {
			errors <- err
		}

		for i, bullet := range bullets {
			bullet.update()
			if err := bullet.draw(renderer); err != nil {
				errors <- err
			}
			if !bullet.active {
				bullets = rmBullet(bullets, i)
			}
		}

		if err := apple.draw(renderer); err != nil {
			errors <- err
		}

		if player.eat(apple) {
			apple = Apple{
				float32(rand.Int31n(ScreenWidth/ObjectSize) * ObjectSize),
				float32(rand.Int31n(ScreenHeight/ObjectSize) * ObjectSize), ObjectSize,
			}
			if err := apple.draw(renderer); err != nil {
				errors <- err
			}
		}

		if player.dead() {
			running = false
		}
		renderer.Present()
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
	//events := make(chan sdl.Event)
	errorC := run()

	runtime.LockOSThread()
	if errorC != nil {
		return errors
	}

	return nil
}

func main() {
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer sdl.Quit()

	window, err = sdl.CreateWindow(
		"Input", 100, 500,
		ScreenWidth, ScreenHeight, sdl.WINDOW_SHOWN,
	)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	player = Player{
		center: sdl.FPoint{X: ObjectSize * 3, Y: ObjectSize * 3},
		size:   ObjectSize,
	}

	//apple = Apple{
	//	float32(rand.Int31n(ScreenWidth/ObjectSize) * ObjectSize),
	//	float32(rand.Int31n(ScreenHeight/ObjectSize) * ObjectSize),
	//	ObjectSize,
	//}
	direction = IDLE

	if err = runWorld(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
