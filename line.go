package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Line draws a line.
func Line(x1, y1, x2, y2 float64, stroke int, c color.Color) Component {
	emptyImage := ebiten.NewImage(1, stroke)
	emptyImage.Fill(c)

	length := math.Hypot(x2-x1, y2-y1)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(length, 1)
	op.GeoM.Translate(0, float64(-stroke/2))
	op.GeoM.Rotate(math.Atan2(y2-y1, x2-x1))
	op.GeoM.Translate(x1, y1)

	return SimpleComponent(func(dctx *DisplayContext) {
		dctx.DrawImage(emptyImage, op)
	})
}

// StaticAntiAliasedLine creates a new anti-aliased line that is meant for static rendering.
func StaticAntiAliasedLine(x1, y1, x2, y2 float64, stroke int, c color.Color) Component {
	if stroke < 1 {
		stroke = 1
	}
	p := ebiten.NewImage(1, 1)
	p.Fill(c)
	line := &AntiAliasedLine{p, stroke}
	return SimpleComponent(func(ctx *DisplayContext) {
		line.Draw(ctx, x1, y1, x2, y2)
	})
}

// DynamicAntiAliasedLine creates a new anti-aliased line that is designed for dynamic rendering.
func DynamicAntiAliasedLine(color color.Color, stroke int) *AntiAliasedLine {
	if stroke < 1 {
		stroke = 1
	}
	p := ebiten.NewImage(1, 1)
	p.Fill(color)
	return &AntiAliasedLine{p, stroke}
}

// AntiAliasedLine is an antialiased line
type AntiAliasedLine struct {
	pixel  *ebiten.Image
	stroke int
}

// Draw draws the anti-aliased line on the image
func (l *AntiAliasedLine) Draw(ctx *DisplayContext, x0, y0, x1, y1 float64) {

	dx := x1 - x0
	dy := y1 - y0

	if math.Abs(dx) > 1e4 {
		return
	} else if math.Abs(dy) > 1e4 {
		return
	}

	ax := dx
	if ax < 0 {
		ax = -ax
	}
	ay := dy
	if ay < 0 {
		ay = -ay
	}
	// plot function set here to handle the two cases of slope
	op := &ebiten.DrawImageOptions{}
	var plot func(int, int, float64)
	if ax < ay {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
		dx, dy = dy, dx
		plot = func(x, y int, c float64) {
			op.ColorM.Scale(1, 1, 1, c)

			// x/y are intensionally switched
			op.GeoM.Translate(float64(y), float64(x))
			ctx.DrawImage(l.pixel, op)

			op.GeoM.Reset()
			op.ColorM.Reset()
		}
	} else {
		plot = func(x, y int, c float64) {
			// l.point(screen, x, y, c)

			op.ColorM.Scale(1, 1, 1, c)
			op.GeoM.Translate(float64(x), float64(y))
			ctx.DrawImage(l.pixel, op)
			op.GeoM.Reset()
			op.ColorM.Reset()
		}
	}
	if x1 < x0 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	gradient := dy / dx

	// handle first endpoint
	xend := l.round(x0)
	yend := y0 + gradient*(xend-x0)
	xgap := l.rfpart(x0 + .5)
	xpxl1 := int(xend) // this will be used in the main loop
	ypxl1 := int(l.ipart(yend))

	intery := yend + gradient // first y-intersection for the main loop

	// handle second endpoint
	xend = l.round(x1)
	yend = y1 + gradient*(xend-x1)
	xgap = l.fpart(x1 + 0.5)
	xpxl2 := int(xend) // this will be used in the main loop
	ypxl2 := int(l.ipart(yend))

	// draw endpoints
	for i := 0; i < l.stroke; i++ {
		plot(xpxl1, ypxl1+i, l.rfpart(yend)*xgap)
		plot(xpxl1, ypxl1+1+i, l.fpart(yend)*xgap)

		plot(xpxl2, i+ypxl2, l.rfpart(yend)*xgap)
		plot(xpxl2, i+ypxl2+1, l.fpart(yend)*xgap)
	}

	// main loop
	for x := xpxl1 + 1; x <= xpxl2-1; x++ {
		for i := 0; i < l.stroke; i++ {
			plot(x, int(l.ipart(intery))+i, l.rfpart(intery))
			plot(x, int(l.ipart(intery))+1+i, l.fpart(intery))
		}
		intery = intery + gradient
	}
}

// integer part of x
func (l *AntiAliasedLine) ipart(x float64) int {
	return int(math.Floor(x))
}

func (l *AntiAliasedLine) round(x float64) float64 {
	return (math.Round(x))
}

// fractional part of x
func (l *AntiAliasedLine) fpart(x float64) float64 {
	return x - float64(l.ipart(x))
}

func (l *AntiAliasedLine) rfpart(x float64) float64 {
	return 1 - l.fpart(x)
}

// VertexLine creates a line component.
func VertexLine(s Stroke) *VertexLineComponent {
	return &VertexLineComponent{0, 0, 0, 0, s}
}

// VertexLineComponent draws a line.
type VertexLineComponent struct {
	x0, y0, x1, y1 float64
	stroke         Stroke
}

// SetColor sets the color of the line.
func (t *VertexLineComponent) SetColor(c color.Color) *VertexLineComponent {
	t.stroke.Color = c
	return t
}

// SetPoints sets the points of the line.
func (t *VertexLineComponent) SetPoints(x0, y0, x1, y1 float64) *VertexLineComponent {
	t.x0, t.y0, t.x1, t.y1 = x0, y0, x1, y1
	return t
}

// SetWidth sets the width of the line.
func (t *VertexLineComponent) SetWidth(w int) *VertexLineComponent {
	t.stroke.Width = w
	return t
}

// Update is a no-op.
func (t *VertexLineComponent) Update(ctx *UpdateContext) error {
	return nil
}

// Display renders the triangles.
func (t *VertexLineComponent) Display(ctx *DisplayContext) {
	vs, is := LineVertices(t.x0, t.y0, t.x1, t.y1, t.stroke)

	op := &ebiten.DrawTrianglesOptions{}
	op.Filter = ebiten.FilterLinear
	ctx.DrawTriangles(vs, is, op)
}

// LineVertices returns the vertices for a line.
func LineVertices(px0, py0, px1, py1 float64, stroke Stroke) ([]ebiten.Vertex, []uint16) {
	x0, y0 := float32(px0), float32(py0)
	x1, y1 := float32(px1), float32(py1)

	width := float32(stroke.Width)

	r0, g0, b0, a0 := stroke.Color.RGBA()
	clr := color.RGBA{uint8(r0), uint8(g0), uint8(b0), uint8(a0)}

	theta := math.Atan2(float64(y1-y0), float64(x1-x0))
	theta += math.Pi / 2
	dx := float32(math.Cos(theta))
	dy := float32(math.Sin(theta))

	r := float32(clr.R) / 0xff
	g := float32(clr.G) / 0xff
	b := float32(clr.B) / 0xff
	a := float32(clr.A) / 0xff

	return []ebiten.Vertex{
		{
			DstX:   x0 - width*dx/2,
			DstY:   y0 - width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x0 + width*dx/2,
			DstY:   y0 + width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1 - width*dx/2,
			DstY:   y1 - width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1 + width*dx/2,
			DstY:   y1 + width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}, []uint16{0, 1, 2, 1, 2, 3}
}
