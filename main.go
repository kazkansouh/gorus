package main

import (
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

/*
  plots a torus

  r and R are the two radius of the torusm and the samples define the
  number of samples to be taken

   r
  <-->
     ------------
    /            \
   /              \
  /     /----\     \
  |    /      \    |
  |    |      |    |
  |    |      |    |
  |    \      /    |
  \     \----/     /
   \              /
    \            /
     -----------
          <------->
              R
*/

const (
	stride    = 9
	lines     = 3
	base_r    = 0.25
	base_R    = 0.6
	samples_r = stride * lines
	samples_R = stride * lines
	step_r    = math.Pi * 2 / float64(samples_r)
	step_R    = math.Pi * 2 / float64(samples_R)
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
	matrix_mutex sync.RWMutex
)

func main() {
	generateModel()

	fin := make(chan bool)

	go func(f chan<- bool) {
		mainloop(func() VAO {
			return makeVao(model_pts, model_col, model_idx)
		}, draw, keyPress)
		f <- true
	}(fin)

	t := time.Tick(time.Second / 16)

loop:
	for {
		select {
		case <-t:
			matrix_mutex.Lock()
			mult(&matrix, rotate(0, 0, 0.5))
			matrix_mutex.Unlock()
		case <-fin:
			log.Println("fin")
			break loop
		}
	}
}

func draw(vao VAO, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	matrix_mutex.RLock()
	gl.UniformMatrix4fv(0, 1, false, &(matrix[0][0]))

	gl.BindVertexArray(vao.v)
	gl.DrawElements(gl.TRIANGLES, vao.n, gl.UNSIGNED_INT, nil)
	matrix_mutex.RUnlock()

	glfw.PollEvents()
	window.SwapBuffers()
}

func keyPress(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	log.Println("Key Press: ", key, "(", action, ")")

	matrix_mutex.Lock()
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
	matrix_mutex.Unlock()
}

func generateModel() {
	// generate a random colour for each stratum of the torus
	rand.Seed(time.Now().Unix())
	colors := [samples_r][3]float32{}
	for i := 0; i < samples_r; i++ {
		colors[i][0] = rand.Float32()
		colors[i][1] = rand.Float32()
		colors[i][2] = rand.Float32()
	}

	// assemble the model
	for i := 0; i < samples_R; i++ {
		sin_R, cos_R := math.Sincos(step_R * float64(i))
		for j := 0; j < samples_r; j++ {
			sin_r, cos_r := math.Sincos(step_r * float64(j))
			r_x := cos_r * base_r
			model_pts = append(
				model_pts,
				float32(cos_R*(base_R+r_x)),
				float32(sin_R*(base_R+r_x)),
				float32(sin_r*base_r))
			model_col = append(
				model_col,
				colors[j][0],
				colors[j][1],
				colors[j][2])
			if n := int32(len(model_pts)/3 - 1); n >= samples_r && ((j+i)%stride == 0 || j == samples_r-1) {
				if (j+i)%stride == 0 {
					if j > 0 {
						model_idx = append(
							model_idx, n, n-1, n-samples_r)
					}
					if j < samples_r-1 {
						model_idx = append(
							model_idx,
							n,
							n-samples_r,
							n-samples_r+1)
					} else {
						model_idx = append(
							model_idx,
							n,
							n-samples_r,
							n-samples_r*2+1)
					}
				} else {
					if (j+i)%stride == stride-1 {
						model_idx = append(
							model_idx,
							n,
							n-samples_r*2+1,
							n-samples_r+1)
					}
				}
			}
		}
	}
	// close the loop
	for j := int32(0); j < samples_r; j++ {
		if n := int32(len(model_pts)/3 - 1); n >= samples_r && (j%stride == 0 || j == samples_r-1) {
			if j > 0 && j%stride == 0 {
				model_idx = append(
					model_idx, j, j-1, n-samples_r+1+j)
			}
			if j < samples_r-1 {
				model_idx = append(
					model_idx,
					j, n-samples_r+1+j, n-samples_r+2+j)
			} else {
				model_idx = append(
					model_idx,
					j, n-samples_r+1, 0)
			}
		}
	}
	log.Println("Triangles: ", model_idx)
}
