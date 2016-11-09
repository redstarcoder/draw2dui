// Copyright (c) 2016, redstarcoder
// Open an OpenGl window and display widgets
package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
	"github.com/redstarcoder/draw2dui"
	"github.com/redstarcoder/draw2dui/widgets"
)

var (
	fps              int
	width, height    int
	redraw           = true
	font             draw2d.FontData
	gc               draw2d.GraphicContext
	offscreen        *glfw.Window
	widgetCollection *draw2dui.WidgetCollection
)

func setGlVars(w, h int) {
	gl.ClearColor(1, 1, 1, 1)
	/* Establish viewing area to cover entire window. */
	gl.Viewport(0, 0, int32(w), int32(h))
	/* PROJECTION Matrix mode. */
	gl.MatrixMode(gl.PROJECTION)
	/* Reset project matrix. */
	gl.LoadIdentity()
	/* Map abstract coords directly to window coords. */
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	/* Invert Y axis so increasing Y goes down. */
	gl.Scalef(1, -1, 1)
	/* Shift origin up to upper-left corner. */
	gl.Translatef(0, float32(-h), 0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	width, height = w, h
}

func reshape(window *glfw.Window, w, h int) {
	setGlVars(w, h)
	/* Recreate graphic context with new width & height. */
	gc = draw2dgl.NewGraphicContext(width, height)
	gc.SetFontData(draw2d.FontData{
		Name:   "luxi",
		Family: draw2d.FontFamilySerif,
		Style:  draw2d.FontStyleNormal /* | draw2d.FontStyleItalic*/})
	gc.SetFontSize(12)

	// Clear buffer
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Ask for redraw
	redraw = true

	// Have the widget collection handle the resize
	widgetCollection.Reshape(w, h)
}

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	width, height = 800, 800

	glfw.WindowHint(glfw.Visible, glfw.False)
	offscreen, err = glfw.CreateWindow(width, height, "", nil, nil)
	if err != nil {
		panic(err)
	}
	offscreen.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		panic(err)
	}
	setGlVars(width, height)

	glfw.WindowHint(glfw.Visible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "Show Widgets", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetSizeCallback(reshape)
	window.SetKeyCallback(onKey)
	window.SetCharCallback(onChar)
	window.SetCursorPosCallback(onMMove)
	window.SetMouseButtonCallback(onMClick)
	window.SetRefreshCallback(onRefresh)

	glfw.SwapInterval(0)

	gc = draw2dgl.NewGraphicContext(width, height)
	gc.SetFontData(draw2d.FontData{
		Name:   "luxi",
		Family: draw2d.FontFamilySerif,
		Style:  draw2d.FontStyleNormal /* | draw2d.FontStyleItalic*/})
	gc.SetFontSize(12)

	// Create widgets and widget collection
	textField := widgets.NewTextField(&gc, window, offscreen, 50, 50, 420, "Testing123456789", 75)
	button := widgets.NewButton(&gc, window, offscreen, 50, 50+gc.GetFontSize()+10, "O:")
	textBox := widgets.NewTextBox(&gc, window, offscreen, 50, 150, 420, 420, "Testing123456789\nTest2\n\n\n\nA very long line is here, it should automatically wrap because it is too long\n\n\n\n\n\n\n\n\n\n\n\n\n\ntest3\n\n\n\n\n\ntest4")
	textBox.InsertLine("INSERT LINE TEST")
	label := widgets.NewLabel(&gc, window, offscreen, 1, 5, "0 fps")
	widgetCollection = draw2dui.NewWidgetCollection(&gc, window, textField, button, label, textBox)

	reshape(window, width, height)
	lastUpdate := time.Now()
	lastDraw := lastUpdate
	drawDelta := time.Duration(0)
	tfps := 0

	drawWait := time.Duration(700/(glfw.GetPrimaryMonitor().GetVideoMode().RefreshRate)*1000) * time.Microsecond
	for !window.ShouldClose() {
		now := time.Now()
		if now.Sub(lastUpdate) >= time.Second {
			if fps != tfps {
				redraw = true
				fps = tfps
				label.SetString(fmt.Sprintf("%d fps", fps))
				if tfps == 0 {
					tfps = -1
				}
			}
			if tfps > 0 {
				tfps = 0
			}
			lastUpdate = now
		}
		if redraw && now.Sub(lastDraw) >= drawWait-drawDelta {
			widgetCollection.Draw()
			gl.Flush() /* Flush before buffer swap. */

			drawDelta = time.Since(lastDraw) - drawWait

			/* Swap buffers then copy the front buffer to the back */
			window.SwapBuffers()
			lastDraw = time.Now()
			gl.ReadBuffer(gl.FRONT)
			gl.DrawBuffer(gl.BACK)
			gl.RasterPos2i(0, int32(height))
			gl.CopyPixels(0, 0, int32(width), int32(height), gl.COLOR)

			/* Debugging code. Copy pixels instead of doing pointer swap, useful for confirming if
			   gl.CopyPixels is working as intended.

			gl.ReadBuffer(gl.BACK)
			gl.DrawBuffer(gl.FRONT)
			gl.RasterPos2i(0, int32(height))
			gl.CopyPixels(0, 0, int32(width), int32(height), gl.COLOR)

			lastDraw = time.Now()

			gl.ReadBuffer(gl.FRONT)
			gl.DrawBuffer(gl.BACK)*/

			tfps++
			redraw = false
		} else if !redraw {
			if fps < 4 {
				time.Sleep(time.Millisecond * 40)
			} else {
				time.Sleep(time.Millisecond * 8)
			}
		} else {
			time.Sleep(time.Millisecond)
		}
		glfw.PollEvents()
		if widgetCollection.Handle() {
			redraw = true
		}
	}
}

func onChar(w *glfw.Window, char rune) {
	_, event := widgetCollection.CharPress(char)
	if event != draw2dui.EventNone {
		redraw = true
	}
}

func onMMove(w *glfw.Window, xpos, ypos float64) {
	_, event := widgetCollection.MMove(xpos, ypos)
	if event != draw2dui.EventNone {
		redraw = true
	}
}

func onMClick(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	widget, event := widgetCollection.MClick(button, action, mods)
	if event != draw2dui.EventNone {
		if event == draw2dui.EventConfirm && widget.Name() == "Button-2" {
			log.Println("Click!")
		}
		redraw = true
	}
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	_, event := widgetCollection.KeyPress(key, action, mods)
	if event != draw2dui.EventNone {
		redraw = true
	}
	switch {
	case key == glfw.KeyEscape && action == glfw.Press:
		w.SetShouldClose(true)
	}
}

func onRefresh(w *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	redraw = true
	widgetCollection.Refresh()
}
