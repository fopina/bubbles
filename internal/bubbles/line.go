package bubbles

import (
	"math"
	"github.com/veandco/go-sdl2/sdl"
)

type Line struct {
	Index,
	X,
	Y,
	LineLength int32
	PopDistanceReturn,
	PopDistance,
	MaxPopDistance float64
	InversePop bool
	Bubble *Bubble
}

func NewLine(bubble *Bubble, index int32) *Line {
	l := &Line{}
	l.Index = index
	l.Bubble = bubble
	l.Reset()
	return l
}

func (l *Line) Reset() {
	l.LineLength = 0
	l.PopDistance = 0
	l.MaxPopDistance = float64(l.Bubble.Radius) * 0.5
	l.PopDistanceReturn = 0
	l.InversePop = false
}

func (l *Line) Update() {
	l.X = int32(float64(l.Bubble.X) + (float64(l.Bubble.Radius) + l.PopDistanceReturn) * math.Cos(2.0 * math.Pi * float64(l.Index) / float64(PopLines)))
	l.Y = int32(float64(l.Bubble.Y) + (float64(l.Bubble.Radius) + l.PopDistanceReturn) * math.Sin(2.0 * math.Pi * float64(l.Index) / float64(PopLines)))
	l.LineLength = int32(float64(l.Bubble.Radius) * l.PopDistance)
}

func (l *Line) Render(renderer *sdl.Renderer) {
	l.Update()
	endx := l.LineLength
	endy := l.LineLength

	if l.X < l.Bubble.X {
		endx = -l.LineLength
	}
	if l.Y < l.Bubble.Y {
		endy = -l.LineLength
	}
	if l.X == l.Bubble.X {
		endx = 0
	}
	if l.Y == l.Bubble.Y {
		endy = 0
	}
	renderer.DrawLine(
		l.X, l.Y,
		l.X + endx,
		l.Y + endy,
	)
}
