package main

import (
	"context"
	"image/color"
	"log"
	"strings"

	"github.com/eliquious/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth, screenHeight = 512, 128
)

func main() {
	ui.EnableHighDPI()

	ctx := context.Background()
	display := ui.New(ctx, &ui.DisplaySettings{
		Title:  "Keyboard Testing",
		Width:  screenWidth,
		Height: screenHeight,
		Debug:  true,
		// HideCursor: true,
		// BackgroundColor: color.RGBA{26, 27, 31, 255},
	})

	display.SetBackground(ui.SolidBackground(color.White))
	// display.Add(ui.Rectangle(ui.Rect(0, 0, screenWidth*2, screenHeight), &ui.RectangleOptions{
	// 	FillColor: color.White,
	// 	Margin:    ui.UniformQuad(32),
	// 	Border:    ui.StrokeBorder(color.Black, 2),
	// }))

	// t := ui.Text("Proverbs 21:16", 45, screenHeight-36-36, &ui.TextOptions{
	// 	Font:      "letter-goth-std-bold.otf",
	// 	FontSize:  32,
	// 	TextColor: color.Black,
	// })
	// display.Add(t)

	dt := ui.DynamicText(&ui.TextOptions{
		Font:      "letter-goth-std-med.otf",
		FontSize:  32,
		TextColor: color.Black,
		Padding:   ui.UniformQuad(8),
	}).SetText("text").SetPosition(0, 0)

	// var t ui.Component
	var text string
	display.AddUpdateHandler(ui.UpdateHandlerFunc(func(ctx *ui.UpdateContext) error {

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			text = ""
		} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			text += " "
		} else if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
			text += "    "
		}

		shift := ebiten.IsKeyPressed(ebiten.KeyShift)
		for k := ebiten.Key(0); k < ebiten.KeyZ; k++ {
			if inpututil.IsKeyJustPressed(k) {
				if shift {
					text += k.String()
				} else {
					text += strings.ToLower(k.String())
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyApostrophe) {
			text += "'"
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackslash) {
			if shift {
				text += "|"
			} else {
				text += "\\"
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			if len(text) > 4 && strings.HasSuffix(text, "    ") {
				text = text[:len(text)-4]
			} else if len(text) > 0 {
				text = text[:len(text)-1]
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyComma) {
			text += ","
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
			text += "="
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyGraveAccent) {
			text += "~"
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyLeftBracket) {
			text += "["
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
			text += "-"
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyPeriod) {
			text += "."
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRightBracket) {
			text += "]"
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySemicolon) {
			text += ";"
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
			text += "/"
		}

		dt.SetText(text)
		return nil
	}))
	display.Add(ui.Container(ui.Rect(0, 0, screenWidth*2, screenHeight), &ui.ContainerOptions{
		Margin:  ui.UniformQuad(32),
		Border:  ui.StrokeBorder(color.Black, 2),
		Padding: ui.UniformQuad(8),
	}, dt))

	// display.Add(ui.Text("Whoever wanders from the way of \nunderstanding will rest in the\nassembly of the dead.", 128, 192, &ui.TextOptions{
	// 	Font:      "letter-goth-std-med.otf",
	// 	FontSize:  32,
	// 	TextColor: color.Black,
	// }))

	if err := display.Show(); err != nil {
		log.Fatal(err)
	}
}
