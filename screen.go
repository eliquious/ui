package ui

import (
	"context"
	"errors"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	// HighDPI stores whether the screen has the display scaled by the DeviceScaleFactor
	highDPI bool
)

// EnableHighDPI enables high DPI display
func EnableHighDPI() {
	highDPI = true
}

// IsHighDPIEnabled returned if display has high DPI enabled.
func IsHighDPIEnabled() bool {
	return highDPI
}

// DisplaySettings stores the display settings.
type DisplaySettings struct {
	Title           string
	Width           int
	Height          int
	BackgroundColor color.Color
	HideCursor      bool
	Debug           bool
}

// New creates a new screen
func New(ctx context.Context, settings *DisplaySettings) *Display {
	ebiten.SetWindowSize(settings.Width, settings.Height)
	ebiten.SetWindowTitle(settings.Title)

	// hide cursor
	if settings.HideCursor {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	}

	display := &Display{ctx: ctx, settings: settings, mouseEventRegistry: DefaultMouseEventRegistry}
	return display
}

// Display represents a display screen
type Display struct {
	ctx                context.Context
	settings           *DisplaySettings
	mouseEventRegistry *MouseEventRegistry

	cursor         Component
	background     Component
	components     []Component
	updateHandlers []UpdateHandler
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (d *Display) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	// scale the display if highdpi is enabled
	if highDPI {
		scale := ebiten.DeviceScaleFactor()
		return int(float64(outsideWidth) * scale), int(float64(outsideHeight) * scale)
	}
	return outsideWidth, outsideHeight
}

// Add adds display components to the display.
func (d *Display) Add(c ...Component) *Display {
	d.components = append(d.components, c...)

	// add mouse handlers
	for i := 0; i < len(c); i++ {
		if h, ok := c[i].(MouseButtonHandler); ok {
			d.AddMouseButtonHandler(h)
		}
		if h, ok := c[i].(MouseMoveHandler); ok {
			d.AddMouseMoveHandler(h)
		}
	}
	return d
}

// AddMouseButtonHandler adds a mouse handler to the screen.
func (d *Display) AddMouseButtonHandler(h MouseButtonHandler) *Display {
	d.mouseEventRegistry.AddButtonHandler(h)
	return d
}

// AddUpdateHandler adds a function to be called on every update.
func (d *Display) AddUpdateHandler(h UpdateHandler) *Display {
	d.updateHandlers = append(d.updateHandlers, h)
	return d
}

// AddMouseMoveHandler adds a mouse move handler to the screen.
func (d *Display) AddMouseMoveHandler(h MouseMoveHandler) *Display {
	d.mouseEventRegistry.AddMoveHandler(h)
	return d
}

// SetCursor sets the display component for the cursor.
func (d *Display) SetCursor(c Component) *Display {
	d.cursor = c
	return d
}

// SetBackground sets the display component for the background.
func (d *Display) SetBackground(c Component) *Display {
	d.background = c
	return d
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (d *Display) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("user exit")
	}
	ctx := NewUpdateContext(d.ctx)

	// update the mouse event registry
	d.mouseEventRegistry.Update()

	// call all update handlers
	for i := 0; i < len(d.updateHandlers); i++ {
		if err := d.updateHandlers[i].Update(ctx); err != nil {
			return err
		}
	}

	// draw background
	if d.background != nil {
		if err := d.background.Update(ctx); err != nil {
			return err
		}
	}

	// update components
	for _, c := range d.components {
		if err := c.Update(ctx); err != nil {
			return err
		}
	}

	// draw cursor
	if d.cursor != nil {
		if err := d.cursor.Update(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (d *Display) Draw(screen *ebiten.Image) {
	if d.settings.BackgroundColor != nil {
		screen.Fill(d.settings.BackgroundColor)
	}

	ctx := NewDisplayContext(d.ctx, screen)
	// draw background component
	if d.background != nil {
		d.background.Display(ctx)
	}

	// draw components
	for _, c := range d.components {
		c.Display(ctx)
	}

	// draw cursor
	if d.cursor != nil {
		d.cursor.Display(ctx)
	}

	if d.settings.Debug {
		mouseX, mouseY := ebiten.CursorPosition()
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f MouseX: %d MouseY: %d", ebiten.CurrentFPS(), mouseX, mouseY))
	}
}

// Show shows the window
func (d *Display) Show() error {
	return ebiten.RunGame(d)
}
