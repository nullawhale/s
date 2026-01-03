package main

import "math"

type Vector struct {
	X float64
	Y float64
}

type PlayerInput struct {
	Left   bool
	Right  bool
	Thrust bool
	Brake  bool
}

func (v Vector) Add(o Vector) Vector {
	return Vector{X: v.X + o.X, Y: v.Y + o.Y}
}

func (v Vector) Sub(o Vector) Vector {
	return Vector{X: v.X - o.X, Y: v.Y - o.Y}
}

func (v Vector) Scale(s float64) Vector {
	return Vector{X: v.X * s, Y: v.Y * s}
}

func (v Vector) Len() float64 {
	return math.Hypot(v.X, v.Y)
}

func (v Vector) Normalize() Vector {
	l := v.Len()
	if l == 0 {
		return Vector{}
	}
	return Vector{X: v.X / l, Y: v.Y / l}
}

func wrapPosition(pos Vector, width, height float64) Vector {
	if pos.X < 0 {
		pos.X += width
	}
	if pos.X >= width {
		pos.X -= width
	}
	if pos.Y < 0 {
		pos.Y += height
	}
	if pos.Y >= height {
		pos.Y -= height
	}
	return pos
}
