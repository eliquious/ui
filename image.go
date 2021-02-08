package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ImageOptions is the options for the image component.
type ImageOptions struct {
	X, Y             float64
	CenterX, CenterY bool
}

// Image draws a static image at the same location.
func Image(img *ebiten.Image, opts *ImageOptions) Component {
	i := DynamicImage(opts)
	i.SetImage(img)
	return i
}

// DynamicImage creates a new DynamicImage component.
func DynamicImage(opts *ImageOptions) *DynamicImageComponent {
	// fmt.Printf("Offset X=%.f, Offset Y=%.f\n", opts.X, opts.Y)
	return &DynamicImageComponent{nil, opts}
}

// DynamicImageComponent is an image component in which the image can be updated.
type DynamicImageComponent struct {
	internal *ebiten.Image
	opts     *ImageOptions
}

// SetImage sets the image to be rendered.
func (d *DynamicImageComponent) SetImage(i *ebiten.Image) {
	d.internal = i
}

// Update is a no op.
func (d *DynamicImageComponent) Update(ctx *UpdateContext) error {
	return nil
}

// Display renders the image to the parent
func (d *DynamicImageComponent) Display(ctx *DisplayContext) {
	if d.internal == nil {
		return
	}

	iw, ih := d.internal.Size()
	dx, dy := d.opts.X, d.opts.Y

	if d.opts.CenterX {
		dx -= float64(iw / 2)
	}
	if d.opts.CenterY {
		dy -= float64(ih / 2)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(dx, dy)
	ctx.DrawImage(d.internal, op)
}
