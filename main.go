package main

import (
	"log"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	model_pts = []float32{
		0, 0.5, -0.5,
		0, -0.5, -0.5,
		0.75, -0.5, -0.5,
		0.75, 0.5, 0.5,
		-0.75, 0.5, -0.5,
		-0.75, 0, -0.5,
		-0.25, 0.5, -0.5,
		-0.25, 0, -0.5,
		-0.75, 0.5, 0,
		-0.75, 0, 0,
		-0.25, 0.5, 0,
		-0.25, 0, 0,
	}
	model_col = []float32{
		1, 1, 0,
		0, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 0,
		0, 0, 1,
		0, 0, 1,
		1, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 1, 0,
		1, 0, 0,
	}
	model_idx = []int32{
		0, 1, 2,
		0, 2, 3,
		0, 3, 1,
		1, 3, 2,
		4, 5, 7,
		4, 7, 6,
		4, 9, 5,
		4, 8, 9,
		6, 7, 11,
		6, 7, 11,
		6, 11, 10,
		4, 6, 8,
		6, 10, 8,
		5, 9, 7,
		7, 9, 11,
		8, 10, 9,
		10, 11, 9,
	}
	matrix = mat4{
		vec4{1.0, 0.0, 0.0, 0.0},
		vec4{0.0, 1.0, 0.0, 0.0},
		vec4{0.0, 0.0, 1.0, 0.0},
		vec4{0.0, 0.0, 0.0, 1.0},
	}
)

func main() {
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

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32, col []float32, idx []int32) VAO {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo_idx uint32
	gl.GenBuffers(1, &vbo_idx)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vbo_idx)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(idx), gl.Ptr(idx), gl.STATIC_DRAW)

	var vbo_pts uint32
	gl.GenBuffers(1, &vbo_pts)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_pts)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	var vbo_col uint32
	gl.GenBuffers(1, &vbo_col)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_col)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(col), gl.Ptr(col), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

	return VAO{vao, int32(len(idx))}
}

func keyPress(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	log.Println("Key Press: ", key, "(", action, ")")

	switch key {
	case 263:
		mult(&matrix, rotate(0, 0, 1))
	case 262:
		mult(&matrix, rotate(0, 0, 359))
	case 264:
		mult(&matrix, rotate(0, 1, 0))
	case 265:
		mult(&matrix, rotate(0, 359, 0))
	case 90:
		mult(&matrix, rotate(1, 0, 0))
	case 88:
		mult(&matrix, rotate(359, 0, 0))
	}
}
