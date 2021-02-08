package ui

import (
	"image"
)

// InteractiveComponent is a test component.
func InteractiveComponent(r image.Rectangle, defaultComponent Component, hover Component, pressed Component, onPressed func()) Component {
	return &interactiveComponent{r: r, defaultComponent: defaultComponent, hoverComponent: hover, pressedComponent: pressed, onPressed: onPressed}
}

type interactiveComponent struct {
	r                image.Rectangle
	defaultComponent Component
	hoverComponent   Component
	pressedComponent Component
	onPressed        func()

	mouseOver bool
	pressed   bool
}

func (i *interactiveComponent) Display(ctx *DisplayContext) {
	if i.pressed {
		i.pressedComponent.Display(ctx)
		return
	}

	if i.mouseOver {
		i.hoverComponent.Display(ctx)
	} else {
		i.defaultComponent.Display(ctx)
	}
}

func (i *interactiveComponent) Update(ctx *UpdateContext) error {
	return nil
}

func (i *interactiveComponent) OnMouseEvent(x, y int, evt MouseEvent) {
	if x > i.r.Min.X && x < i.r.Max.X && y > i.r.Min.Y && y < i.r.Max.Y {
		if evt.EventType == MousePressEvent {
			i.pressed = !i.pressed
			i.onPressed()
		}
	}
}

func (i *interactiveComponent) OnMouseMove(x, y int) {
	if x > i.r.Min.X && x < i.r.Max.X && y > i.r.Min.Y && y < i.r.Max.Y {
		i.mouseOver = true
	} else {
		i.mouseOver = false
	}
}
