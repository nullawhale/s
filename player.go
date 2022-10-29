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
	rays   []*Ray
}

type Ray struct {
	start sdl.FPoint
	end   sdl.FPoint
	dir   float64
}

func (p *Player) draw(renderer *sdl.Renderer) (err error) {

	if err = renderer.SetDrawColor(255, 255, 255, 255); err != nil {
		return err
	}

	err = renderer.DrawLinesF([]sdl.FPoint{
		rotate(p.center, sdl.FPoint{X: p.center.X - p.size, Y: p.center.Y - p.size}, p.angle),
		rotate(p.center, sdl.FPoint{X: p.center.X + p.size, Y: p.center.Y - p.size}, p.angle),
		rotate(p.center, sdl.FPoint{X: p.center.X, Y: p.center.Y + p.size*2}, p.angle),
		rotate(p.center, sdl.FPoint{X: p.center.X - p.size, Y: p.center.Y - p.size}, p.angle),
	})
	if err != nil {
		return err
	}

	/*renderer.DrawLinesF([]sdl.FPoint{
		rotate(p.center, sdl.FPoint{X: p.center.X, Y: p.center.Y + p.size*2}, p.angle),
		rotate(p.center, sdl.FPoint{X: p.center.X, Y: p.center.Y + 200}, p.angle),
	})*/

	//renderer.DrawLineF(rayStart.X, rayStart.Y, rayEnd.X, rayEnd.Y)

	for _, ray := range p.rays {
		rayStart := rotate(p.center, sdl.FPoint{X: ray.start.X, Y: ray.start.Y + p.size*2}, p.angle)
		rayEnd := rotate(p.center, sdl.FPoint{X: ray.end.X, Y: ray.end.Y}, p.angle)
		if err = renderer.DrawLineF(rayStart.X, rayStart.Y, rayEnd.X, rayEnd.Y); err != nil {
			return err
		}
	}

	//if err = renderer.SetDrawColor(255, 0, 0, 255); err != nil {
	//	return err
	//}
	//
	//renderer.DrawPointF(p.center.X, p.center.Y)

	return
}

func rotate(orig sdl.FPoint, p sdl.FPoint, a float64) sdl.FPoint {
	sin := float32(math.Sin(a))
	cos := float32(math.Cos(a))

	newX := cos*(p.X-orig.X) - sin*(p.Y-orig.Y) + orig.X
	newY := sin*(p.X-orig.X) + cos*(p.Y-orig.Y) + orig.Y

	return sdl.FPoint{X: newX, Y: newY}
}

func (p *Player) eat(a Asteroid) bool {
	if p.center.X == a.x && p.center.Y == a.y {
		//p.body = append(p.body, Part{a.x + 1, a.y + 1})
		return true
	}
	return false
}

func (p *Player) dead() bool {
	//if p.body.X == p.body.X && p.body.Y == p.body.Y {
	//	return true
	//}
	return false
}

func (p *Player) update(dir Direction) {
	p.center.X += p.dx
	p.center.Y += p.dy

	if p.center.X >= float32(WorldMap.World.ScreenWidth) {
		p.center.X = float32(WorldMap.World.ScreenWidth)
	}
	if p.center.X <= 0 {
		p.center.X = 0
	}
	if p.center.Y >= float32(WorldMap.World.ScreenHeight) {
		p.center.Y = float32(WorldMap.World.ScreenHeight)
	}
	if p.center.Y <= 0 {
		p.center.Y = 0
	}

	switch dir {
	case LEFT:
		p.angle -= WorldMap.Const.RotationSpeed
		break
	case RIGHT:
		p.angle += WorldMap.Const.RotationSpeed
		break
	case UP:
		//if p.a < 1 {
		//	p.a += 1
		//}
		p.dx = -WorldMap.Const.ObjectSpeed * float32(math.Sin(p.angle))
		p.dy = WorldMap.Const.ObjectSpeed * float32(math.Cos(p.angle))
	case DOWN:
		//if p.a > 0 {
		//	p.a -= 0
		//}
		p.dx = WorldMap.Const.ObjectSpeed * float32(math.Sin(p.angle))
		p.dy = -WorldMap.Const.ObjectSpeed * float32(math.Cos(p.angle))
	case IDLE:
		p.a = 0
		p.dx = 0
		p.dy = 0
	}

	for _, ray := range p.rays {
		ray.start.X += p.dx
		ray.start.Y += p.dy
		ray.end.X += p.dx
		ray.end.Y += p.dy
		//ray.start = sdl.FPoint{X: p.center.X, Y: p.center.Y}

		//for i, wall := range walls {
		//	collided, point := collisionLineLine(wall, ray)
		//	if collided {
		//		ray.end.X = point.X + p.dx
		//		ray.end.Y = point.Y + p.dy
		//		fmt.Printf("(%d) %t {%f, %f}\n", i, collided, point.X, point.Y)
		//	}
		//}
		//fmt.Println(ray)
	}

	if p.angle > math.Pi {
		p.angle -= 2 * math.Pi
	}
	if p.angle < -math.Pi {
		p.angle += 2 * math.Pi
	}
}

func collisionLineLine(wall *Wall, ray *Ray) (cl bool, pt sdl.FPoint) {
	a := wall.start
	b := wall.end
	c := ray.start
	d := ray.end

	var den, uA, uB float32
	den = (b.X-a.X)*(d.Y-c.Y) - (b.Y-a.Y)*(d.X-c.X)
	uA = (a.Y-c.Y)*(d.X-c.X) - (a.X-c.X)*(d.Y-c.Y)
	uB = (a.Y-c.Y)*(b.X-a.X) - (a.X-c.X)*(b.Y-a.Y)

	t := uA / den
	u := uB / den

	if den == 0 {
		return false, pt
	}
	if (t >= 0 && t <= 1) && (u >= 0 && u <= 1) {
		pt.X = a.X + t*(b.X-a.X)
		pt.Y = a.Y + t*(b.Y-a.Y)
		return true, pt
	}

	return cl, pt
}
