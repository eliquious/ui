package ui

import (
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// TextOptions contains the options for a text component.
type TextOptions struct {
	Font             string
	FontSize         float64
	TextColor        color.Color
	BackgroundColor  color.Color
	Padding          Quad
	CenterX, CenterY bool
}

// Text creates a new component for rentering text.
func Text(msg string, x, y int, opts *TextOptions) Component {
	if opts.TextColor == nil {
		opts.TextColor = color.Black
	}
	if opts.FontSize == 0 {
		opts.FontSize = 12
	}

	// load text
	ff, err := NewFontFace(opts.Font, opts.FontSize)
	if err != nil {
		log.Fatalf("failed to load font: %s err=%s", opts.Font, err)
	}

	// text bounds
	bounds := text.BoundString(ff, strings.ToUpper(msg))

	// container image
	tImage := ebiten.NewImage(
		bounds.Dx()+opts.Padding.Left+opts.Padding.Right,
		bounds.Dy()+opts.Padding.Top+opts.Padding.Bottom,
	)

	// background color
	if opts.BackgroundColor != nil {
		tImage.Fill(opts.BackgroundColor)
	}

	// draw text
	text.Draw(tImage, msg, ff, -int(bounds.Min.X)+opts.Padding.Left, -int(bounds.Min.Y)+opts.Padding.Top, opts.TextColor)

	return Image(tImage, &ImageOptions{
		X:       float64(x - opts.Padding.Left),
		Y:       float64(y - opts.Padding.Top),
		CenterX: opts.CenterX,
		CenterY: opts.CenterY,
	})
}

// DynamicText craetes a dynamic text component.
func DynamicText(opts *TextOptions) *DynamicTextComponent {
	if opts.TextColor == nil {
		opts.TextColor = color.Black
	}
	if opts.FontSize == 0 {
		opts.FontSize = 12
	}

	// load text
	ff, err := NewFontFace(opts.Font, opts.FontSize)
	if err != nil {
		log.Fatalf("failed to load font: %s err=%s", opts.Font, err)
	}

	return &DynamicTextComponent{
		opts:     opts,
		fontFace: ff,
		tImage:   ebiten.NewImage(1, 1),
	}
}

// DynamicTextComponent is a text component that is optimized for dynamic text.
type DynamicTextComponent struct {
	opts     *TextOptions
	fontFace font.Face
	tImage   *ebiten.Image

	dirty bool
	text  string
	x, y  float64
}

// SetText updates the text.
func (d *DynamicTextComponent) SetText(s string) *DynamicTextComponent {
	if s != d.text {
		d.dirty = true
		d.text = s
	}
	return d
}

// SetPosition updates the position.
func (d *DynamicTextComponent) SetPosition(x, y float64) *DynamicTextComponent {
	d.x, d.y = x, y
	return d
}

// Update updates the internal image.
func (d *DynamicTextComponent) Update(ctx *UpdateContext) error {
	if d.dirty {

		// text bounds
		bounds := text.BoundString(d.fontFace, strings.ToUpper(d.text))

		// container image
		d.tImage = ebiten.NewImage(
			bounds.Dx()+d.opts.Padding.Left+d.opts.Padding.Right,
			bounds.Dy()+d.opts.Padding.Top+d.opts.Padding.Bottom,
		)

		// background color
		if d.opts.BackgroundColor != nil {
			d.tImage.Fill(d.opts.BackgroundColor)
		}

		// draw text
		text.Draw(d.tImage, d.text, d.fontFace, -bounds.Min.X+d.opts.Padding.Left, -bounds.Min.Y+d.opts.Padding.Top, d.opts.TextColor)
	}

	return nil
}

// Display renders the text component
func (d *DynamicTextComponent) Display(ctx *DisplayContext) {
	iw, ih := d.tImage.Size()
	dx, dy := 0., 0.
	if d.opts.CenterX || d.opts.CenterY {
		dx, dy = float64(-iw/2), float64(-ih/2)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(d.x, d.y)

	// center image
	if d.opts.CenterX && d.opts.CenterY {
		op.GeoM.Translate(dx, dy)
	} else if d.opts.CenterY {
		op.GeoM.Translate(0, dy)
	} else if d.opts.CenterX {
		op.GeoM.Translate(dx, 0)
	}

	ctx.DrawImage(d.tImage, op)
}
