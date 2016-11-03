// Copyright (c) 2016, redstarcoder
package widgets

import (
	"image/color"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/redstarcoder/draw2dglkit"
	"github.com/redstarcoder/draw2dui"
)

type Button struct {
	x, y, width, height        float64
	enabled, redraw, hasCursor bool
	shape                      *draw2d.Path
	window, offscreen          *glfw.Window
	gc                         *draw2d.GraphicContext // gc may be overwritten so it's a pointer to an interface
	name, text                 string
}

func NewButton(gc *draw2d.GraphicContext, window, offscreen *glfw.Window, x, y float64, text string) *Button {
	Button := &Button{
		gc:        gc,
		window:    window,
		offscreen: offscreen,
		x:         x,
		y:         y,
		height:    (*gc).GetFontSize() + 13,
		enabled:   true,
		shape:     &draw2d.Path{},
		redraw:    true,
		name:      draw2dui.NameWidget("Button"),
		text:      text,
	}
	Button.reshape()
	return Button
}

// reshape recreates btn's path, which is used for drawing it to the screen
func (btn *Button) reshape() {
	// Recalulate width
	_, _, btn.width, _ = (*btn.gc).GetStringBounds(btn.text)
	btn.width += 6

	draw2dkit.Rectangle(btn.shape, btn.x, btn.y, btn.x+btn.width-1, btn.y+btn.height-1)
	btn.redraw = true
}

// Name returns btn's name
func (btn *Button) Name() string {
	return btn.name
}

// Draw draws the widget, selected determines if the widget displays as selected or not, and forceRedraw
// forces a full redraw of the widget.
// TODO draw dotted border if selected
func (btn *Button) Draw(selected, forceRedraw bool) {
	if btn.redraw || forceRedraw {
		gc := *btn.gc
		gc.Save()
		btn.clear(gc, false)
		gl.LineWidth(1)
		var fg, bg color.RGBA
		if btn.hasCursor {
			fg = color.RGBA{255, 255, 255, 0xff}
			bg = color.RGBA{0, 0, 0, 0xff}
		} else {
			fg = color.RGBA{0, 0, 0, 0xff}
			bg = color.RGBA{255, 255, 255, 0xff}
		}
		gc.SetFillColor(bg)
		gc.SetStrokeColor(fg)
		gc.FillStroke(btn.shape)
		gc.SetFillColor(fg)
		gc.FillStringAt(btn.text, btn.x+3, btn.y+6+gc.GetFontSize())
		gc.Restore()

		btn.redraw = false
	}
}

// clear clears btn's border to make redrawing look more consistent without a complete screen clear. If fill
// is true, it also fills the shape with white.
func (btn *Button) clear(gc draw2d.GraphicContext, fill bool) {
	gl.LineWidth(3)
	if !fill {
		gc.SetStrokeColor(color.RGBA{255, 255, 255, 0xff})
		gc.Stroke(btn.shape)
	} else {
		gc.Save()
		gc.SetStrokeColor(color.RGBA{255, 255, 255, 0xff})
		gc.SetFillColor(color.RGBA{255, 255, 255, 0xff})
		gc.FillStroke(btn.shape)
		gc.Restore()
	}
}

// Handle returns false
func (btn *Button) Handle(selected bool) bool {
	return false
}

// KeyPress has the widget process a KeyPress event
func (btn *Button) KeyPress(key glfw.Key, action glfw.Action, mods glfw.ModifierKey) draw2dui.Event {
	if action == glfw.Release {
		return draw2dui.EventNone
	}
	if key == glfw.KeyEnter {
		// TODO make the button look like it got pressed
		return draw2dui.EventConfirm
	}
	return draw2dui.EventNone
}

// CharPress returns draw2dui.EventNone
func (btn *Button) CharPress(char rune) draw2dui.Event {
	return draw2dui.EventNone
}

// MMove has the widget process a MouseMove event
func (btn *Button) MMove(xpos, ypos float64) draw2dui.Event {
	inside := btn.IsInside(xpos, ypos)
	if btn.hasCursor && !inside {
		btn.hasCursor = false
		btn.redraw = true
		return draw2dui.EventAction
	} else if !inside {
		return draw2dui.EventNone
	}
	if !btn.hasCursor {
		btn.window.SetCursor(glfw.CreateStandardCursor(int(glfw.HandCursor)))
		btn.hasCursor = true
		btn.redraw = true
		return draw2dui.EventHasCursor
	} else {
		return draw2dui.EventHasCursor
	}
	return draw2dui.EventNone
}

// MClick has the widget process a MouseClick event
func (btn *Button) MClick(xpos, ypos float64, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) draw2dui.Event {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		btn.redraw = true
		if !btn.IsInside(xpos, ypos) {
			return draw2dui.EventNone
		}
	} else {
		return draw2dui.EventNone
	}
	// TODO make the button look like it got pressed
	return draw2dui.EventConfirm
}

// SetPos changes the widget's x, y coordinates
func (btn *Button) SetPos(x, y float64) {
	btn.clear(*btn.gc, true)
	btn.x, btn.y = x, y
	btn.reshape()
}

// GetPos retrieves the widget's x, y coordinates
func (btn *Button) GetPos() (float64, float64) {
	return btn.x, btn.y
}

// SetDimensions sets btn's drawn width and height
func (btn *Button) SetDimensions(w, h float64) {
	btn.clear(*btn.gc, true)
	btn.width, btn.height = w, h
	btn.reshape() // reshape overwrites width
}

// GetDimensions returns btn's drawn width and height
func (btn *Button) GetDimensions() (float64, float64) {
	return btn.width, btn.height
}

// IsInside checks if point x, y is inside of the widget's boundaries. It uses btn.offscreen as a pallet
func (btn *Button) IsInside(x, y float64) bool {
	return draw2dglkit.IsPointInShape(*btn.gc, btn.offscreen, x, y, btn.shape)
}

// SetString sets btn's text, using btn.maxlen as the max length
func (btn *Button) SetString(s string) {
	btn.text = s
	btn.redraw = true
}

// GetString returns btn's text
func (btn *Button) GetString() string {
	return btn.text
}

// SetInt does nothing
func (btn *Button) SetInt(i int) {
}

// GetInt returns -1
func (btn *Button) GetInt() int {
	return -1
}

// SetData does nothing
func (btn *Button) SetData(d interface{}) {
}

// GetData returns nil
func (btn *Button) GetData() interface{} {
	return nil
}

// SetEnabled enables or disables the widget
func (btn *Button) SetEnabled(enabled bool) {
	if btn.enabled != enabled {
		btn.enabled = enabled
		btn.redraw = true
	}
}

// GetEnabled returns whether the widget is enabled or not
func (btn *Button) GetEnabled() bool {
	return btn.enabled
}
