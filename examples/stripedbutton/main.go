package main

import (
	"context"
	"fmt"
	"image/color"
	"log"

	"github.com/eliquious/ui"
)

func main() {
	ui.EnableHighDPI()

	screenWidth, screenHeight := 512, 384
	rX, rY, rWidth, rHeight := 64, 64, 256, 64
	r := ui.Rect(rX, rY, rWidth, rHeight)

	ctx := context.Background()
	display := ui.New(ctx, &ui.DisplaySettings{
		Title:  "Testing",
		Width:  screenWidth,
		Height: screenHeight,
		Debug:  true,
	})

	display.SetBackground(ui.SolidBackground(color.White))

	display.Add(ui.StripedRect(r, &ui.StripedRectOptions{
		StripeColor: color.RGBA{255, 255, 255, 255},
		Stroke:      8,
		Angle:       45,
		Padding:     ui.UniformQuad(6),
	}))

	display.Add(ui.InteractiveComponent(r,
		ui.Rectangle(r, &ui.RectangleOptions{
			Border:    ui.StrokeBorder(color.Black, 2),
			FillColor: &color.RGBA{250, 250, 250, 255},
		}),
		ui.Rectangle(r, &ui.RectangleOptions{FillColor: &color.RGBA{0, 0, 0, 10},
			Border: ui.StrokeBorder(color.Black, 2),
		}),
		ui.Rectangle(r, &ui.RectangleOptions{FillColor: &color.RGBA{225, 225, 225, 255},
			Border: ui.StrokeBorder(color.Black, 4),
		}),
		func() {
			fmt.Println("Button clicked")
		},
	))

	display.Add(ui.StripedRect(ui.Rect(rX, rY+96, rWidth, rHeight), &ui.StripedRectOptions{
		StripeColor: color.Black,
		Stroke:      8,
		Angle:       45,
		Padding:     ui.UniformQuad(5),
	}))

	if err := display.Show(); err != nil {
		log.Fatal(err)
	}
}
