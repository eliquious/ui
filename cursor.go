package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// TriangleCursor renders a triangle at the mouse position
func TriangleCursor(c color.Color) Component {
	stroke := 4
	line := DynamicAntiAliasedLine(c, stroke)

	height := 18.
	width := 10.
	if IsHighDPIEnabled() {
		scale := ebiten.DeviceScaleFactor()
		height *= scale
		width *= scale
	}

	return SimpleComponent(func(ctx *DisplayContext) {
		mouseX, mouseY := ebiten.CursorPosition()

		if mouseX > 1e4 || mouseX < -1e4 {
			return
		} else if mouseY > 1e4 || mouseY < -1e4 {
			return
		}

		line.Draw(ctx, float64(mouseX), float64(mouseY), float64(mouseX), float64(mouseY)+height)
		line.Draw(ctx, float64(mouseX+stroke), float64(mouseY+stroke/2-1), float64(mouseX+stroke)+width, float64(mouseY)+height/2)
		// line.Draw(screen, float64(mouseX), float64(mouseY)+18, float64(mouseX)+10, float64(mouseY)+9)
	})
}

// ImageCursor renders an image at the mouse position
func ImageCursor(i *ebiten.Image, center bool) Component {
	iw, ih := i.Size()

	op := &ebiten.DrawImageOptions{}
	return SimpleComponent(func(ctx *DisplayContext) {
		mouseX, mouseY := ebiten.CursorPosition()
		op.GeoM.Translate(float64(mouseX), float64(mouseY))

		if center {
			op.GeoM.Translate(float64(-iw/2), float64(-ih/2))
		}

		// draw cursor
		ctx.DrawImage(i, op)
		op.GeoM.Reset()
	})
}
