package ui

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// UniformQuad returns a Quad where each side is the same.
func UniformQuad(side int) Quad {
	return Quad{side, side, side, side}
}

// Quad stores the padding around the element.
type Quad struct {
	Left, Right, Top, Bottom int
}

// String returns the string representation.
func (q Quad) String() string {
	return fmt.Sprintf("Left=%d, Right=%d, Top=%d, Bottom=%d", q.Left, q.Right, q.Top, q.Bottom)
}

// Stroke represents a line stroke.
type Stroke struct {
	Color color.Color
	Width int
}

func (s Stroke) String() string {
	return fmt.Sprintf("Color=%s, Width=%d", s.Color, s.Width)
}

// StrokeBorder creates a border with the same stroke on all sides.
func StrokeBorder(c color.Color, width int) Border {
	stroke := Stroke{c, width}
	return Border{stroke, stroke, stroke, stroke}
}

// Border represents a border for a rectangle.
type Border struct {
	Left, Right, Top, Bottom Stroke
}

// RectangleOptions stores the rectangle options
type RectangleOptions struct {
	FillColor color.Color
	CenterX   bool
	CenterY   bool
	Margin    Quad
	Border    Border
}

// Rect creates a new image.Rectangle using X,Y,W,H coordinates.
func Rect(x, y, w, h int) image.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}

// Rectangle draws a rectangle. All border lines are within the rectangle.
func Rectangle(r image.Rectangle, opts *RectangleOptions) Component {

	// add margin
	x, y, w, h := r.Min.X, r.Min.Y, r.Dx(), r.Dy()
	borderRect := Rect(x+opts.Margin.Left, y+opts.Margin.Top, w-opts.Margin.Left-opts.Margin.Right, h-opts.Margin.Top-opts.Margin.Bottom)

	// new cached image
	rImage := ebiten.NewImage(borderRect.Dx(), borderRect.Dy())
	if opts.FillColor != nil {
		rImage.Fill(opts.FillColor)
	}

	ctx := NewDisplayContext(context.Background(), rImage)

	left := float64(opts.Border.Left.Width / 2)
	top := float64(opts.Border.Top.Width / 2)
	right := float64(borderRect.Dx() - opts.Border.Right.Width/2)
	bottom := float64(borderRect.Dy() - opts.Border.Bottom.Width/2)

	// left
	if opts.Border.Left.Width > 0 {
		strokeColor := opts.Border.Left.Color
		stroke := opts.Border.Left.Width

		line := Line(left, top, left, bottom, stroke, strokeColor)
		line.Display(ctx)
	}

	// right
	if opts.Border.Right.Width > 0 {
		strokeColor := opts.Border.Right.Color
		stroke := opts.Border.Right.Width

		line := Line(right, top, right, bottom, stroke, strokeColor)
		line.Display(ctx)
	}

	// top
	if opts.Border.Top.Width > 0 {
		strokeColor := opts.Border.Top.Color
		stroke := opts.Border.Top.Width

		line := Line(0, top, float64(borderRect.Dx()), top, stroke, strokeColor)
		line.Display(ctx)
	}

	// bottom
	if opts.Border.Bottom.Width > 0 {
		strokeColor := opts.Border.Bottom.Color
		stroke := opts.Border.Bottom.Width

		line := Line(0, bottom-1, float64(borderRect.Dx()), bottom-1, stroke, strokeColor)
		line.Display(ctx)
	}

	return Image(rImage, &ImageOptions{
		CenterX: opts.CenterX,
		CenterY: opts.CenterY,
		X:       float64(borderRect.Min.X),
		Y:       float64(borderRect.Min.Y),
	})
}

// NewSlantImage draws a slant onto an image
func NewSlantImage(angle float64, stroke, height int, c color.Color) *ebiten.Image {
	rad := angle * math.Pi / 180.
	hypot := float64(height) / math.Sin(rad)
	width := int(hypot * math.Cos(rad))
	slant := ebiten.NewImage(int(width)+stroke, height)

	path := &vector.Path{}

	// top left of stroke
	path.MoveTo(float32(width), float32(0))

	// top right
	path.LineTo(float32(int(width)+stroke), float32(0))

	// bottom slant line
	path.LineTo(float32(stroke), float32(height))

	// bottom line
	path.LineTo(float32(0), float32(height))

	// top slant line
	path.LineTo(float32(width), float32(0))

	// fill
	path.Fill(slant, &vector.FillOptions{Color: c})

	return slant
}

// StripedRectOptions is the stripe rect options
type StripedRectOptions struct {
	StripeColor     color.Color
	Stroke          int
	Angle           float64
	BackgroundColor color.Color
	Padding         Quad
}

// StripedRect draws a rectangle with striped
func StripedRect(rect image.Rectangle, opts *StripedRectOptions) Component {
	if opts.StripeColor == nil {
		opts.StripeColor = color.Black
	}
	if opts.Stroke == 0 {
		opts.Stroke = 4
	}
	if opts.Angle == 0 {
		opts.Angle = 75
	}

	x, y := rect.Min.X+opts.Padding.Left, rect.Min.Y+opts.Padding.Top

	// must account for padding on both sides
	w, h := rect.Dx()-opts.Padding.Right-opts.Padding.Left, rect.Dy()-opts.Padding.Bottom-opts.Padding.Top
	clipped := ebiten.NewImage(w, h)
	if opts.BackgroundColor != nil {
		clipped.Fill(opts.BackgroundColor)
	}

	// create slant image
	slant := NewSlantImage(opts.Angle, opts.Stroke, h, opts.StripeColor)
	sw, _ := slant.Size()

	// draw stripes
	op := &ebiten.DrawImageOptions{}
	for x := -sw + opts.Stroke; x < w+sw; x += opts.Stroke * 2 {
		op.GeoM.Translate(float64(x), 0.)
		clipped.DrawImage(slant, op)
		op.GeoM.Reset()
	}

	return Image(clipped, &ImageOptions{X: float64(x), Y: float64(y)})
}

// CircleOptions is the options for the circle.
type CircleOptions struct {
	FillColor color.Color
	Stroke    Stroke
	Radius    float64
}

// Circle draws a circle. The width and height of the rectangle should be the same.
func Circle(x, y, r float64, opts *CircleOptions) Component {
	dCircle := DynamicCircle(opts)
	dCircle.SetPosition(x, y)
	dCircle.SetRadius(r)
	return dCircle
}

// DynamicCircle creates a cicle componet which can change radius and position.
func DynamicCircle(opts *CircleOptions) *DynamicCircleComponent {
	if opts.Stroke.Color == nil {
		opts.Stroke.Color = color.Black
	}
	return &DynamicCircleComponent{opts, 0, 0, float32(opts.Radius)}
}

// DynamicCircleComponent is a circle component which can change position and radius.
type DynamicCircleComponent struct {
	opts         *CircleOptions
	x, y, radius float32
}

// SetRadius sets the radius of the circle.
func (d *DynamicCircleComponent) SetRadius(r float64) {
	d.radius = float32(r)
}

// SetPosition sets the position of the cicle..
func (d *DynamicCircleComponent) SetPosition(x, y float64) {
	d.x, d.y = float32(x), float32(y)
}

// Display draws the circle on the screen.
func (d *DynamicCircleComponent) Display(ctx *DisplayContext) {
	if d.radius == 0 {
		return
	}

	c := RGBA(d.opts.Stroke.Color)
	cr := float32(c.R) / 0xff
	cg := float32(c.G) / 0xff
	cb := float32(c.B) / 0xff
	ca := float32(c.A) / 0xff

	cx, cy := float32(d.x), float32(d.y)

	r := d.radius
	circum := 2 * math.Pi * r
	numSegments := int(circum) / 2

	if numSegments < 8 {
		numSegments = 8
	}
	theta := 2 * math.Pi / float64(numSegments)
	width := float32(d.opts.Stroke.Width)

	tanFactor := float32(math.Tan(theta)) //calculate the tangential factor
	radFactor := float32(math.Cos(theta)) //calculate the radial factor

	// inside and outside radii
	ri, ro := float32(r-width/2), float32(r+width/2)

	vertices := make([]ebiten.Vertex, 0, numSegments*2+2)
	indices := make([]uint16, 0, numSegments/2*3+2)

	// initialize at 0 deg
	xi, yi, xo, yo := ri, float32(0.), ro, float32(0.)

	vert := func(x, y float32) ebiten.Vertex {
		return ebiten.Vertex{DstX: x, DstY: y, SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}
	}

	// add initial inside/outside verts
	vertices = append(vertices, vert(xi+cx, yi+cy), vert(xo+cx, yo+cy))

	for ii := 0; ii < numSegments; ii++ {

		// inside vert
		txi := -yi
		tyi := xi
		xi += txi * tanFactor
		yi += tyi * tanFactor
		xi *= radFactor
		yi *= radFactor

		// outside vert
		txo := -yo
		tyo := xo
		xo += txo * tanFactor
		yo += tyo * tanFactor
		xo *= radFactor
		yo *= radFactor

		// add verts
		vertices = append(vertices, vert(xi+cx, yi+cy), vert(xo+cx, yo+cy))

		l := uint16(len(vertices)) - 4
		indices = append(indices, l, l+1, l+2, l+1, l+2, l+3)
	}

	op := &ebiten.DrawTrianglesOptions{Filter: ebiten.FilterNearest}

	if d.opts.Stroke.Width > 0 {
		ctx.DrawTriangles(vertices, indices, op)
	}

	if d.opts.FillColor != nil {
		c = RGBA(d.opts.FillColor)
		cr = float32(c.R) / 0xff
		cg = float32(c.G) / 0xff
		cb = float32(c.B) / 0xff
		ca = float32(c.A) / 0xff
		fillIndices := make([]uint16, 0, (numSegments)*3)
		fillVertices := make([]ebiten.Vertex, numSegments+1)
		fillVertices[0] = vert(cx, cy)
		for i := 0; i < numSegments; i++ {
			fillVertices[i+1] = vertices[i*2]
			fillVertices[i+1].ColorR = cr
			fillVertices[i+1].ColorG = cg
			fillVertices[i+1].ColorB = cb
			fillVertices[i+1].ColorA = ca
		}

		for i := 0; i < numSegments-1; i++ {
			fillIndices = append(fillIndices, 0, uint16(i), uint16(i+1))
		}
		fillIndices = append(fillIndices, 0, uint16(numSegments-1), uint16(1))
		ctx.DrawTriangles(fillVertices, fillIndices, op)
	}

	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(d.x-d.radius-float64(d.opts.Stroke.Width)/2, d.y-d.radius-float64(d.opts.Stroke.Width)/2)
	// ctx.DrawImage(d.bufferedImage, op)
}

// Update is a no-op.
func (d *DynamicCircleComponent) Update(ctx *UpdateContext) error {
	return nil
}

// TriLine draws a line using triangles.
func TriLine(x1, y1, x2, y2 float64, stroke int, c color.Color) Component {
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

// RectVertices returns the vertices for a rectangle.
func RectVertices(px0, py0, px1, py1 int, c color.Color) ([]ebiten.Vertex, []uint16) {
	x0, y0 := float32(px0), float32(py0)
	x1, y1 := float32(px1), float32(py1)

	r0, g0, b0, a0 := c.RGBA()
	clr := color.RGBA{uint8(r0), uint8(g0), uint8(b0), uint8(a0)}

	r := float32(clr.R) / 0xff
	g := float32(clr.G) / 0xff
	b := float32(clr.B) / 0xff
	a := float32(clr.A) / 0xff

	return []ebiten.Vertex{
		{
			// bottom left
			DstX:   x0,
			DstY:   y1,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			// top left
			DstX:   x0,
			DstY:   y0,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			// top right
			DstX:   x1,
			DstY:   y0,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			// bottom right
			DstX:   x1,
			DstY:   y1,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}, []uint16{0, 1, 2, 1, 2, 3}
}
