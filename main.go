package main

import (
	"log"
	"math"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

/*
  disc
*/

const (
	bunit = 0.5
	sunit = 0.25
	steps = 5
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

	for i := 0; i < steps; i++ {
		sin, cos := math.Sincos(step * float64(i))
		model_pts = append(model_pts, float32(cos*sunit), float32(sin*sunit), 0)
		model_col = append(model_col, 1, 0, 0)
		if n := int32(len(model_pts)/3 - 1); n >= 2 {
			model_idx = append(model_idx, n, n-2, n-1)
		}
		model_pts = append(model_pts, float32(cos*bunit), float32(sin*bunit), 0)
		model_col = append(model_col, 0, 0, 1)
		if n := int32(len(model_pts)/3 - 1); n >= 2 {
			model_idx = append(model_idx, n, n-1, n-2)
		}
	}
	if n := int32(len(model_pts)/3 - 1); n >= 2 {
		model_idx = append(model_idx, 0, n-1, n, 1, 0, n)
	}
	log.Println("Triangles: ", model_idx)
}
