package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// RGBA converts a color.Color to color.RGBA
func RGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// FPSDisplay renders the FPS and mouse position in the top left corner in Arial font.
func FPSDisplay() Component {
	size := 12.
	if IsHighDPIEnabled() {
		size = 24.
	}

	return &fpsDisplay{DynamicText(&TextOptions{
		Font:      "arial.ttf",
		FontSize:  size,
		TextColor: color.White,
		Padding:   UniformQuad(12),
		// BackgroundColor: color.White,
	})}
}

type fpsDisplay struct {
	component *DynamicTextComponent
}

// OnMouseMove updates the text on mouse move.
func (f *fpsDisplay) OnMouseMove(x, y int) {
	mouseX, mouseY := ebiten.CursorPosition()
	msg := fmt.Sprintf("FPS: %.2f TPS: %.2f X: %d Y: %d", ebiten.CurrentFPS(), ebiten.CurrentTPS(), mouseX, mouseY)
	f.component.SetText(msg)
}

func (f *fpsDisplay) Update(ctx *UpdateContext) error {
	return f.component.Update(ctx)
}

func (f *fpsDisplay) Display(ctx *DisplayContext) {
	f.component.Display(ctx)
}
