package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

/*
  plots a flat disk with a hole in it
*/

const (
	base  = 0.5
	rings = 50
	steps = 50
)

var (
	model_pts = []float32{}
	model_col = []float32{}
	model_idx = []int32{}
	matrix    = mat4{
		vec4{1.0, 0.0, 0.0, 0.0},
		vec4{0.0, 1.0, 0.0, 0.0},
		vec4{0.0, 0.0, 1.0, 0.0},
		vec4{0.0, 0.0, 0.0, 1.0},
	}
)

func main() {
	generateModel()

	mainloop(func() VAO {
		return makeVao(model_pts, model_col, model_idx)
	}, draw, keyPress)
}

func draw(vao VAO, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.UniformMatrix4fv(0, 1, false, &(matrix[0][0]))

	gl.BindVertexArray(vao.v)
	gl.DrawElements(gl.TRIANGLES, vao.n, gl.UNSIGNED_INT, nil)

	glfw.PollEvents()
	window.SwapBuffers()
}

func keyPress(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	log.Println("Key Press: ", key, "(", action, ")")

	switch key {
	case 263:
		mult(&matrix, rotate(0, 0, float32(action)))
	case 262:
		mult(&matrix, rotate(0, 0, float32(360-action)))
	case 264:
		mult(&matrix, rotate(0, float32(action), 0))
	case 265:
		mult(&matrix, rotate(0, float32(360-action), 0))
	case 90:
		mult(&matrix, rotate(float32(action), 0, 0))
	case 88:
		mult(&matrix, rotate(float32(360-action), 0, 0))
	}
}

func generateModel() {
	const step = math.Pi * 2 / float64(steps)
	const step_r = math.Pi * 2 / float64(rings)

	rand.Seed(time.Now().Unix())
	colors := [rings][3]float32{}
	for i := 0; i < rings; i++ {
		colors[i][0] = rand.Float32()
		colors[i][1] = rand.Float32()
		colors[i][2] = rand.Float32()
	}

	for i := 0; i < steps; i++ {
		sin_R, cos_R := math.Sincos(step * float64(i))
		for j := 0; j < rings; j++ {
			sin_r, cos_r := math.Sincos(step_r * float64(j))
			r_x := cos_r * 0.2
			model_pts = append(
				model_pts,
				float32(cos_R*(base+r_x)),
				float32(sin_R*(base+r_x)),
				float32(sin_r*0.2))
			model_col = append(model_col, colors[j][0], colors[j][1], colors[j][2])
			if n := int32(len(model_pts)/3 - 1); n >= rings {
				if j > 0 {
					model_idx = append(model_idx, n, n-1, n-rings)
				}
				if j < rings-1 {
					model_idx = append(model_idx, n, n-rings, n-rings+1)
				}
			}
		}
	}
	// close the loop
	for j := int32(0); j < rings; j++ {
		if n := int32(len(model_pts)/3 - 1); n >= rings {
			if j > 0 {
				model_idx = append(model_idx, j, j-1, n-rings+1+j)
			}
			if j < rings-1 {
				model_idx = append(model_idx, j, n-rings+1+j, n-rings+2+j)
			}
		}
	}
	log.Println("Triangles: ", model_idx)
}
