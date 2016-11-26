// Copyright (c) 2016, redstarcoder
package draw2dui

import (
	"os"
	"runtime"
	"testing"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
)

var (
	width, height int
	gc            draw2d.GraphicContext
)

// TODO make dummy widgets, use them for tests

func getNewWidgetCollection() (window, offscreen *glfw.Window, wc *WidgetCollection) {
	width, height = 800, 600

	//glfw.WindowHint(glfw.Visible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "++WidgetCollectionTest++", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	reshape(width, height)

	//glfw.WindowHint(glfw.Visible, glfw.False)
	offscreen, err = glfw.CreateWindow(width, height, "Offscreen (You shouldnt see this)", nil, nil)
	if err != nil {
		panic(err)
	}
	offscreen.MakeContextCurrent()
	reshape(width, height)

	window.MakeContextCurrent()

	return window, offscreen, NewWidgetCollection(&gc, window)
}

func BenchmarkNewWidgetCollection(b *testing.B) {
	width, height = 800, 600

	//glfw.WindowHint(glfw.Visible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "++WidgetCollectionTest++", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	reshape(width, height)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewWidgetCollection(&gc, window)
	}
	b.StopTimer()
	window.Destroy()
}

func TestWidgetCollectionDraw(t *testing.T) {
	window, offscreen, wc := getNewWidgetCollection()
	wc.forceRedraw = true
	wc.Draw()
	if wc.forceRedraw {
		t.Fail()
	}

	window.Destroy()
	offscreen.Destroy()
}

func TestWidgetCollectionHandle(t *testing.T) {
	window, offscreen, wc := getNewWidgetCollection()
	wc.forceRedraw = true

	if wc.Handle() {
		t.Error("Handle shouldn't request a redraw.")
	}

	window.Destroy()
	offscreen.Destroy()
}

func TestWidgetCollectionKeyPress(t *testing.T) {
	window, offscreen, wc := getNewWidgetCollection()
	wc.forceRedraw = true

	w, ev := wc.KeyPress(0, 0, 0)
	if w != nil || ev != EventNone {
		t.Fail()
	}

	window.Destroy()
	offscreen.Destroy()
}

func TestWidgetCollectionCharPress(t *testing.T) {
	window, offscreen, wc := getNewWidgetCollection()
	wc.forceRedraw = true

	w, ev := wc.CharPress(0)
	if w != nil || ev != EventNone {
		t.Fail()
	}

	window.Destroy()
	offscreen.Destroy()
}

func TestWidgetCollectionMMove(t *testing.T) {
	window, offscreen, wc := getNewWidgetCollection()
	wc.forceRedraw = true

	w, ev := wc.MMove(0, 0)
	if w != nil || ev != EventNone || !wc.hasCursor {
		t.Fail()
	}

	window.Destroy()
	offscreen.Destroy()
}

func TestWidgetCollectionMClick(t *testing.T) {
	window, offscreen, wc := getNewWidgetCollection()
	wc.forceRedraw = false

	w, ev := wc.MClick(0, glfw.Release, 0)
	if w != nil || ev != EventNone || wc.forceRedraw {
		t.Error("Failed Release test.")
	}

	wc.selected = "test"
	w, ev = wc.MClick(0, glfw.Press, 0)
	if w != nil || ev != EventSelected || !wc.forceRedraw || wc.selected != "" {
		t.Error("Failed Press test.")
	}

	window.Destroy()
	offscreen.Destroy()
}

func TestWidgetCollectionReshape(t *testing.T) {
	window, offscreen, wc := getNewWidgetCollection()
	wc.forceRedraw = false

	wc.Reshape(0, 0)
	if !wc.forceRedraw {
		t.Fail()
	}

	window.Destroy()
	offscreen.Destroy()
}

func TestWidgetCollectionRefresh(t *testing.T) {
	window, offscreen, wc := getNewWidgetCollection()
	wc.forceRedraw = false

	wc.Refresh()
	if !wc.forceRedraw {
		t.Fail()
	}

	window.Destroy()
	offscreen.Destroy()
}

func TestNameWidget(t *testing.T) {
	widgetCount = 0
	if NameWidget("test") != "test-1" {
		t.Fail()
	}
}

func TestMain(m *testing.M) {
	r := m.Run()
	glfw.Terminate()
	os.Exit(r)
}

func reshape(w, h int) {
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
	/* Recreate graphic context with new width & height. */
	gc = draw2dgl.NewGraphicContext(width, height)
}

func init() {
	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	err = gl.Init()
	if err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Visible, glfw.False)
}
