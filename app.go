package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width              = 500
	height             = 500
	vertexShaderSource = `
#version 140
in vec3 vp;
in vec3 c;
out vec3 C;
uniform mat4 trans;
void main() {
    gl_Position = trans * vec4(vp, 1.0);
    C = c;
}
` + "\x00"
	fragmentShaderSource = `
#version 140
in vec3 C;
void main(void) {
    gl_FragColor = vec4 (C, 1.0);
}
` + "\x00"
)

type VAO struct {
	v uint32
	n int32
}

func mainloop(buildModel func() VAO, drawModel func(VAO, *glfw.Window, uint32), cbfun glfw.KeyCallback) {
	runtime.LockOSThread()

	window := initGlfw(cbfun)
	defer glfw.Terminate()

	program := initOpenGL()
	vao := buildModel()

	for !window.ShouldClose() {
		drawModel(vao, window, program)
	}
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw(cbfun glfw.KeyCallback) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	//glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Gorus", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetKeyCallback(cbfun)

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	gl.Enable(gl.CULL_FACE)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
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
