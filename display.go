package main

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

type Display struct {
	surf   *DisplaySurface
	style  tcell.Style
	x      float64
	y      float64
	width  float64
	height float64
	name   string
	config DisplayConfig
}

func (disp *Display) draw(style ...tcell.Style) {
	applied_style := disp.style
	if len(style) != 0 {
		applied_style = style[len(style)-1]
	}
	x, y, width, height := disp.calc_coords()
	disp.surf.DrawHLine(x+1, y, applied_style, width-2)
	disp.surf.DrawHLine(x+1, y+height-1, applied_style, width-2)
	disp.surf.DrawVLine(x, y+1, applied_style, height-2)
	disp.surf.DrawVLine(x+width-1, y+1, applied_style, height-2)

	disp.surf.SetContent(x, y, tcell.RuneULCorner, nil, applied_style)
	disp.surf.SetContent(x+width-1, y, tcell.RuneURCorner, nil, applied_style)
	disp.surf.SetContent(x+width-1, y+height-1, tcell.RuneLRCorner, nil, applied_style)
	disp.surf.SetContent(x, y+height-1, tcell.RuneLLCorner, nil, applied_style)

	disp.surf.DrawHorizontal(x+(width-len(disp.name))/2, y+(height-1)/2, applied_style, disp.name)
	// FillBox(disp.scr, disp.x+1, disp.y+1, disp.x+disp.width-2, disp.y+disp.height-2, applied_style, '#')
}

func (disp *Display) calc_coords() (int, int, int, int) {
	x := disp.x / disp.surf.ppc.x
	y := disp.y / disp.surf.ppc.y
	width := disp.width / disp.surf.ppc.x
	height := disp.height / disp.surf.ppc.y
	return int(math.Floor(x)), int(math.Floor(y)), int(math.Floor(width)), int(math.Floor(height))
}

func (disp *Display) clear() {
	x, y, width, height := disp.calc_coords()
	disp.surf.FillBox(x, y, x+width-1, y+height-1, tcell.StyleDefault, ' ')
}

func (disp *Display) move(x, y int) {
	// disp.clear()
	disp.x, disp.y = disp.surf.ppc.x*float64(x), disp.surf.ppc.y*float64(y)
}

func (disp *Display) move_rel(x, y int) {
	// disp.clear()
	disp.x, disp.y = disp.x+(disp.surf.ppc.x*float64(x)), disp.y+(disp.surf.ppc.y*float64(y))
}

func (disp *Display) distance(pos [2]float64) float64 {
	c := disp.center()
	deltax_2 := math.Pow(pos[0]-c[0], 2)
	deltay_2 := math.Pow(pos[1]-c[1], 2)
	return math.Sqrt(deltax_2 + deltay_2)
}

func (disp *Display) center() [2]float64 {
	x := float64(disp.x) + (float64(disp.width-1) / 2)
	y := float64(disp.y) + (float64(disp.height-1) / 2)
	return [2]float64{x, y}
}
