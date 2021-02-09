package ui

import (
	"image"
)

// Button combines the component and mouse handlers in one interface.
type Button interface {
	Component
	MouseButtonHandler
	MouseMoveHandler
}

// ButtonState represents a button's internal state.
type ButtonState bool

// These are the two available button states.
const (
	ButtonDown ButtonState = true
	ButtonUp               = false
)

// MomentaryButton creates an interactive component which responds to mouse events. The onPress function returns the state of the button that should be set after the function had completed.
// Normally, the onPress function should always return true as the button state should remain pressed until the mouse button has released. However, during screen transitions, etc the
// button state may require the button state to remain unpressed in true momentary fashion.
func MomentaryButton(r image.Rectangle, defaultComponent Component, hover Component, pressed Component, onPress func() ButtonState, onRelease func()) Button {
	return &momentaryButton{r: r, defaultComponent: defaultComponent, hoverComponent: hover, pressedComponent: pressed, onPress: onPress, onRelease: onRelease}
}

// momentaryButton is a button which responds to momentary presses.
type momentaryButton struct {
	r                image.Rectangle
	defaultComponent Component
	hoverComponent   Component
	pressedComponent Component
	onPress          func() ButtonState
	onRelease        func()

	mouseOver bool
	pressed   ButtonState
}

// Display renders the button. If the pressed and hover components are defined they will be rendered instead of the default component.
func (i *momentaryButton) Display(ctx *DisplayContext) {
	if i.pressed {
		if i.pressedComponent != nil {
			i.pressedComponent.Display(ctx)
			return
		}
	}

	if i.mouseOver {
		if i.hoverComponent != nil {
			i.hoverComponent.Display(ctx)
			return
		}
	}
	i.defaultComponent.Display(ctx)
}

// Update is a no-op.
func (i *momentaryButton) Update(ctx *UpdateContext) error {
	return nil
}

// OnMouseEvent calls the onPress and onRelease handlers.
func (i *momentaryButton) OnMouseEvent(x, y int, evt MouseEvent) {
	if x > i.r.Min.X && x < i.r.Max.X && y > i.r.Min.Y && y < i.r.Max.Y {
		if evt.EventType == MousePressEvent {
			i.pressed = !i.pressed

			if i.onPress != nil {
				i.pressed = i.onPress()
			}
		}
	}

	// mouse release
	if evt.EventType == MouseReleaseEvent && i.pressed {
		i.pressed = false

		if i.onRelease != nil {
			i.onRelease()
		}
	}
}

// OnMouseMove toggles the mouse over effect.
func (i *momentaryButton) OnMouseMove(x, y int) {
	if x > i.r.Min.X && x < i.r.Max.X && y > i.r.Min.Y && y < i.r.Max.Y {
		i.mouseOver = true
	} else {
		i.mouseOver = false
	}
}

// ToggleButton creates an interactive component which responds to mouse events and toggles state.
func ToggleButton(r image.Rectangle, defaultComponent Component, hover Component, pressed Component, onPress func(), onRelease func()) Button {
	return &toggleButton{r: r, defaultComponent: defaultComponent, hoverComponent: hover, pressedComponent: pressed, onPress: onPress, onRelease: onRelease}
}

// toggleButton is a button which toggles state when pressed.
type toggleButton struct {
	r                image.Rectangle
	defaultComponent Component
	hoverComponent   Component
	pressedComponent Component
	onPress          func()
	onRelease        func()

	mouseOver bool
	pressed   bool
}

// Display renders the button. If the pressed and hover components are defined they will be rendered instead of the default component.
func (i *toggleButton) Display(ctx *DisplayContext) {
	if i.pressed {
		if i.pressedComponent != nil {
			i.pressedComponent.Display(ctx)
			return
		}
	}

	if i.mouseOver {
		if i.hoverComponent != nil {
			i.hoverComponent.Display(ctx)
			return
		}
	}
	i.defaultComponent.Display(ctx)
}

// Update is a no-op.
func (i *toggleButton) Update(ctx *UpdateContext) error {
	return nil
}

// OnMouseEvent calls the onPress and onRelease handlers.
func (i *toggleButton) OnMouseEvent(x, y int, evt MouseEvent) {
	if x > i.r.Min.X && x < i.r.Max.X && y > i.r.Min.Y && y < i.r.Max.Y {
		if evt.EventType == MousePressEvent {

			if i.onPress != nil && i.pressed == false {
				i.onPress()
			} else if i.onRelease != nil && i.pressed {
				i.onRelease()
			}
			i.pressed = !i.pressed
		}
	}
}

// OnMouseMove toggles the mouse over effect.
func (i *toggleButton) OnMouseMove(x, y int) {
	if x > i.r.Min.X && x < i.r.Max.X && y > i.r.Min.Y && y < i.r.Max.Y {
		i.mouseOver = true
	} else {
		i.mouseOver = false
	}
}
