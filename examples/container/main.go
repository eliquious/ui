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

	ctx := context.Background()
	screenWidth, screenHeight := 512, 512
	display := ui.New(ctx, &ui.DisplaySettings{
		Title:  "Mouse",
		Width:  screenWidth,
		Height: screenHeight,
		// HideCursor:      true,
		BackgroundColor: color.White,
	})
	display.Add(ui.Container(ui.Rect(0, 0, 512, 512), &ui.ContainerOptions{
		Margin:    ui.UniformQuad(32),
		Border:    ui.StrokeBorder(color.Black, 8),
		FillColor: color.RGBA{255, 0, 0, 50},
		// CenterX:   true,
		// CenterY:   true,
		// Padding: ui.UniformQuad(8),
	},
		ui.Text("01010", 248, 248, &ui.TextOptions{
			Font:      "letter-goth-std-med.otf",
			FontSize:  48,
			CenterX:   true,
			CenterY:   true,
			TextColor: color.Black,
			// BackgroundColor: color.Black,
			// Padding: ui.UniformQuad(4),
		}),
		ui.Rectangle(ui.Rect(248, 248, 32, 32), &ui.RectangleOptions{
			Border:    ui.StrokeBorder(color.Black, 2),
			FillColor: color.RGBA{255, 0, 0, 50},
			CenterX:   true,
			CenterY:   true,
		}),
		ui.Rectangle(ui.Rect(248, 248, 64, 64), &ui.RectangleOptions{
			Border:    ui.StrokeBorder(color.Black, 2),
			FillColor: color.RGBA{255, 0, 0, 50},
			CenterX:   true,
			CenterY:   true,
		}),
		ui.Rectangle(ui.Rect(248, 248, 128, 128), &ui.RectangleOptions{
			Border:    ui.StrokeBorder(color.Black, 2),
			FillColor: color.RGBA{255, 0, 0, 50},
			CenterX:   true,
			CenterY:   true,
		}),
		ui.Rectangle(ui.Rect(248, 248, 256, 256), &ui.RectangleOptions{
			Border:    ui.StrokeBorder(color.Black, 2),
			FillColor: color.RGBA{255, 0, 0, 50},
			CenterX:   true,
			CenterY:   true,
		}),
		ui.Rectangle(ui.Rect(248, 248, 360, 360), &ui.RectangleOptions{
			Border:    ui.StrokeBorder(color.Black, 2),
			FillColor: color.RGBA{255, 0, 0, 50},
			CenterX:   true,
			CenterY:   true,
		}),
	))

	offsetX, offsetY := 32, 512
	for i := 0; i < 12; i++ {
		offsetX += 160

		display.Add(ui.Container(ui.Rect(offsetX-160, offsetY, 128, 128), &ui.ContainerOptions{
			Border: ui.StrokeBorder(color.Black, 4),
		},
			ui.Text(fmt.Sprintf("CH%d", i), 56, 56, &ui.TextOptions{
				Font:      "letter-goth-std-med.otf",
				FontSize:  32,
				CenterX:   true,
				CenterY:   true,
				TextColor: color.Black,
			}),
		))

		if i%4 == 3 {
			offsetY += 128 + 32
			offsetX = 32
		}
	}

	// display.Add(ui.Container(ui.Rect(32, 512+128, 128, 128), &ui.ContainerOptions{
	// 	Border:    ui.StrokeBorder(color.Black, 4),
	// 	FillColor: color.RGBA{255, 0, 0, 50},
	// },
	// 	ui.Text("CH1", 56, 56, &ui.TextOptions{
	// 		Font:      "letter-goth-std-med.otf",
	// 		FontSize:  32,
	// 		CenterX:   true,
	// 		CenterY:   true,
	// 		TextColor: color.Black,
	// 	}),
	// ))

	display.Add(ui.FPSDisplay())

	if err := display.Show(); err != nil {
		log.Fatal(err)
	}
}
