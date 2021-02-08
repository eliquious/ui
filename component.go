package ui

// Component is a display component for the screen.
type Component interface {
	Update(ctx *UpdateContext) error
	Display(ctx *DisplayContext)
}

// RenderFunc is a render function for a simple component.
type RenderFunc func(ctx *DisplayContext)

// UpdateFunc is an update function for a simple component.
type UpdateFunc func(ctx *UpdateContext) error

// SimpleComponent creates a simple component from a render function.
func SimpleComponent(fn RenderFunc) Component {
	return &simpleComponent{fn, nil}
}

// ComponentFunc creates a simple component from render and update functions.
func ComponentFunc(update UpdateFunc, render RenderFunc) Component {
	return &simpleComponent{render, update}
}

type simpleComponent struct {
	renderFunc RenderFunc
	updateFunc UpdateFunc
}

func (s *simpleComponent) Update(ctx *UpdateContext) error {
	if s.updateFunc != nil {
		return s.updateFunc(ctx)
	}
	return nil
}

func (s *simpleComponent) Display(ctx *DisplayContext) {
	s.renderFunc(ctx)
}

// StackedComponent creates a component which updates and renders all the child components.
func StackedComponent(c ...Component) Component {
	return &stackedComponent{c}
}

type stackedComponent struct {
	children []Component
}

func (s *stackedComponent) Update(ctx *UpdateContext) error {
	for i := 0; i < len(s.children); i++ {
		if err := s.children[i].Update(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *stackedComponent) Display(ctx *DisplayContext) {
	for i := 0; i < len(s.children); i++ {
		s.children[i].Display(ctx)
	}
}

// UpdateHandler is an interface for non-display components.
type UpdateHandler interface {
	Update(ctx *UpdateContext) error
}

// UpdateHandlerFunc creates an update handler
func UpdateHandlerFunc(f UpdateFunc) UpdateHandler {
	return &updateHandler{f}
}

type updateHandler struct {
	update UpdateFunc
}

func (u *updateHandler) Update(ctx *UpdateContext) error {
	return u.update(ctx)
}
