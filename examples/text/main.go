package main

import (
	"context"
	"image/color"
	"log"

	"github.com/eliquious/ui"
)

const (
	screenWidth, screenHeight = 512, 384
)

func main() {
	ui.EnableHighDPI()

	ctx := context.Background()
	display := ui.New(ctx, &ui.DisplaySettings{
		Title:  "Testing",
		Width:  screenWidth,
		Height: screenHeight,
		Debug:  true,
		// HideCursor: true,
		// BackgroundColor: color.RGBA{26, 27, 31, 255},
	})

	display.SetBackground(ui.SolidBackground(color.White))

	display.Add(ui.Text("Proverbs 21:16", 128, 128, &ui.TextOptions{
		Font:      "letter-goth-std-bold.otf",
		FontSize:  48,
		TextColor: color.Black,
	}))
	display.Add(ui.Text("Whoever wanders from the way of \nunderstanding will rest in the\nassembly of the dead.", 128, 192, &ui.TextOptions{
		Font:      "letter-goth-std-med.otf",
		FontSize:  32,
		TextColor: color.Black,
	}))

	display.Add(ui.Text(`
    DATE   |       DESCRIPTION        |  CV2  | AMOUNT
-----------+--------------------------+-------+----------
  1/1/2014 | Domain name              |  2233 | $10.98
  1/1/2014 | January Hosting          |  2233 | $54.95
  1/4/2014 | February Hosting         |  2233 | $51.00
  1/4/2014 | February Extra Bandwidth |  2233 | $30.00
-----------+--------------------------+-------+----------
                                        TOTAL | $146.93
                                      --------+----------`,
		32, 384, &ui.TextOptions{
			Font:      "letter-goth-std-med.otf",
			FontSize:  24,
			TextColor: color.Black,
		}))

	if err := display.Show(); err != nil {
		log.Fatal(err)
	}
}
