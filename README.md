# Asteroids (Ebiten)

A minimal Asteroids clone written in Go using Ebiten.

## Requirements

- Go 1.25+

## Run

```bash
go run .
```

## Controls

- Left/Right: rotate
- Up: thrust
- Down: brake
- Space: shoot
- Enter/Space: start or restart
- Esc: pause/resume

## Notes

- Asteroids have 5 sizes; each hit shrinks them until the smallest disappears.
- Score increases with smaller asteroids (1..5 points).
