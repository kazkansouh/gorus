package main

import (
	"log"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

/*
      8                  9
       .-----------------.
      /                 /|
     /                 / |
   0/                1/  |
   .-----------------.   |
   |                 |   |
   |  .-----------.  |   |
   |  | 4/      5 |  |   |11
   |  | /   +     |  |   .
   |  |/6       7 |  |  /
   |  .-----------.  | /
   |                 |/
   .-----------------.
   2                 3


      1                  0
       .-----------------.
      /                 /|
     /                 / |
   9/                8/  |
   .-----------------.   |
   |                 |   |
   |  .-----------.  |   |
   |  | 13     12 |  |   |2
   |  | /   +     |  |   .
   |  |/15     14 |  |  /
   |  .-----------.  | /
   |                 |/
   .-----------------.
   11                10

*/

const (
	bunit = 0.5
	sunit = 0.25
)

var (
	model_pts = []float32{
		-bunit, bunit, -bunit,
		bunit, bunit, -bunit,
		-bunit, -bunit, -bunit,
		bunit, -bunit, -bunit,
		-sunit, sunit, -bunit,
		sunit, sunit, -bunit,
		-sunit, -sunit, -bunit,
		sunit, -sunit, -bunit,
		-bunit, bunit, bunit,
		bunit, bunit, bunit,
		-bunit, -bunit, bunit,
		bunit, -bunit, bunit,
		-sunit, sunit, bunit,
		sunit, sunit, bunit,
		-sunit, -sunit, bunit,
		sunit, -sunit, bunit,
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
		1, 1, 0,
		0, 1, 1,
	}
	model_idx = []int32{
		15, 14, 6,
		15, 6, 7,
		12, 4, 14,
		4, 6, 14,
		5, 13, 7,
		13, 15, 7,
		13, 5, 4,
		4, 12, 13,
		0, 4, 1,
		1, 4, 5,
		1, 5, 3,
		3, 5, 7,
		3, 7, 2,
		2, 7, 6,
		2, 6, 0,
		0, 6, 4,
		9, 8, 0,
		9, 0, 1,
		9, 1, 11,
		1, 3, 11,
		0, 8, 2,
		8, 10, 2,
		11, 3, 2,
		11, 2, 10,
		8, 9, 13,
		8, 13, 12,
		8, 12, 10,
		12, 14, 10,
		10, 14, 11,
		11, 14, 15,
		11, 15, 9,
		15, 13, 9,
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
