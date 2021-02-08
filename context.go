package ui

import (
	"context"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type contextKey string

// // WithBounds adds parent image bounds to the context.
// func WithBounds(ctx context.Context, r image.Rectangle) context.Context {
// 	return context.WithValue(ctx, contextKeyBounds, r)
// }

// // Size returns the size of the parent image.
// func Size(ctx context.Context) (int, int) {
// 	val := ctx.Value(contextKeyBounds)
// 	r, ok := val.(image.Rectangle)
// 	if !ok {
// 		return 0, 0
// 	}
// 	s := r.Size()
// 	return s.X, s.Y
// }

// // Bounds returns the bounds of the parent image.
// func Bounds(ctx context.Context) (image.Rectangle, bool) {
// 	r, ok := ctx.Value(contextKeyBounds).(image.Rectangle)
// 	return r, ok
// }

// NewDisplayContext creates a new display context.
func NewDisplayContext(ctx context.Context, i *ebiten.Image) *DisplayContext {
	src := ebiten.NewImage(3, 3)
	src.Fill(color.White)
	return &DisplayContext{ctx, 0, 0, i, src}
}

// DisplayContext manages drawing the component in the parent component.
type DisplayContext struct {
	context context.Context

	dx, dy     float64
	parent     *ebiten.Image
	emptyImage *ebiten.Image
}

// Image returns the image for the parent context.
func (c *DisplayContext) Image() *ebiten.Image {
	return c.parent
}

// Context returns the context.Context for the DisplayContext
func (c *DisplayContext) Context() context.Context {
	return c.context
}

// Translate creates a new context after translated.
func (c *DisplayContext) Translate(x, y float64) *DisplayContext {
	return &DisplayContext{c.context, c.dx + x, c.dy + y, c.parent, c.emptyImage}
}

// DrawImage draws the image on the parent image.
func (c *DisplayContext) DrawImage(i *ebiten.Image, op *ebiten.DrawImageOptions) {
	if c.dx > 0 || c.dy > 0 {
		op.GeoM.Translate(c.dx, c.dy)
	}
	c.parent.DrawImage(i, op)
}

// DrawTriangles draws triangles on the parent image.
func (c *DisplayContext) DrawTriangles(vs []ebiten.Vertex, is []uint16, op *ebiten.DrawTrianglesOptions) {
	src := c.emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	c.parent.DrawTriangles(vs, is, src, op)
}

// CursorPosition returns the mouse position.
func (c *DisplayContext) CursorPosition() (int, int) {
	return ebiten.CursorPosition()
}

// NewUpdateContext creates a new UpdateContext with the provided context.Context.
func NewUpdateContext(ctx context.Context) *UpdateContext {
	return &UpdateContext{ctx}
}

// UpdateContext provides a simple context and a way to pass information to child components during update.
type UpdateContext struct {
	context context.Context
}

// Context returns the context.Context.
func (u *UpdateContext) Context() context.Context {
	return u.context
}

// CursorPosition returns the mouse position.
func (u *UpdateContext) CursorPosition() (int, int) {
	return ebiten.CursorPosition()
}
