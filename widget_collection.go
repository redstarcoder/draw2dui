// Copyright (c) 2016, redstarcoder
package draw2dui

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
)

// TODO Add a method for removing widgets?

// WidgetCollection is a struct for managing many widgets at once. It has many helper methods for handling
// mouse and keyboard events.
type WidgetCollection struct {
	gc      *draw2d.GraphicContext
	window  *glfw.Window
	widgets map[string]Widget

	mx, my      float64
	hasCursor   bool
	forceRedraw bool
	selected    string
}

// NewWidgetCollection creates a new widget collection and registers all the widgets
func NewWidgetCollection(gc *draw2d.GraphicContext, window *glfw.Window, widgets ...Widget) *WidgetCollection {
	wc := &WidgetCollection{
		gc:      gc,
		window:  window,
		widgets: make(map[string]Widget, len(widgets)),
	}
	for _, w := range widgets {
		wc.Register(w)
	}
	return wc
}

// Register adds a widget to the collection
func (wc *WidgetCollection) Register(widget Widget) {
	if len(wc.selected) == 0 {
		wc.selected = widget.Name()
	}
	wc.widgets[widget.Name()] = widget
}

// Draw draws all the widgets to the screen. If force is true, it will force a redraw of all widgets in the
// collection.
func (wc *WidgetCollection) Draw() {
	for _, w := range wc.widgets {
		w.Draw(w.Name() == wc.selected, wc.forceRedraw)
	}
	wc.forceRedraw = false
}

// Handle processes all the idle events for every widget in the collection. Returns whether it requests a
// call to WidgetCollection.Draw or not.
func (wc *WidgetCollection) Handle() (redraw bool) {
	for _, w := range wc.widgets {
		if w.Handle(w.Name() == wc.selected) {
			redraw = true
		}
	}
	return
}

// KeyPress has the selected widget process a KeyPress event, returning the selected widget and the event if
// it isn't EventNone.
func (wc *WidgetCollection) KeyPress(key glfw.Key, action glfw.Action, mods glfw.ModifierKey) (Widget, Event) {
	if len(wc.selected) == 0 {
		return nil, EventNone
	}
	if w := wc.widgets[wc.selected]; w.KeyPress(key, action, mods) == EventAction {
		return w, EventAction
	}
	return nil, EventNone
}

// CharPress has the selected widget process a character, returning the selected widget and the event if it
// isn't EventNone.
func (wc *WidgetCollection) CharPress(char rune) (Widget, Event) {
	if len(wc.selected) == 0 {
		return nil, EventNone
	}
	if w := wc.widgets[wc.selected]; w.CharPress(char) == EventAction {
		return w, EventAction
	}
	return nil, EventNone
}

// MMove has all the widgets in the collection process a MouseMove event, returning the a widget and event
// if the cursor changes. Always returns the moused-over widget, unless there isn't one, then it returns a
// widget that returned EventAction, if any.
func (wc *WidgetCollection) MMove(xpos, ypos float64) (widget Widget, event Event) {
	wc.mx, wc.my = xpos, ypos
	hasCursor := true
	for _, w := range wc.widgets {
		if ev := w.MMove(xpos, ypos); ev == EventHasCursor {
			if wc.hasCursor {
				wc.hasCursor = false
				event = EventHasCursor
			}
			hasCursor = false
			widget = w
		} else if widget == nil && ev == EventAction {
			widget = w
			event = EventAction
		}
	}
	if hasCursor && !wc.hasCursor {
		wc.window.SetCursor(glfw.CreateStandardCursor(int(glfw.ArrowCursor)))
		wc.hasCursor = true
	}
	return
}

// MClick has the all widgets in the collection process a MouseClick event, returning the a widget and event
// if it isn't EventNone.
func (wc *WidgetCollection) MClick(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) (Widget, Event) {
	for _, w := range wc.widgets {
		if w.MClick(wc.mx, wc.my, button, action, mods) == EventSelected {
			if w.Name() != wc.selected {
				wc.selected = w.Name()
				wc.forceRedraw = true
			}
			return w, EventSelected
		}
	}
	if action == glfw.Press {
		if len(wc.selected) > 0 {
			wc.selected = ""
			wc.forceRedraw = true
		}
		return nil, EventSelected
	}
	return nil, EventNone
}

// Reshape should be called whenever the draw2d.GraphicContext is resized
// TODO handle w & h
func (wc *WidgetCollection) Reshape(w, h int) {
	wc.forceRedraw = true
}

// Refresh should be called when a window refresh occurs (See: glfw.RefreshCallback)
func (wc *WidgetCollection) Refresh() {
	wc.forceRedraw = true
}
