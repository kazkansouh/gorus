package main

import (
	"math"
)

type vec4 [4]float32
type mat4 [4]vec4

func mult(a *mat4, b *mat4) {
	var r mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			accum := float32(0)
			for x := 0; x < 4; x++ {
				accum += a[i][x] * b[x][j]
			}
			r[i][j] = accum
		}
	}
	*a = r
}

func rotate(deg_z float32, deg_x float32, deg_y float32) *mat4 {
	sin_z, cos_z := math.Sincos(float64(deg_z) * math.Pi / 180)
	sin_x, cos_x := math.Sincos(float64(deg_x) * math.Pi / 180)
	sin_y, cos_y := math.Sincos(float64(deg_y) * math.Pi / 180)
	z := mat4{
		vec4{float32(cos_z), float32(-sin_z), 0, 0},
		vec4{float32(sin_z), float32(cos_z), 0, 0},
		vec4{0, 0, 1, 0},
		vec4{0, 0, 0, 1},
	}
	x := mat4{
		vec4{1, 0, 0, 0},
		vec4{0, float32(cos_x), float32(-sin_x), 0},
		vec4{0, float32(sin_x), float32(cos_x), 0},
		vec4{0, 0, 0, 1},
	}
	y := mat4{
		vec4{float32(cos_y), 0, float32(-sin_y), 0},
		vec4{0, 1, 0, 0},
		vec4{float32(sin_y), 0, float32(cos_y), 0},
		vec4{0, 0, 0, 1},
	}
	mult(&z, &x)
	mult(&z, &y)
	return &z
}
