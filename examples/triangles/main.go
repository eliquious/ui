package main

import (
	"context"
	"image/color"
	"log"

	"github.com/eliquious/ui"
)

const (
	screenWidth, screenHeight = 512, 512
)

func main() {
	ui.EnableHighDPI()

	ctx := context.Background()
	display := ui.New(ctx, &ui.DisplaySettings{
		Title:  "Triangle Testing",
		Width:  screenWidth,
		Height: screenHeight,
	})
	display.Add(ui.FPSDisplay())
	display.SetBackground(ui.SolidBackground(color.White))

	dCircle0 := ui.DynamicCircle(&ui.CircleOptions{
		FillColor: color.RGBA{255, 0, 0, 50},
		Stroke:    ui.Stroke{Color: color.Black, Width: 4},
		Radius:    32,
	})
	dCircle1 := ui.DynamicCircle(&ui.CircleOptions{
		FillColor: color.RGBA{255, 0, 0, 50},
		Stroke:    ui.Stroke{Color: color.Black, Width: 4},
		Radius:    64,
	})
	dCircle2 := ui.DynamicCircle(&ui.CircleOptions{
		FillColor: color.RGBA{255, 0, 0, 50},
		Stroke:    ui.Stroke{Color: color.Black, Width: 4},
		Radius:    128,
	})
	dCircle3 := ui.DynamicCircle(&ui.CircleOptions{
		FillColor: color.RGBA{255, 0, 0, 50},
		Stroke:    ui.Stroke{Color: color.Black, Width: 4},
		Radius:    256,
	})

	display.AddUpdateHandler(ui.UpdateHandlerFunc(func(ctx *ui.UpdateContext) error {
		mouseX, mouseY := ctx.CursorPosition()
		cx, cy := float64(mouseX), float64(mouseY)
		dCircle0.SetPosition(cx, cy)
		dCircle1.SetPosition(cx, cy)
		dCircle2.SetPosition(cx, cy)
		dCircle3.SetPosition(cx, cy)
		return nil
	}))
	display.Add(dCircle0, dCircle1, dCircle2, dCircle3)

	if err := display.Show(); err != nil {
		log.Fatal(err)
	}
}
