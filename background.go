package ui

import (
	"image/color"
)

// SolidBackground fills the background of the parent
func SolidBackground(c color.Color) Component {
	return SimpleComponent(func(ctx *DisplayContext) {
		ctx.Image().Fill(c)
	})
}
