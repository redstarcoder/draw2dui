// Copyright (c) 2016, redstarcoder
package widgets

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dbase"
	"github.com/llgcode/draw2d/draw2dgl"
	"golang.org/x/image/math/fixed"
)

// FIXME Do away with this function
// BUG(x) loadCurrentFont should take a GraphicContext interface
func loadCurrentFont(gc *draw2dgl.GraphicContext) (*truetype.Font, error) {
	font := draw2d.GetFont(gc.GetFontData())
	if font == nil {
		font = draw2d.GetFont(draw2dbase.DefaultFontData)
	}
	if font == nil {
		return nil, errors.New("No font set, and no default font available.")
	}
	gc.SetFont(font)
	gc.SetFontSize(gc.GetFontSize())
	return font, nil
}

func fUnitsToFloat64(x fixed.Int26_6) float64 {
	scaled := x << 2
	return float64(scaled/256) + float64(scaled%256)/256.0
}

// fillStringAtWidth draws the text at the specified point (x, y), stopping before width is exceeded
func fillStringAtWidth(_gc draw2d.GraphicContext, text string, x, y, width float64) float64 {
	gc := _gc.(*draw2dgl.GraphicContext)
	f, err := loadCurrentFont(gc)
	if err != nil {
		log.Println(err)
		return 0.0
	}
	startx := x
	prev, hasPrev := truetype.Index(0), false
	fontName := gc.GetFontName()
	for _, r := range text {
		index := f.Index(r)
		if hasPrev {
			x += fUnitsToFloat64(f.Kern(fixed.Int26_6(gc.Current.Scale), prev, index))
		}
		glyph := draw2dbase.FetchGlyph(gc, fontName, r)
		if x+glyph.Width-startx > width {
			break
		}
		x += glyph.Fill(gc, x, y)
		prev, hasPrev = index, true
	}
	return x - startx
}

// fillStringAtWidthCursor draws the text at the specified point (x, y), stopping before width is exceeded.
func fillStringAtWidthCursor(_gc draw2d.GraphicContext, c *Cursor, x, y, width float64) {
	gc := _gc.(*draw2dgl.GraphicContext)
	var last_i int = 1
	f, err := loadCurrentFont(gc)
	if err != nil {
		log.Println(err)
		return
	}
	startx := x
	prev, hasPrev := truetype.Index(0), false
	fontName := gc.GetFontName()
	cx := x
	cindex := c.i - c.iOffset
	text := c.text[c.iOffset:]
	for i, r := range text {
		index := f.Index(r)
		if hasPrev {
			x += fUnitsToFloat64(f.Kern(fixed.Int26_6(gc.Current.Scale), prev, index))
		}
		if i == cindex {
			cx = x + 1
		}
		glyph := draw2dbase.FetchGlyph(gc, fontName, r)
		if x+glyph.Width-startx > width {
			last_i = i
			break
		}
		x += glyph.Fill(gc, x, y)
		prev, hasPrev = index, true
	}
	if c.drawCursor {
		if cindex == len(text) {
			cx = x + 1
		}
		gl.LineWidth(2)
		gc.MoveTo(cx, y-gc.GetFontSize()-1)
		gc.LineTo(cx, y)
		gc.Stroke()
		gl.LineWidth(1)
	}
	if last_i > 1 {
		c.iEdge = last_i + c.iOffset
	} else {
		if c.iOffset > 0 {
			c.iEdge = len(c.text)
		} else {
			c.iEdge = len(c.text) + 100
		}
	}
}

// TODO calculate iEdge & iOffset using something like MoveToX
type Cursor struct {
	text              string   // text is the text stored in the field
	textLines         []string // textLines is the text stored in the field, stored as lines
	i, iOffset, iEdge int      // i is the position of the text cursor
	iY, maxLines      int      // iY is the y position of i, maxLines is the max visible lines
	lastBlink         time.Time
	drawCursor        bool
}

func (c *Cursor) GenLines(_gc draw2d.GraphicContext, width float64) {
	gc := _gc.(*draw2dgl.GraphicContext)
	f, err := loadCurrentFont(gc)
	if err != nil {
		log.Println(err)
		return
	}
	x := float64(3)
	ii := 0
	lastLine := 0
	prev, hasPrev := truetype.Index(0), false
	fontName := gc.GetFontName()
	c.textLines = make([]string, 0, 127)
	for i, r := range c.text {
		index := f.Index(r)
		if hasPrev {
			x += fUnitsToFloat64(f.Kern(fixed.Int26_6(gc.Current.Scale), prev, index))
		}
		glyph := draw2dbase.FetchGlyph(gc, fontName, r)
		if r == '\n' || x+glyph.Width > width {
			c.textLines = append(c.textLines, string(c.text[lastLine:lastLine+ii]))
			ii = 0
			x = 3
			if r != '\n' {
				lastLine = i
			} else {
				lastLine = i + 1
			}
			prev, hasPrev = truetype.Index(0), false
		}
		if r != '\n' {
			ii++
			x += glyph.Width
			prev, hasPrev = index, true
		}
	}
	if ii != 0 {
		c.textLines = append(c.textLines, string(c.text[lastLine:lastLine+ii]))
	}
}

func (c *Cursor) Blink() (redraw bool) {
	if c.i > c.iEdge {
		c.iOffset += c.i - c.iEdge
		c.iEdge = c.i
		redraw = true
	}
	if now := time.Now(); now.Sub(c.lastBlink) >= time.Millisecond*667 {
		c.lastBlink = now
		c.drawCursor = c.drawCursor == false
		redraw = true
	}
	return
}

func (c *Cursor) Insert(s string) {
	c.text = strings.Join([]string{c.text[:c.i], s, c.text[c.i:]}, "")
	c.MoveRight()
}

// GenLines must be called after
func (c *Cursor) InsertLine(s string) {
	c.text = strings.Join(c.textLines, "\n") + "\n" + s
}

func (c *Cursor) Backspace() bool {
	if c.i == 0 {
		return false
	}
	c.text = strings.Join([]string{c.text[:c.i-1], c.text[c.i:]}, "")
	if c.iOffset > 0 && len(c.text)+1 == c.iEdge {
		c.i--
		c.iEdge--
		c.iOffset--
		c.drawCursor = true
	} else {
		c.MoveLeft()
	}
	return true
}

// TODO Make MoveLeft and MoveRight take a parameter
func (c *Cursor) MoveLeft() bool {
	if c.i > 0 {
		c.i--
		if c.i < c.iOffset {
			c.iEdge--
			c.iOffset--
		}
		c.drawCursor = true
		return true
	}
	return false
}

func (c *Cursor) MoveRight() bool {
	if c.i < len(c.text) {
		c.i++
		if c.i > c.iEdge {
			c.iEdge++
			c.iOffset++
		}
		c.drawCursor = true
		return true
	}
	return false
}

func (c *Cursor) MoveTo(i int) bool {
	if i > c.i {
		for c.i < i {
			if !c.MoveRight() {
				break
			}
		}
		return true
	} else if i < c.i {
		for c.i > i {
			if !c.MoveLeft() {
				break
			}
		}
		return true
	}
	return false
}

func (c *Cursor) MoveToX(_gc draw2d.GraphicContext, x, mx, width float64) {
	gc := _gc.(*draw2dgl.GraphicContext)
	f, err := loadCurrentFont(gc)
	if err != nil {
		log.Println(err)
		return
	}
	prev, hasPrev := truetype.Index(0), false
	width += x
	fontName := gc.GetFontName()
	text := c.text[c.iOffset:]
	for i, r := range text {
		index := f.Index(r)
		if hasPrev {
			x += fUnitsToFloat64(f.Kern(fixed.Int26_6(gc.Current.Scale), prev, index))
		}
		glyph := draw2dbase.FetchGlyph(gc, fontName, r)
		if x+glyph.Width > mx || x+glyph.Width > width {
			c.i = i + c.iOffset
			return
		}
		x += glyph.Width
		prev, hasPrev = index, true
	}
	c.i = len(text) + c.iOffset
}
