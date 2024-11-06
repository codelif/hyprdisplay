package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

const GRID_SPACING = 2

func InitSurface() (tcell.Screen, *DisplaySurface) {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	surf := DisplaySurface{scr: s, pos: Position{x: 0, y: 0}}
	return s, &surf
}

func Application(config []DisplayConfig) {
	scr, surf := InitSurface()
	scr.EnableMouse()
	style := tcell.StyleDefault
	w, h := scr.Size()

	surf.pan_viewport(w/2, h/2)

	scr.SetStyle(style)
	scr.Clear()

	generate_displays(surf, config)

	surf.select_disp(0)
	for i, disp := range surf.displays {
		if disp.config.primary {
			surf.select_disp(i)
			break
		}
	}

	redraw_objects(surf, w, h)
mainloop:
	for {
		scr.Show()
		ev := scr.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			scr.Clear()
			scr.Sync()
			w, h = scr.Size()
			redraw_objects(surf, w, h)
		case *tcell.EventKey:
			mod, key, ch := ev.Modifiers(), ev.Key(), ev.Rune()
			if key == tcell.KeyEscape || ch == 'q' {
				break mainloop
			}
			if key == tcell.KeyCtrlN || (mod == tcell.ModCtrl && (key == tcell.KeyRight || key == tcell.KeyDown)) {
				surf.next_disp()
			} else if key == tcell.KeyCtrlB || (mod == tcell.ModCtrl && (key == tcell.KeyLeft || key == tcell.KeyUp)) {
				surf.prev_disp()
			}

			if mod == tcell.ModShift {
				camera_panned := true
				switch key {
				case tcell.KeyRight:
					surf.pan_viewport_rel(-2, 0)
				case tcell.KeyLeft:
					surf.pan_viewport_rel(2, 0)
				case tcell.KeyUp:
					surf.pan_viewport_rel(0, 1)
				case tcell.KeyDown:
					surf.pan_viewport_rel(0, -1)
				default:
					camera_panned = false
				}
				if camera_panned {
					scr.Clear()
				}
			}

			if (mod == tcell.ModNone) && (key == tcell.KeyRune) {
				switch ch {
				case '=':
					surf.inc_ppc()
				case '-':
					surf.dec_ppc()
				}
			}

			if (mod == tcell.ModNone) && (key == tcell.KeyEnter) {
				surf.snap_to_nearest()
			}

			surf.clear_disps()
			handle_move_box_input(surf.selected, mod, key, ch)
			redraw_objects(surf, w, h)
		}
	}
	scr.Fini()
	for _, line := range export(surf) {
		fmt.Println(line)
	}
	os.Exit(0)
}

func redraw_objects(surf *DisplaySurface, w, h int) {
	grid_axes(surf, w, h)
	// surf.clear_disps()
	surf.render()
	status_bar(surf, w, h)
}

func handle_move_box_input(disp *Display, mod tcell.ModMask, key tcell.Key, ch rune) {
	rune_mappings := func(ch rune) tcell.Key {
		var key tcell.Key
		switch ch {
		case 'k':
			key = tcell.KeyUp
		case 'j':
			key = tcell.KeyDown
		case 'h':
			key = tcell.KeyLeft
		case 'l':
			key = tcell.KeyRight
		}

		return key
	}

	mod_allowed := func(mod tcell.ModMask) bool {
		contains := false
		for _, v := range []tcell.ModMask{tcell.ModNone, tcell.ModAlt} {
			if mod == v {
				contains = true
			}
		}
		return contains
	}
	if !mod_allowed(mod) {
		return
	}

	if key == tcell.KeyRune {
		key = rune_mappings(ch)
	}

	delta_horizontal := 2
	delta_vertical := 1

	if mod == tcell.ModAlt {
		delta_horizontal = 1
		delta_vertical = 1
	}

	switch key {
	case tcell.KeyUp:
		disp.move_rel(0, -delta_vertical)
	case tcell.KeyDown:
		disp.move_rel(0, delta_vertical)
	case tcell.KeyLeft:
		disp.move_rel(-delta_horizontal, 0)
	case tcell.KeyRight:
		disp.move_rel(delta_horizontal, 0)
	}
}

func status_bar(surf *DisplaySurface, w, h int) {
	style := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	leftText := "Selected: " + surf.selected.name
	x, y, _, _ := surf.selected.calc_coords()
	rightText := fmt.Sprintf("Position: %v, %v", x, y)

	// disp2_x, disp2_y := disps[1].center()
	// centerText := fmt.Sprintf("Distance: %v", disps[0].distance(disp2_x, disp2_y))

	FillBox(surf.scr, 0, h-1, w-1, h-1, style, ' ')
	DrawHorizontal(surf.scr, 0, h-1, style, leftText)
	// DrawHorizontal(scr, (w-len(centerText)-1)/2, h-1, style, centerText)
	DrawHorizontal(surf.scr, w-len(rightText), h-1, style, rightText)
}

func grid_axes(surf *DisplaySurface, w, h int) {
	x_pad := surf.pos.x % (GRID_SPACING * 2)
	y_pad := surf.pos.y % GRID_SPACING
	for i := -surf.pos.x + x_pad; i < -surf.pos.x+w; i = i + (GRID_SPACING * 2) {
		for j := -surf.pos.y + y_pad; j < -surf.pos.y+h; j = j + GRID_SPACING {
			surf.SetContent(i, j, '·', nil, tcell.StyleDefault.Dim(true))
		}
	}
	DrawHLine(surf.scr, 0, surf.pos.y, tcell.StyleDefault.Dim(true), w)
	DrawVLine(surf.scr, surf.pos.x, 0, tcell.StyleDefault.Dim(true), h)
	surf.SetContent(0, 0, tcell.RunePlus, nil, tcell.StyleDefault.Dim(true))
}
