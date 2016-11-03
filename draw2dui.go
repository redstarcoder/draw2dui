// Copyright (c) 2016, redstarcoder
// Package draw2dui offers useful tools for drawing and handling UIs in Golang using draw2d with OpenGL.
package draw2dui

import (
	"fmt"
	"github.com/go-gl/glfw/v3.1/glfw"
	"sync/atomic"
)

var widgetCount int32

// Widget is an interface for draw2dui widgets
type Widget interface {
	// Name returns the widget's name
	Name() string
	// Draw draws the widget, selected determines if the widget displays as selected or not, and forceRedraw
	// forces a full redraw of the widget.
	Draw(selected, forceRedraw bool)
	// Handle processes the widget's idle actions, selected determines if the widget behaves as selected or
	// not. Returns if it needs a draw or not.
	Handle(selected bool) bool
	// KeyPress has the widget process a KeyPress event
	KeyPress(key glfw.Key, action glfw.Action, mods glfw.ModifierKey) Event
	// CharPress has the widget process a character
	CharPress(char rune) Event
	// MMove has the widget process a MouseMove event
	MMove(xpos, ypos float64) Event
	// MClick has the widget process a MouseClick event
	MClick(xpos, ypos float64, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) Event
	// SetPos changes the widget's x, y coordinates
	SetPos(x, y float64)
	// GetPos retrieves the widget's x, y coordinates
	GetPos() (float64, float64)
	// SetDimensions changes the widget's width and height
	SetDimensions(w, h float64)
	// GetDimensions returns the widget's width and height
	GetDimensions() (float64, float64)
	// IsInside checks if point x, y is inside of the widget's boundaries
	IsInside(x, y float64) bool
	// SetString sets the widget's string, if it supports it
	SetString(s string)
	// GetString returns a string, if the widget supports it
	GetString() string
	// SetInt sets the widget's int, if it supports it
	SetInt(i int)
	// GetInt returns an int, if the widget supports it
	GetInt() int
	// SetData sets the widget's data if it supports it. This must be a type supported by the widget.
	SetData(d interface{})
	// GetData returns an interface{}, if the widget supports it
	GetData() interface{}
	// SetEnabled enables or disables the widget
	SetEnabled(enabled bool)
	// GetEnabled returns whether the widget is enabled or not
	GetEnabled() bool
}

// NameWidget returns a unique widget name. It is thread-safe.
func NameWidget(w string) string {
	return fmt.Sprintf("%s-%d", w, atomic.AddInt32(&widgetCount, 1))
}

// Event is for widget events
type Event int

const (
	// EventNone means the widget returned normally
	EventNone Event = iota
	// EventNext means the next widget should be selected
	EventNext
	// EventPrevious means the previous widget should be sekected
	EventPrevious
	// EventExit means the widget or the application should be closed
	EventExit
	// EventConfirm means the user has confirmed an action (typically they hit enter)
	EventConfirm
	// EventAction means the user modified the widget in some way (changed a dropdown, changed text)
	EventAction
	// EventSelect means the widget was selected
	EventSelected
	// EventHasCursor means the widget currently controls the cursor
	EventHasCursor
)
