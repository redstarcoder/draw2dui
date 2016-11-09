// Copyright (c) 2016, redstarcoder
package widgets

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/redstarcoder/draw2dglkit"
	"github.com/redstarcoder/draw2dui"
	"image/color"
	"strings"
	"unicode/utf8"
)

// TODO text highlighting
type TextBox struct {
	cursor                     *Cursor
	x, y, width, height        float64
	maxlen                     int
	enabled, redraw, hasCursor bool
	shape                      *draw2d.Path
	window, offscreen          *glfw.Window
	gc                         *draw2d.GraphicContext // gc may be overwritten so it's a pointer to an interface
	name                       string
}

// NewTextBox creates a new TextBox widget
// BUG(x) TextBox does not support enabled state
// BUG(x) TextBox default is disabled, should be enabled
// BUG(x) TextBox should have a scrollbar
// BUG(x) TextBox text wrapping should be optional
func NewTextBox(gc *draw2d.GraphicContext, window, offscreen *glfw.Window, x, y, width, height float64, text string) *TextBox {
	textBox := &TextBox{
		cursor:    &Cursor{text: text},
		gc:        gc,
		window:    window,
		offscreen: offscreen,
		x:         x,
		y:         y,
		width:     width,
		height:    height,
		maxlen:    0x7ffffffe,
		//		enabled:   true,
		shape:  &draw2d.Path{},
		redraw: true,
		name:   draw2dui.NameWidget("TextBox"),
	}
	textBox.cursor.GenLines(*gc, width)
	textBox.reshape()
	return textBox
}

func (tb *TextBox) InsertLine(s string) {
	tb.cursor.InsertLine(s)
	tb.cursor.GenLines(*tb.gc, tb.width)
}

// reshape recreates tf's path, which is used for drawing it to the screen
func (tb *TextBox) reshape() {
	tb.shape = &draw2d.Path{}
	tb.cursor.maxLines = int((tb.height - 2) / ((*tb.gc).GetFontSize() + 3))
	draw2dkit.Rectangle(tb.shape, tb.x, tb.y, tb.x+tb.width-1, tb.y+tb.height-1)
	tb.redraw = true
}

// Name returns tf's name
func (tb *TextBox) Name() string {
	return tb.name
}

// Draw draws the widget, selected determines if the widget displays as selected or not, and forceRedraw
// forces a full redraw of the widget.
func (tb *TextBox) Draw(selected, forceRedraw bool) {
	if tb.redraw || forceRedraw {
		gc := *tb.gc
		gc.Save()
		tb.clear(gc, false)
		gl.LineWidth(1)
		gc.SetFillColor(color.RGBA{255, 255, 255, 0xff})
		gc.SetStrokeColor(color.RGBA{0, 0, 0, 0xff})
		gc.FillStroke(tb.shape)
		gc.SetFillColor(color.RGBA{0, 0, 0, 0xff})
		y := tb.y + float64(tb.cursor.maxLines)*(gc.GetFontSize()+3)
		for i := 0; i < tb.cursor.maxLines && i+tb.cursor.iY < len(tb.cursor.textLines); i++ {
			gc.FillStringAt(tb.cursor.textLines[len(tb.cursor.textLines)-1-i-tb.cursor.iY], tb.x+1, y)
			y -= gc.GetFontSize() + 3
		}
		gc.Restore()

		tb.redraw = false
	}
}

// clear clears tf's border to make redrawing look more consistent without a complete screen clear. If fill
// is true, it also fills the shape with white.
func (tb *TextBox) clear(gc draw2d.GraphicContext, fill bool) {
	gl.LineWidth(3)
	if !fill {
		gc.SetStrokeColor(color.RGBA{255, 255, 255, 0xff})
		gc.Stroke(tb.shape)
	} else {
		gc.Save()
		gc.SetStrokeColor(color.RGBA{255, 255, 255, 0xff})
		gc.SetFillColor(color.RGBA{255, 255, 255, 0xff})
		gc.FillStroke(tb.shape)
		gc.Restore()
	}
}

// Handle processes tf's cursor
func (tb *TextBox) Handle(selected bool) bool {
	if selected {
		if tb.cursor.Blink() {
			tb.redraw = true
			return true
		}
	}
	return false
}

// KeyPress has the widget process a KeyPress event
func (tb *TextBox) KeyPress(key glfw.Key, action glfw.Action, mods glfw.ModifierKey) draw2dui.Event {
	if action == glfw.Release {
		return draw2dui.EventNone
	}
	switch key {
	default:
		return draw2dui.EventNone
	/*case glfw.KeyLeft:
		if tb.cursor.MoveLeft() {
			tb.redraw = true
			return draw2dui.EventAction
		}
	case glfw.KeyRight:
		if tb.cursor.MoveRight() {
			tb.redraw = true
			return draw2dui.EventAction
		}
	case glfw.KeyBackspace:
		if tb.cursor.Backspace() {
			tb.redraw = true
			return draw2dui.EventAction
		}
	case glfw.KeyEnter:
		return draw2dui.EventConfirm*/
	case glfw.KeyUp:
		if tb.cursor.iY < len(tb.cursor.textLines)-1 {
			tb.cursor.iY++
			tb.redraw = true
			return draw2dui.EventAction
		}
	case glfw.KeyDown:
		if tb.cursor.iY > 0 {
			tb.cursor.iY--
			tb.redraw = true
			return draw2dui.EventAction
		}
	}
	return draw2dui.EventNone
}

// CharPress adds a character to the TextBox
func (tb *TextBox) CharPress(char rune) draw2dui.Event {
	if !utf8.ValidRune(char) || len(tb.GetString()) >= tb.maxlen {
		return draw2dui.EventNone
	}
	tb.cursor.Insert(string(char))
	tb.redraw = true
	return draw2dui.EventAction
}

// MMove has the widget process a MouseMove event
func (tb *TextBox) MMove(xpos, ypos float64) draw2dui.Event {
	if !tb.IsInside(xpos, ypos) {
		tb.hasCursor = false
		return draw2dui.EventNone
	}
	if !tb.hasCursor {
		tb.window.SetCursor(glfw.CreateStandardCursor(int(glfw.IBeamCursor)))
		tb.hasCursor = true
	}
	return draw2dui.EventHasCursor
}

// MClick has the widget process a MouseClick event
func (tb *TextBox) MClick(xpos, ypos float64, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) draw2dui.Event {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		tb.redraw = true
		if !tb.IsInside(xpos, ypos) {
			return draw2dui.EventNone
		}
	} else {
		return draw2dui.EventNone
	}
	tb.cursor.drawCursor = false // gets swapped to true before next draw
	tb.cursor.MoveToX(*tb.gc, tb.x, xpos, tb.width)
	return draw2dui.EventSelected
}

// SetPos changes the widget's x, y coordinates
func (tb *TextBox) SetPos(x, y float64) {
	tb.clear(*tb.gc, true)
	tb.x, tb.y = x, y
	tb.reshape()
}

// GetPos retrieves the widget's x, y coordinates
func (tb *TextBox) GetPos() (float64, float64) {
	return tb.x, tb.y
}

// SetDimensions sets tf's drawn width and height
func (tb *TextBox) SetDimensions(w, h float64) {
	tb.clear(*tb.gc, true)
	tb.width, tb.height = w, h
	tb.reshape()
	tb.SetString(strings.Join(tb.cursor.textLines, "\n"))
}

// GetDimensions returns tf's drawn width and height
func (tb *TextBox) GetDimensions() (float64, float64) {
	return tb.width, tb.height
}

// IsInside checks if point x, y is inside of the widget's boundaries. It uses tb.offscreen as a pallet
func (tb *TextBox) IsInside(x, y float64) bool {
	return draw2dglkit.IsPointInShape(*tb.gc, tb.offscreen, x, y, tb.shape)
}

// SetString sets tf's text, using tb.maxlen as the max length
func (tb *TextBox) SetString(s string) {
	if len(s) > tb.maxlen {
		tb.cursor.text = s[:tb.maxlen]
	} else {
		tb.cursor.text = s
	}
	tb.cursor.GenLines(*tb.gc, tb.width)
	//if tb.GetInt() > len(tb.GetString()) {
	//	tb.SetInt(len(tb.GetString()))
	//}
	tb.redraw = true
}

// GetString returns tf's text
// BUG(x) needs to rebuild cursor.text from cursor.textLines
func (tb *TextBox) GetString() string {
	tb.cursor.text = strings.Join(tb.cursor.textLines, "\n")
	return tb.cursor.text
}

// SetInt moves tf's text cursor
func (tb *TextBox) SetInt(i int) {
	if tb.cursor.MoveTo(i) {
		tb.redraw = true
	}
}

// GetInt returns tf's text cursor's position
func (tb *TextBox) GetInt() int {
	return tb.cursor.i
}

// SetData does nothing
func (tb *TextBox) SetData(d interface{}) {
}

// GetData returns nil
func (tb *TextBox) GetData() interface{} {
	return nil
}

// SetEnabled enables or disables the widget
func (tb *TextBox) SetEnabled(enabled bool) {
	panic("not implemented.")
	if tb.enabled != enabled {
		tb.enabled = enabled
		tb.redraw = true
	}
}

// GetEnabled returns whether the widget is enabled or not
func (tb *TextBox) GetEnabled() bool {
	return tb.enabled
}

// TODO text highlighting
type TextField struct {
	cursor                     *Cursor
	x, y, width, height        float64
	maxlen                     int
	enabled, redraw, hasCursor bool
	shape                      *draw2d.Path
	window, offscreen          *glfw.Window
	gc                         *draw2d.GraphicContext // gc may be overwritten so it's a pointer to an interface
	name                       string
}

// NewTextField creates a new TextField widget
func NewTextField(gc *draw2d.GraphicContext, window, offscreen *glfw.Window, x, y, width float64, text string, maxlen int) *TextField {
	textField := &TextField{
		cursor:    &Cursor{text: text},
		gc:        gc,
		window:    window,
		offscreen: offscreen,
		x:         x,
		y:         y,
		width:     width,
		height:    (*gc).GetFontSize() + 7,
		maxlen:    maxlen,
		enabled:   true,
		shape:     &draw2d.Path{},
		redraw:    true,
		name:      draw2dui.NameWidget("TextField"),
	}
	textField.reshape()
	return textField
}

// reshape recreates tf's path, which is used for drawing it to the screen
func (tf *TextField) reshape() {
	tf.shape = &draw2d.Path{}
	draw2dkit.Rectangle(tf.shape, tf.x, tf.y, tf.x+tf.width-1, tf.y+tf.height-1)
	tf.redraw = true
}

// Name returns tf's name
func (tf *TextField) Name() string {
	return tf.name
}

// Draw draws the widget, selected determines if the widget displays as selected or not, and forceRedraw
// forces a full redraw of the widget.
func (tf *TextField) Draw(selected, forceRedraw bool) {
	if tf.redraw || forceRedraw {
		gc := *tf.gc
		gc.Save()
		tf.clear(gc, false)
		gl.LineWidth(1)
		gc.SetFillColor(color.RGBA{255, 255, 255, 0xff})
		gc.SetStrokeColor(color.RGBA{0, 0, 0, 0xff})
		gc.FillStroke(tf.shape)
		gc.SetFillColor(color.RGBA{0, 0, 0, 0xff})
		if selected {
			fillStringAtWidthCursor(*tf.gc, tf.cursor, tf.x+1, tf.y+3+gc.GetFontSize(), tf.width-2)
		} else {
			fillStringAtWidth(*tf.gc, tf.cursor.text[tf.cursor.iOffset:], tf.x+1, tf.y+3+gc.GetFontSize(), tf.width-2)
		}
		gc.Restore()

		tf.redraw = false
	}
}

// clear clears tf's border to make redrawing look more consistent without a complete screen clear. If fill
// is true, it also fills the shape with white.
func (tf *TextField) clear(gc draw2d.GraphicContext, fill bool) {
	gl.LineWidth(3)
	if !fill {
		gc.SetStrokeColor(color.RGBA{255, 255, 255, 0xff})
		gc.Stroke(tf.shape)
	} else {
		gc.Save()
		gc.SetStrokeColor(color.RGBA{255, 255, 255, 0xff})
		gc.SetFillColor(color.RGBA{255, 255, 255, 0xff})
		gc.FillStroke(tf.shape)
		gc.Restore()
	}
}

// Handle processes tf's cursor
func (tf *TextField) Handle(selected bool) bool {
	if selected {
		if tf.cursor.Blink() {
			tf.redraw = true
			return true
		}
	}
	return false
}

// KeyPress has the widget process a KeyPress event
func (tf *TextField) KeyPress(key glfw.Key, action glfw.Action, mods glfw.ModifierKey) draw2dui.Event {
	if action == glfw.Release {
		return draw2dui.EventNone
	}
	switch key {
	default:
		return draw2dui.EventNone
	case glfw.KeyLeft:
		if tf.cursor.MoveLeft() {
			tf.redraw = true
			return draw2dui.EventAction
		}
	case glfw.KeyRight:
		if tf.cursor.MoveRight() {
			tf.redraw = true
			return draw2dui.EventAction
		}
	case glfw.KeyBackspace:
		if tf.cursor.Backspace() {
			tf.redraw = true
			return draw2dui.EventAction
		}
	case glfw.KeyEnter:
		return draw2dui.EventConfirm
	}
	return draw2dui.EventNone
}

// CharPress adds a character to the textfield
func (tf *TextField) CharPress(char rune) draw2dui.Event {
	if !utf8.ValidRune(char) || len(tf.GetString()) >= tf.maxlen {
		return draw2dui.EventNone
	}
	tf.cursor.Insert(string(char))
	tf.redraw = true
	return draw2dui.EventAction
}

// MMove has the widget process a MouseMove event
func (tf *TextField) MMove(xpos, ypos float64) draw2dui.Event {
	if !tf.IsInside(xpos, ypos) {
		tf.hasCursor = false
		return draw2dui.EventNone
	}
	if !tf.hasCursor {
		tf.window.SetCursor(glfw.CreateStandardCursor(int(glfw.IBeamCursor)))
		tf.hasCursor = true
	}
	return draw2dui.EventHasCursor
}

// MClick has the widget process a MouseClick event
func (tf *TextField) MClick(xpos, ypos float64, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) draw2dui.Event {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		tf.redraw = true
		if !tf.IsInside(xpos, ypos) {
			return draw2dui.EventNone
		}
	} else {
		return draw2dui.EventNone
	}
	tf.cursor.drawCursor = false // gets swapped to true before next draw
	tf.cursor.MoveToX(*tf.gc, tf.x, xpos, tf.width)
	return draw2dui.EventSelected
}

// SetPos changes the widget's x, y coordinates
func (tf *TextField) SetPos(x, y float64) {
	tf.clear(*tf.gc, true)
	tf.x, tf.y = x, y
	tf.reshape()
}

// GetPos retrieves the widget's x, y coordinates
func (tf *TextField) GetPos() (float64, float64) {
	return tf.x, tf.y
}

// SetDimensions sets tf's drawn width and height
func (tf *TextField) SetDimensions(w, h float64) {
	tf.clear(*tf.gc, true)
	tf.width, tf.height = w, h
	tf.reshape()
}

// GetDimensions returns tf's drawn width and height
func (tf *TextField) GetDimensions() (float64, float64) {
	return tf.width, tf.height
}

// IsInside checks if point x, y is inside of the widget's boundaries. It uses tf.offscreen as a pallet
func (tf *TextField) IsInside(x, y float64) bool {
	return draw2dglkit.IsPointInShape(*tf.gc, tf.offscreen, x, y, tf.shape)
}

// SetString sets tf's text, using tf.maxlen as the max length
func (tf *TextField) SetString(s string) {
	if len(s) > tf.maxlen {
		tf.cursor.text = s[:tf.maxlen]
	} else {
		tf.cursor.text = s
	}
	if tf.GetInt() > len(tf.GetString()) {
		// BUG(x) shortening strings can end up with bad iOffset / iEdge
		tf.SetInt(len(tf.GetString()))
	}
	tf.redraw = true
}

// GetString returns tf's text
func (tf *TextField) GetString() string {
	return tf.cursor.text
}

// SetInt moves tf's text cursor
func (tf *TextField) SetInt(i int) {
	if tf.cursor.MoveTo(i) {
		tf.redraw = true
	}
}

// GetInt returns tf's text cursor's position
func (tf *TextField) GetInt() int {
	return tf.cursor.i
}

// SetData does nothing
func (tf *TextField) SetData(d interface{}) {
}

// GetData returns nil
func (tf *TextField) GetData() interface{} {
	return nil
}

// SetEnabled enables or disables the widget
func (tf *TextField) SetEnabled(enabled bool) {
	if tf.enabled != enabled {
		tf.enabled = enabled
		tf.redraw = true
	}
}

// GetEnabled returns whether the widget is enabled or not
func (tf *TextField) GetEnabled() bool {
	return tf.enabled
}

type Label struct {
	x, y, width, height float64
	redraw              bool
	shape               *draw2d.Path
	window, offscreen   *glfw.Window
	gc                  *draw2d.GraphicContext // gc may be overwritten so it's a pointer to an interface
	name, text          string
}

func NewLabel(gc *draw2d.GraphicContext, window, offscreen *glfw.Window, x, y float64, text string) *Label {
	Label := &Label{
		gc:        gc,
		window:    window,
		offscreen: offscreen,
		x:         x,
		y:         y,
		height:    (*gc).GetFontSize() + 6,
		shape:     &draw2d.Path{},
		redraw:    true,
		name:      draw2dui.NameWidget("Label"),
		text:      text,
	}
	Label.reshape()
	return Label
}

// reshape recreates lbl's path, which is used for drawing it to the screen
func (lbl *Label) reshape() {
	lbl.shape = &draw2d.Path{}
	var x float64
	x, _, lbl.width, _ = (*lbl.gc).GetStringBounds(lbl.text)
	lbl.width += x + 5
	draw2dkit.Rectangle(lbl.shape, lbl.x, lbl.y, lbl.x+lbl.width-1, lbl.y+lbl.height-1)
	lbl.redraw = true
}

// Name returns lbl's name
func (lbl *Label) Name() string {
	return lbl.name
}

// Draw draws the widget, selected determines if the widget displays as selected or not, and forceRedraw
// forces a full redraw of the widget.
func (lbl *Label) Draw(selected, forceRedraw bool) {
	if lbl.redraw || forceRedraw {
		gc := *lbl.gc
		gc.Save()
		gc.BeginPath()
		gl.LineWidth(1)
		fg := color.RGBA{0, 0, 0, 0xff}
		bg := color.RGBA{255, 255, 255, 0xff}
		gc.SetFillColor(bg)
		gc.Fill(lbl.shape)
		gc.SetFillColor(fg)
		gc.FillStringAt(lbl.text, lbl.x+1, lbl.y+1+gc.GetFontSize())
		gc.Restore()

		lbl.redraw = false
	}
}

// clear fills the shape with white
func (lbl *Label) clear(gc draw2d.GraphicContext) {
	gc.Save()
	gc.BeginPath()
	gc.SetFillColor(color.RGBA{255, 255, 255, 0xff})
	gc.Fill(lbl.shape)
	gc.Restore()
}

// Handle returns false
func (lbl *Label) Handle(selected bool) bool {
	return false
}

// KeyPress returns draw2dui.EventNone
func (lbl *Label) KeyPress(key glfw.Key, action glfw.Action, mods glfw.ModifierKey) draw2dui.Event {
	return draw2dui.EventNone
}

// CharPress returns draw2dui.EventNone
func (lbl *Label) CharPress(char rune) draw2dui.Event {
	return draw2dui.EventNone
}

// MMove returns draw2dui.EventNone
func (lbl *Label) MMove(xpos, ypos float64) draw2dui.Event {
	return draw2dui.EventNone
}

// MClick returns draw2dui.EventNone
func (lbl *Label) MClick(xpos, ypos float64, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) draw2dui.Event {
	return draw2dui.EventNone
}

// SetPos changes the widget's x, y coordinates
func (lbl *Label) SetPos(x, y float64) {
	lbl.clear(*lbl.gc)
	lbl.x, lbl.y = x, y
	lbl.reshape()
}

// GetPos retrieves the widget's x, y coordinates
func (lbl *Label) GetPos() (float64, float64) {
	return lbl.x, lbl.y
}

// SetDimensions sets lbl's drawn width and height
func (lbl *Label) SetDimensions(w, h float64) {
	lbl.clear(*lbl.gc)
	lbl.width, lbl.height = w, h
	lbl.reshape() // reshape overwrites width
}

// GetDimensions returns lbl's drawn width and height
func (lbl *Label) GetDimensions() (float64, float64) {
	return lbl.width, lbl.height
}

// IsInside checks if point x, y is inside of the widget's boundaries. It uses lbl.offscreen as a pallet
func (lbl *Label) IsInside(x, y float64) bool {
	return draw2dglkit.IsPointInShape(*lbl.gc, lbl.offscreen, x, y, lbl.shape)
}

// SetString sets lbl's text
func (lbl *Label) SetString(s string) {
	lbl.clear(*lbl.gc)
	lbl.text = s
	lbl.reshape()
}

// GetString returns lbl's text
func (lbl *Label) GetString() string {
	return lbl.text
}

// SetInt does nothing
func (lbl *Label) SetInt(i int) {
}

// GetInt returns -1
func (lbl *Label) GetInt() int {
	return -1
}

// SetData does nothing
func (lbl *Label) SetData(d interface{}) {
}

// GetData returns nil
func (lbl *Label) GetData() interface{} {
	return nil
}

// SetEnabled does nothing
func (lbl *Label) SetEnabled(enabled bool) {
}

// GetEnabled returns true
func (lbl *Label) GetEnabled() bool {
	return true
}
