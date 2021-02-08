package ui

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// MouseEventType enumerates the type of mouse event.
type MouseEventType int

// String returns the string representation of the event type
func (evt MouseEventType) String() string {
	switch evt {
	case MousePressEvent:
		return "MousePressed"
	case MouseReleaseEvent:
		return "MouseReleased"
	}
	return "Unknown"
}

const (

	// MousePressEvent occurs when a mouse button is pressed
	MousePressEvent MouseEventType = iota

	// MouseReleaseEvent occurs when a mouse button is released.
	MouseReleaseEvent

	// MouseMoveEvent occurs when the mouse moves.
	MouseMoveEvent
)

// MouseEvent stores the mouse butten and event type.
type MouseEvent struct {
	Button    ebiten.MouseButton
	EventType MouseEventType
}

func (evt MouseEvent) String() string {
	return fmt.Sprintf("MouseButtonEvent: Button=%d Event=%s", evt.Button, evt.EventType.String())
}

// MouseButtonHandler is dispatched whenever mouse button events occur.
type MouseButtonHandler interface {
	OnMouseEvent(x, y int, evt MouseEvent)
}

// MouseMoveHandler is dispatched whenever mouse moves.
type MouseMoveHandler interface {
	OnMouseMove(x, y int)
}

// DefaultMouseEventRegistry is the root mouse event registry.
var DefaultMouseEventRegistry = NewMouseEventRegistry(0, 0)

// NewMouseEventRegistry creates a new mouse event registry with the given origin.
func NewMouseEventRegistry(x, y int) *MouseEventRegistry {
	return &MouseEventRegistry{
		make([]MouseButtonHandler, 0),
		make([]MouseMoveHandler, 0),
		image.Pt(x, y),
		image.Pt(0, 0),
	}
}

// MouseEventRegistry stores all the mouse button handlers.
type MouseEventRegistry struct {
	handlers     []MouseButtonHandler
	moveHandlers []MouseMoveHandler
	origin       image.Point

	lastMousePosition image.Point
}

// AddButtonHandler adds a button handler to the registry.
func (r *MouseEventRegistry) AddButtonHandler(h MouseButtonHandler) {
	r.handlers = append(r.handlers, h)
}

// AddMoveHandler adds a move handler to the registry.
func (r *MouseEventRegistry) AddMoveHandler(h MouseMoveHandler) {
	r.moveHandlers = append(r.moveHandlers, h)
}

// Update gets the latest mouse events and dispatches them to the handlers
func (r *MouseEventRegistry) Update() {
	mouseX, mouseY := ebiten.CursorPosition()

	// left
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		r.Dispatch(mouseX, mouseY, MouseEvent{Button: ebiten.MouseButtonLeft, EventType: MousePressEvent})
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		r.Dispatch(mouseX, mouseY, MouseEvent{Button: ebiten.MouseButtonLeft, EventType: MouseReleaseEvent})
	}

	// middle
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		r.Dispatch(mouseX, mouseY, MouseEvent{Button: ebiten.MouseButtonMiddle, EventType: MousePressEvent})
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
		r.Dispatch(mouseX, mouseY, MouseEvent{Button: ebiten.MouseButtonMiddle, EventType: MouseReleaseEvent})
	}

	// right
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		r.Dispatch(mouseX, mouseY, MouseEvent{Button: ebiten.MouseButtonRight, EventType: MousePressEvent})
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		r.Dispatch(mouseX, mouseY, MouseEvent{Button: ebiten.MouseButtonRight, EventType: MouseReleaseEvent})
	}

	// dispatch move handlers
	ex, ey := mouseX-r.origin.X, mouseY-r.origin.Y
	for i := 0; i < len(r.moveHandlers); i++ {
		if mouseX != r.lastMousePosition.X || mouseY != r.lastMousePosition.Y {
			r.moveHandlers[i].OnMouseMove(ex, ey)
			// fmt.Printf("calling mouse move %d %d\n", ex, ey)
		}
	}
	r.lastMousePosition = image.Pt(mouseX, mouseY)
}

// Dispatch emits a handler event to the handlers when fired.
func (r *MouseEventRegistry) Dispatch(x, y int, evt MouseEvent) {
	ex, ey := x-r.origin.X, y-r.origin.Y
	for i := 0; i < len(r.handlers); i++ {
		r.handlers[i].OnMouseEvent(ex, ey, evt)
	}
}

// MouseHandlerFunc creates a MouseButtonHandler from a function.
func MouseHandlerFunc(h func(x, y int, evt MouseEvent)) MouseButtonHandler {
	return &simpleMouseHandler{buttonHandler: h}
}

// MouseMoveHandlerFunc creates a MouseMoveHandler from a function.
func MouseMoveHandlerFunc(h func(x, y int)) MouseMoveHandler {
	return &simpleMouseHandler{moveHandler: h}
}

type simpleMouseHandler struct {
	buttonHandler func(x, y int, evt MouseEvent)
	moveHandler   func(x, y int)
}

func (s *simpleMouseHandler) OnMouseEvent(x, y int, evt MouseEvent) {
	if s.buttonHandler != nil {
		s.buttonHandler(x, y, evt)
	}
}

func (s *simpleMouseHandler) OnMouseMove(x, y int) {
	if s.moveHandler != nil {
		s.moveHandler(x, y)
	}
}
