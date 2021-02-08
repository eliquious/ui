package main

import (
	"context"
	"image/color"
	"log"

	"github.com/eliquious/ui"
)

func main() {
	ui.EnableHighDPI()

	screenWidth, screenHeight := 512, 384
	display := ui.New(context.Background(), &ui.DisplaySettings{
		Title:           "Path Line",
		Width:           screenWidth,
		Height:          screenHeight,
		BackgroundColor: color.White,
	})
	display.Add(ui.FPSDisplay())

	display.Add(ui.Line(64, 64, 256, 256, 12, color.Black))
	display.Add(ui.StaticAntiAliasedLine(64, 64+32, 256, 256+32, 32, color.Black))

	if err := display.Show(); err != nil {
		log.Fatal(err)
	}
}
