package main

import (
	"math"
)

const MAX_INT = int(^uint(0) >> 1)

func prepare_coords(surf *DisplaySurface) ([][]int, [][]int, [][]int, [][]int) {
	var x_a [][]int
	var x_b [][]int
	var y_a [][]int
	var y_b [][]int

	for i, disp := range surf.displays {
		x, y, w, h := disp.calc_coords()
		x_axis := []int{x, w - 1, i}
		y_axis := []int{y, h - 1, i}
		if disp == surf.selected {
			x_a = append(x_a, x_axis)
			y_a = append(y_a, y_axis)
			continue
		}
		x_b = append(x_b, x_axis)
		y_b = append(y_b, y_axis)
	}

	return x_a, x_b, y_a, y_b
}

func get_distances(p []int, s []int) (int, int, int, int) {
	pxsxd := math.Abs(float64(p[0] - s[0] + 1))
	pwsxd := math.Abs(float64(p[0] + p[1] - s[0] + 1))
	pxswd := math.Abs(float64(p[0] - s[0] - s[1] + 1))
	pwswd := math.Abs(float64(p[0] + p[1] - s[0] - s[1] + 1))

	return int(math.Floor(pxsxd)), int(math.Floor(pwsxd)), int(math.Floor(pxswd)), int(math.Floor(pwswd))
}

func closest_edge_delta(p []int, s []int) (int, int) {
	px, pw := p[0], p[0]+p[1]-1
	sx, sw := s[0], s[0]+s[1]-1

	pxsx, pwsx, pxsw, pwsw := get_distances(p, s)

	if px == sw || pw == sx {
		return 0, 0
	}

	if !(px < sw && pw > sx) {
		if pxsx > pwsx {
			if pwsx < pwsw {
				return sx - pw, pwsx
			} else {
				return sw - pw, pwsw
			}
		} else {
			if pxsx < pxsw {
				return sx - px, pxsx
			} else {
				return sw - px, pxsw
			}
		}
	} else {
		if px < sx {
			return sx - pw, pwsx
		} else {
			return sw - px, pxsw
		}
	}
}

func get_closest_edge_by_axis(p []int, axis [][]int, displays []*Display, a int) int {
	delta := MAX_INT
	delta_abs := MAX_INT
	for _, x := range axis {
		disp := displays[x[2]]
		py, ph := p[1], p[3]
		_, y, _, h := disp.calc_coords()
		if a == 1 {
			py, ph = p[0], p[2]
			y, _, h, _ = disp.calc_coords()
		}

		temp_delta, temp_delta_abs := closest_edge_delta([]int{py, ph}, []int{y, h})
		if temp_delta_abs < delta_abs {
			delta = temp_delta
			delta_abs = temp_delta_abs
		}
	}
	return delta
}

func get_closest_edge(surf *DisplaySurface) (int, int) {
	x_a, x_b, y_a, y_b := prepare_coords(surf)
	x_axis := get_intersections(x_a, x_b)
	y_axis := get_intersections(y_a, y_b)

	px, py, pw, ph := surf.selected.calc_coords()

	delta_y := get_closest_edge_by_axis([]int{px, py, pw, ph}, x_axis, surf.displays, 0)
	delta_x := get_closest_edge_by_axis([]int{px, py, pw, ph}, y_axis, surf.displays, 1)

	if delta_x == MAX_INT && delta_y == MAX_INT {
		delta_x, delta_y = 0, 0
	}
	return delta_x, delta_y
}

func get_intersections(A [][]int, B [][]int) [][]int {
	var res [][]int
	a, b := 0, 0
	if len(A) != 0 && len(B) != 0 {
		for a < len(A) && b < len(B) {
			if A[a][0] <= B[b][0]+B[b][1] && A[a][0]+A[a][1] >= B[b][0] {
				res = append(res, B[b])
			}

			b++
		}
	}
	return res
}
