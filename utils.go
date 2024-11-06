package main

import (
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func print(text string) {
	exec.Command("notify-send",
		"-h", "int:x-hyprnotify-font-size:25", text).Run()
}

func DrawHorizontal(scr tcell.Screen, x, y int, style tcell.Style, text string) {
	for _, char := range text {
		scr.SetContent(x, y, char, nil, style)
		x++
	}
}

func DrawVertical(scr tcell.Screen, x, y int, style tcell.Style, text string) {
	for _, char := range text {
		scr.SetContent(x, y, char, nil, style)
		y++
	}
}

func DrawVLine(scr tcell.Screen, x, y int, style tcell.Style, n int) {
	DrawVertical(scr, x, y, style, strings.Repeat(string(tcell.RuneVLine), n))
}

func DrawHLine(scr tcell.Screen, x, y int, style tcell.Style, n int) {
	DrawHorizontal(scr, x, y, style, strings.Repeat(string(tcell.RuneHLine), n))
}

func FillBox(scr tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, c rune) {
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if x1 > x2 {
		x1, x2 = x2, x1
	}

	for y := y1; y <= y2; y++ {
		DrawHorizontal(scr, x1, y, style, strings.Repeat(string(c), x2-x1+1))
	}
}
