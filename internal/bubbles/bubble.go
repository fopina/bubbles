package bubbles

import (
	"math/rand"
	"math"
	
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/gfx"
)

type MousePos struct {
	X,
	Y int32
}

type Bubble struct {
	X,
	Y,
	MaxX,
	MaxY,
	Radius,
	XOff,
	YOff int32
    DistanceBetweenWaves float64
    Rotation,
    RotationStep,
    Count int32
    Color sdl.Color
    Popping bool
    Lines [PopLines]*Line
}

func NewBubble(MaxX, MaxY int32) *Bubble {
	b := &Bubble{}
	b.MaxX = MaxX
	b.MaxY = MaxY
	b.Rotation = rand.Int31n(MaxRotation * 2) - MaxRotation
	b.RotationStep = BubbleSpeed
	b.Color = sdl.Color{255,255,255,255}
	for i := range b.Lines {
		b.Lines[i] = NewLine(b, int32(i))
	}
	b.Reset()
	return b
}

func (b *Bubble) Reset() {
	b.X = 0
	b.Y = 0
	b.Radius = MinRadius + rand.Int31n(MaxRadiusDelta)
	b.XOff = rand.Int31n(b.MaxX) - b.Radius
	b.YOff = rand.Int31n(b.MaxY)
	b.DistanceBetweenWaves = 50 + float64(rand.Int31n(40))
	b.Count = b.MaxY + b.YOff
	b.Popping = false
	for _, l := range b.Lines {
		l.Reset()
	}
}

func (b *Bubble) Render(renderer *sdl.Renderer) {
	sdl.Do(func() {
		if b.Popping {
			for _, l := range b.Lines {
				if float64(l.LineLength) < l.MaxPopDistance && !l.InversePop {
            		l.PopDistance += 0.06
	          	} else {
	            	if(l.PopDistance >= 0) {
	              		l.InversePop = true
	              		l.PopDistanceReturn += BubbleSpeed
	              		l.PopDistance -= 0.03
	            	} else {
	            		b.Reset()
	            	}
	            }
	            l.Render(renderer)
	        }
		} else {
			gfx.CircleColor(renderer, b.X, b.Y, b.Radius, b.Color)
			gfx.ArcColor(
				renderer, b.X, b.Y, b.Radius - 3,
				b.Rotation - MaxRotation,
				b.Rotation - MaxRotation + 90,
				b.Color,
			)
		}
	})
}

func (b *Bubble) Update(mouse *MousePos) {
	b.X = int32(math.Sin(float64(b.Count) / b.DistanceBetweenWaves) * 50) + b.XOff;
	b.Y = b.Count
	if b.Count < -b.Radius {
		b.Count = b.MaxY + b.YOff
	} else {
		b.Count -= BubbleSpeed
	}
	b.Rotation += b.RotationStep
	if math.Abs(float64(b.Rotation)) >= MaxRotation {
		b.RotationStep *= -1
	}

	if mouse.X >= b.X - b.Radius && mouse.X <= b.X + b.Radius {
		if mouse.Y >= b.Y - b.Radius && mouse.Y <= b.Y + b.Radius {
			b.Popping = true
		}
	}
}
