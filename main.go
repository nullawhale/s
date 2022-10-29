package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math"
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

var window *sdl.Window
var renderer *sdl.Renderer
var player Player

var bullets []*Bullet
var asteroids []*Asteroid
var direction Direction
var WorldMap Map
var walls []*Wall

func run() <-chan error {
	ticker := time.NewTicker(time.Second / 60)
	errors := make(chan error)
	running := true

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
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
					direction = IDLE
				case sdl.K_f:
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

		for _, asteroid := range asteroids {
			if err := asteroid.draw(renderer); err != nil {
				errors <- err
			}
		}

		for _, wall := range walls {
			if err := wall.draw(renderer); err != nil {
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

	WorldMap.LoadMap()

	width := WorldMap.World.ScreenWidth
	height := WorldMap.World.ScreenHeight

	window, err = sdl.CreateWindow(
		"Input", 100, 500,
		width, height, sdl.WINDOW_SHOWN,
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
		center: sdl.FPoint{
			X: WorldMap.Player.Coordinates.X,
			Y: WorldMap.Player.Coordinates.Y,
		},
		angle: WorldMap.Player.Angle * math.Pi / 180,
		size:  float32(WorldMap.Const.ObjectSize),
	}

	player.rays = append(player.rays, &Ray{
		player.center,
		sdl.FPoint{X: player.center.X, Y: player.center.Y + float32(100)},
		player.angle,
	})

	if len(WorldMap.Asteroids) > 0 {
		for _, asteroid := range WorldMap.Asteroids {
			asteroids = append(asteroids, &Asteroid{
				asteroid.Coordinates.X,
				asteroid.Coordinates.Y,
				asteroid.Radius,
			})
		}
	}

	if len(WorldMap.Walls) > 0 {
		for _, wall := range WorldMap.Walls {
			walls = append(walls, &Wall{
				start: sdl.FPoint{X: wall.Start.X, Y: wall.Start.Y},
				end:   sdl.FPoint{X: wall.End.X, Y: wall.End.Y},
			})
		}
	}

	direction = IDLE

	if err = runWorld(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
