package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// ContainerOptions stores the options for a container.
type ContainerOptions struct {
	FillColor color.Color
	CenterX   bool
	CenterY   bool
	Margin    Quad
	Border    Border
	Padding   Quad
}

// Container creates a container component.
func Container(r image.Rectangle, opts *ContainerOptions, children ...Component) Component {
	x, y, w, h := r.Min.X, r.Min.Y, r.Dx(), r.Dy()
	borderRect := Rect(
		0, 0,
		w-opts.Margin.Left-opts.Margin.Right-opts.Border.Left.Width-opts.Border.Right.Width,
		h-opts.Margin.Top-opts.Margin.Bottom-opts.Border.Top.Width-opts.Border.Bottom.Width,
	)

	rect := Rectangle(borderRect, &RectangleOptions{
		FillColor: opts.FillColor,
		CenterX:   opts.CenterX,
		CenterY:   opts.CenterY,
		Border:    opts.Border,
	})
	// fmt.Printf("Margin=(%s) Padding=(%s)\n", opts.Margin, opts.Padding)
	// fmt.Printf("X=%d, Y=%d, W=%d, H=%d\n", x, y, w, h)

	internalComponents := StackedComponent(children...)

	internalLeft := opts.Border.Left.Width + opts.Padding.Left
	internalTop := opts.Border.Top.Width + opts.Padding.Right

	internalRect := Rect(
		internalLeft,
		internalTop,
		borderRect.Dx()-opts.Padding.Left-opts.Padding.Right,
		borderRect.Dy()-opts.Padding.Top-opts.Padding.Bottom,
	)

	// fmt.Printf("Offset X=%d, Offset Y=%d\n", internalLeft, internalTop)

	interiorWidth := internalRect.Dx()
	interiorHeight := internalRect.Dy()

	bufferImage := ebiten.NewImage(interiorWidth, interiorHeight)
	// fmt.Printf("Interior W=%d, Interior H=%d\n", interiorWidth, interiorHeight)

	// interiorImage := DynamicImage(&ImageOptions{
	// 	X:       float64(internalRect.Min.X),
	// 	Y:       float64(internalRect.Min.Y),
	// 	CenterX: opts.CenterX,
	// 	CenterY: opts.CenterY,
	// 	// X: float64(x + opts.Margin.Left), //+ opts.Border.Left.Width
	// 	// Y: float64(y + opts.Margin.Top),  //+ opts.Border.Top.Width
	// })

	return ComponentFunc(func(ctx *UpdateContext) error {
		return internalComponents.Update(ctx)
	}, func(ctx *DisplayContext) {
		ctx = ctx.Translate(float64(x), float64(y))

		// margin context
		marginContext := ctx.Translate(float64(opts.Margin.Left), float64(opts.Margin.Top))
		rect.Display(marginContext)

		bufferImage.Fill(color.Transparent)
		bufferCtx := NewDisplayContext(ctx.Context(), bufferImage)

		// render sub components onto buffered image
		internalComponents.Display(bufferCtx)

		// render buffered image
		paddingCtx := marginContext.Translate(float64(internalRect.Min.X), float64(internalRect.Min.Y))
		paddingCtx.DrawImage(bufferImage, &ebiten.DrawImageOptions{})
	})
}
