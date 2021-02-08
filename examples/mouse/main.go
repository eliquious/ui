package main

import (
	"math"

	"github.com/eliquious/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"context"
	"image/color"
	"log"
)

func main() {
	ui.EnableHighDPI()

	screenWidth, screenHeight := 512, 384
	display := ui.New(context.Background(), &ui.DisplaySettings{
		Title:  "Mouse",
		Width:  screenWidth,
		Height: screenHeight,
		// HideCursor:      true,
		BackgroundColor: color.White,
	})
	display.Add(ui.FPSDisplay())

	display.Add(ui.SimpleComponent(func(ctx *ui.DisplayContext) {
		c := ui.Circle(float64(512), float64(384), 256., &ui.CircleOptions{
			Stroke: ui.Stroke{
				Color: color.Black,
				Width: 2,
			},
		})
		c.Display(ctx)
	}))

	drawCircle := func(img *ebiten.Image, r float64, c color.Color) {
		path := &vector.Path{}

		segments := 360
		path.MoveTo(float32(r*2), float32(r))

		radStep := 2 * math.Pi / float64(segments)
		for i := 0; i < segments; i++ {
			rad := radStep * float64(i)
			x1, y1 := math.Cos(rad)*float64(r), math.Sin(rad)*float64(r)
			path.LineTo(float32(x1+r), float32(y1+r))
		}
		path.Fill(img, &vector.FillOptions{
			Color: c,
		})
	}

	r := 64.
	buffered := ebiten.NewImage(128, 128)
	drawCircle(buffered, r, color.Black)

	bufferedInside := ebiten.NewImage(128, 128)
	drawCircle(bufferedInside, r-2, color.White)

	// mask
	maskedFgImage := ebiten.NewImage(128, 128)
	maskedFgImage.Fill(color.White)
	op := &ebiten.DrawImageOptions{}
	op.CompositeMode = ebiten.CompositeModeCopy
	maskedFgImage.DrawImage(buffered, op)

	centerOp := &ebiten.DrawImageOptions{}
	centerOp.GeoM.Translate(float64(2), float64(2))
	centerOp.CompositeMode = ebiten.CompositeModeSourceIn
	maskedFgImage.DrawImage(bufferedInside, centerOp)
	// buffered.DrawImage(bufferedInside, centerOp)

	// op = &ebiten.DrawImageOptions{}
	// op.CompositeMode = ebiten.CompositeModeSourceIn
	// maskedFgImage.DrawImage(bufferedInside, op)

	display.Add(ui.SimpleComponent(func(ctx *ui.DisplayContext) {
		x, y := ctx.CursorPosition()
		// buffered.Fill(color.Transparent)
		// bufferedInside.Fill(color.Transparent)
		// path := &vector.Path{}
		// r := 64.
		// segments := 360

		// // offsetX, offsetY := float64(screenWidth), float64(screenHeight)
		// path.MoveTo(float32(r*2), float32(r))

		// radStep := 2 * math.Pi / float64(segments)
		// for i := 0; i < segments; i++ {
		// 	rad := radStep * float64(i)
		// 	x1, y1 := math.Cos(rad)*float64(r), math.Sin(rad)*float64(r)
		// 	// path.LineTo(float32(x1+offsetX), float32(y1+offsetY))
		// 	path.LineTo(float32(x1+r), float32(y1+r))

		// 	// x2, y2 := math.Cos(rad+radStep)*r, math.Sin(rad+radStep)*r

		// 	// line.Draw(ctx, x1+r, y1+r, x2+r, y2+r)
		// }
		// path.Fill(buffered, &vector.FillOptions{
		// 	Color: color.Black,
		// })

		// drawCircle(bufferedInside, r-4, color.Black)

		// centerOp := &ebiten.DrawImageOptions{}
		// centerOp.GeoM.Translate(float64(4), float64(4))
		// centerOp.CompositeMode = ebiten.CompositeModeCopy
		// buffered.DrawImage(bufferedInside, centerOp)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x)-r, float64(y)-r)
		// op.GeoM.Translate(float64(screenWidth)-r, float64(screenHeight)-r)

		ctx.DrawImage(buffered, op)
		ctx.DrawImage(maskedFgImage, op)
		// ctx.DrawImage(buffered, op)
	}))

	// button handler
	display.AddMouseButtonHandler(ui.MouseHandlerFunc(func(x, y int, evt ui.MouseEvent) {
		log.Printf("X: %d Y: %d Event: %s\n", x, y, evt)
	}))

	if err := display.Show(); err != nil {
		log.Fatal(err)
	}
}
