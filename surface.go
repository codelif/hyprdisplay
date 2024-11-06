package main

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type DisplaySurface struct {
	scr      tcell.Screen
	pos      Position
	displays []*Display
	selected *Display
	ppc      PixelPerCell
}

func (surf *DisplaySurface) pan_viewport(x, y int) {
	surf.pos.x, surf.pos.y = x, y
}

func (surf *DisplaySurface) pan_viewport_rel(x, y int) {
	surf.pos.x, surf.pos.y = surf.pos.x+x, surf.pos.y+y
}

func (surf *DisplaySurface) real(x, y int) (int, int) {
	return surf.pos.x + x, surf.pos.y + y
}

func (surf *DisplaySurface) render() {
	for _, disp := range surf.displays {
		if disp != surf.selected {
			disp.draw()
		}
	}
	surf.selected.draw()
}

func (surf *DisplaySurface) snap_to_nearest() {
	x, y := get_closest_edge(surf)
	if math.Abs(float64(x)/2) > math.Abs(float64(y)) {
		x = 0
	} else {
		y = 0
	}

	surf.clear_disps()
	surf.selected.move_rel(x, y)
}

func (surf *DisplaySurface) set_ppc(x, y float64) {
	print(fmt.Sprint(x, y))
	surf.clear_disps()
	x = min(148, x)
	x = max(8, x)

	y = min(296, y)
	y = max(16, y)

	surf.ppc = PixelPerCell{x: x, y: y}
}

func (surf *DisplaySurface) inc_ppc() {
	x := surf.ppc.x - 5.0
	y := x * 2
	surf.set_ppc(x, y)
}

func (surf *DisplaySurface) dec_ppc() {
	x := surf.ppc.x + 5.0
	y := x * 2
	surf.set_ppc(x, y)
}

func (surf *DisplaySurface) clear_disps() {
	for _, disp := range surf.displays {
		disp.clear()
	}
}

func (surf *DisplaySurface) select_disp(index int) error {
	if index >= len(surf.displays) || index < 0 {
		return errors.New("disps index is out of bounds")
	}
	if surf.selected != nil {
		surf.selected.style = surf.selected.style.Dim(true)
	}
	surf.selected = surf.displays[index]
	surf.selected.style = surf.selected.style.Dim(false)
	return nil
}

func (surf *DisplaySurface) next_disp() {
	for i, disp := range surf.displays {
		if disp == surf.selected {
			if i == len(surf.displays)-1 {
				i = -1
			}
			surf.select_disp(i + 1)
			break
		}
	}
}

func (surf *DisplaySurface) prev_disp() {
	for i, disp := range surf.displays {
		if disp == surf.selected {
			if i == 0 {
				i = len(surf.displays)
			}
			surf.select_disp(i - 1)
			break
		}
	}
}

func (surf *DisplaySurface) NewDisplay(x, y, width, height float64, name string, config DisplayConfig, style ...tcell.Style) *Display {
	new_style := tcell.StyleDefault
	if len(style) != 0 {
		new_style = style[len(style)-1]
	}
	disp := Display{surf: surf, x: x, y: y, width: width, height: height, style: new_style, name: name, config: config}

	surf.displays = append(surf.displays, &disp)
	return &disp
}

func (surf *DisplaySurface) SetContent(x, y int, primary rune, combining []rune, style tcell.Style) {
	x, y = surf.real(x, y)
	surf.scr.SetContent(x, y, primary, combining, style)
}

func (surf *DisplaySurface) DrawHorizontal(x, y int, style tcell.Style, text string) {
	for _, char := range text {
		surf.SetContent(x, y, char, nil, style)
		x++
	}
}

func (surf *DisplaySurface) DrawVertical(x, y int, style tcell.Style, text string) {
	for _, char := range text {
		surf.SetContent(x, y, char, nil, style)
		y++
	}
}

func (surf *DisplaySurface) DrawVLine(x, y int, style tcell.Style, n int) {
	surf.DrawVertical(x, y, style, strings.Repeat(string(tcell.RuneVLine), n))
}

func (surf *DisplaySurface) DrawHLine(x, y int, style tcell.Style, n int) {
	surf.DrawHorizontal(x, y, style, strings.Repeat(string(tcell.RuneHLine), n))
}

func (surf *DisplaySurface) FillBox(x1, y1, x2, y2 int, style tcell.Style, c rune) {
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if x1 > x2 {
		x1, x2 = x2, x1
	}

	for y := y1; y <= y2; y++ {
		surf.DrawHorizontal(x1, y, style, strings.Repeat(string(c), x2-x1+1))
	}
}
