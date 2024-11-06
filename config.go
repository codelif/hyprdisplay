package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type DisplayConfig struct {
	name    string
	pos     Position
	res     Resolution
	sres    ScaledResolution
	scale   float64
	primary bool
}

func NewDisplayConfig(name string, pos Position, res Resolution, scale float64, primary bool) DisplayConfig {
	sres := ScaledResolution{x: float64(res.x) / scale, y: float64(res.y) / scale}
	return DisplayConfig{name: name, pos: pos, res: res, sres: sres, scale: scale, primary: primary}
}

type Position struct {
	x, y int
}

type Resolution struct {
	x, y int
}

type ScaledResolution struct {
	x, y float64
}

type PixelPerCell struct {
	x, y float64
}

func generate_displays(surf *DisplaySurface, disp_confs []DisplayConfig) {
	ppc := PixelPerCell{x: 56, y: 112}
	surf.ppc = ppc
	for _, disp_conf := range disp_confs {
		x := float64(disp_conf.pos.x)
		y := float64(disp_conf.pos.y)
		width := disp_conf.sres.x
		height := disp_conf.sres.y
		style := tcell.StyleDefault
		if !disp_conf.primary {
			style = style.Dim(true)
		}
		surf.NewDisplay(x, y, width, height, disp_conf.name, disp_conf, style)
	}
}

func largest_disp(disps []DisplayConfig) DisplayConfig {
	largest_area := 0.0
	largest_disp := DisplayConfig{}
	for _, disp := range disps {
		area := disp.sres.x * disp.sres.y

		if area > largest_area {
			largest_area = area
			largest_disp = disp
		}
	}

	return largest_disp
}

func export(surf *DisplaySurface) []string {
	var primary *Display
	for _, disp := range surf.displays {
		if disp.config.primary {
			primary = disp
		}
	}

	// px, py, pw, ph := primary.calc_coords()

	template := "monitor = %s,preferred,%d,%d"
	var configs []string
	for _, disp := range surf.displays {
		name := disp.name
		x := disp.x - primary.x
		y := disp.y - primary.y
		// dx,dy,dw,dh := disp.calc_coords()
		configs = append(configs, fmt.Sprintf(template, name, int(x), int(y)))
	}

	return configs
}
