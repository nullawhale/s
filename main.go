package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

const (
	ScreenWidth       = 800
	ScreenHeight      = 600
	RotationSpeed     = 0.06
	BulletSpeed       = 6.5
	BulletLife        = 60
	FireCooldownTicks = 10
	InitialLives      = 3
	PlayerRadius      = 12
)

type Game struct {
	state     GameState
	player    Player
	bullets   []*Bullet
	asteroids []*Asteroid
	score     int
	lives     int
	wave      int
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func (g *Game) Update() error {
	switch g.state {
	case StateMenu:
		if isStartPressed() {
			g.startGame()
		}
		return nil
	case StatePaused:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StatePlaying
		}
		return nil
	case StateGameOver:
		if isStartPressed() {
			g.startGame()
		}
		return nil
	case StatePlaying:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StatePaused
			return nil
		}
	}

	input := PlayerInput{
		Left:   ebiten.IsKeyPressed(ebiten.KeyLeft),
		Right:  ebiten.IsKeyPressed(ebiten.KeyRight),
		Thrust: ebiten.IsKeyPressed(ebiten.KeyUp),
		Brake:  ebiten.IsKeyPressed(ebiten.KeyDown),
	}
	g.player.update(input)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.player.fireCooldown == 0 {
		g.spawnBullet()
		g.player.fireCooldown = FireCooldownTicks
	}

	for i := 0; i < len(g.bullets); i++ {
		bullet := g.bullets[i]
		bullet.update()
		if !bullet.active {
			g.bullets = rmBullet(g.bullets, i)
			i--
		}
	}

	for _, asteroid := range g.asteroids {
		asteroid.update()
	}

	g.handleBulletAsteroidCollisions()
	g.handlePlayerAsteroidCollisions()

	if len(g.asteroids) == 0 {
		g.wave++
		g.asteroids = g.spawnAsteroids(g.player.pos)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	switch g.state {
	case StateMenu:
		drawCenteredText(screen, "ASTEROIDS", 180)
		drawCenteredText(screen, "Press ENTER or SPACE to start", 240)
		drawCenteredText(screen, "Arrows to move, SPACE to shoot", 270)
		return
	case StatePaused:
		g.drawWorld(screen)
		drawCenteredText(screen, "PAUSED", 260)
		return
	case StateGameOver:
		g.drawWorld(screen)
		drawCenteredText(screen, "GAME OVER", 230)
		drawCenteredText(screen, fmt.Sprintf("Score: %d", g.score), 255)
		drawCenteredText(screen, "Press ENTER or SPACE to restart", 290)
		return
	case StatePlaying:
		g.drawWorld(screen)
	}
}

func (g *Game) drawWorld(screen *ebiten.Image) {
	g.player.draw(screen)
	for _, bullet := range g.bullets {
		bullet.draw(screen)
	}
	for _, asteroid := range g.asteroids {
		asteroid.draw(screen)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Wave: %d", g.wave), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score), 10, 24)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Lives: %d", g.lives), 10, 38)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) startGame() {
	g.score = 0
	g.lives = InitialLives
	g.wave = 1
	g.bullets = nil
	g.player = newPlayer()
	g.asteroids = g.spawnAsteroids(g.player.pos)
	g.state = StatePlaying
}

func (g *Game) spawnBullet() {
	forward := Vector{X: math.Sin(g.player.angle), Y: -math.Cos(g.player.angle)}
	start := g.player.pos.Add(forward.Scale(g.player.radius + 4))
	vel := forward.Scale(BulletSpeed).Add(g.player.vel)
	g.bullets = append(g.bullets, &Bullet{
		pos:    start,
		vel:    vel,
		angle:  g.player.angle,
		active: true,
		life:   BulletLife,
	})
}

func (g *Game) handleBulletAsteroidCollisions() {
	for bi := 0; bi < len(g.bullets); bi++ {
		bullet := g.bullets[bi]
		if !bullet.active {
			continue
		}
		for ai := 0; ai < len(g.asteroids); ai++ {
			asteroid := g.asteroids[ai]
			if bullet.pos.Sub(asteroid.pos).Len() <= asteroid.radius()+2 {
				bullet.active = false
				g.score += asteroid.points()
				if asteroid.size+1 < len(asteroidRadii) {
					asteroid.size++
					asteroid.vel = asteroid.vel.Normalize().Scale(0.8 + float64(asteroid.size)*0.45)
					if asteroid.vel.Len() == 0 {
						asteroid.vel = randomAsteroidVel(asteroid.size)
					}
				} else {
					g.asteroids = append(g.asteroids[:ai], g.asteroids[ai+1:]...)
					ai--
				}
				break
			}
		}
	}
}

func (g *Game) handlePlayerAsteroidCollisions() {
	if g.player.invulnTicks > 0 {
		return
	}
	for _, asteroid := range g.asteroids {
		if g.player.pos.Sub(asteroid.pos).Len() <= g.player.radius+asteroid.radius() {
			g.lives--
			if g.lives <= 0 {
				g.state = StateGameOver
				return
			}
			g.player = newPlayer()
			g.player.invulnTicks = 120
			return
		}
	}
}

func newPlayer() Player {
	return Player{
		pos:    Vector{X: ScreenWidth / 2, Y: ScreenHeight / 2},
		radius: PlayerRadius,
	}
}

func (g *Game) spawnAsteroids(avoid Vector) []*Asteroid {
	count := 3 + g.wave
	asteroids := make([]*Asteroid, 0, count)

	for range count {
		pos := Vector{X: rng.Float64() * ScreenWidth, Y: rng.Float64() * ScreenHeight}
		asteroids = append(asteroids, &Asteroid{
			pos:  pos,
			vel:  randomAsteroidVel(0),
			size: 0,
		})
	}

	return asteroids
}

func randomAsteroidVel(size int) Vector {
	angle := rng.Float64() * 2 * math.Pi
	speed := 0.6 + float64(size)*0.45
	return Vector{X: math.Cos(angle) * speed, Y: math.Sin(angle) * speed}
}

func isStartPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

func drawCenteredText(screen *ebiten.Image, text string, y int) {
	charWidth := 7
	x := max((ScreenWidth-len(text)*charWidth)/2, 10)
	ebitenutil.DebugPrintAt(screen, text, x, y)
}

func main() {
	game := &Game{state: StateMenu}
	game.player = newPlayer()
	game.wave = 1
	game.asteroids = game.spawnAsteroids(game.player.pos)
	game.lives = InitialLives

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Asteroids")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
