// singlepoint2
package main

import (
	gl "github.com/chsc/gogl/gl42"
	"github.com/ginuerzh/gogl/utils"
	"log"
)

const (
	width        = 640
	height       = 480
	title        = "OpenGL SuperBible - Shader Triangle"
	majorVersion = 3
	minorVersion = 0
	debug        = true
)

var (
	program gl.Uint
	vao     gl.Uint
	buffer  gl.Uint
	vss     = `#version 130

			uniform float offset = 0.5;
			in vec4 position;

			void main(void)
			{
				gl_Position = position;
			}`

	fss = `#version 130

			out vec4 color;

			void main(void)
			{
				color = vec4(0.0, 0.8, 1.0, 1.0);
			}`
)

func ptr2Slice(ptr gl.Pointer, size int) []float32 {
	return ((*[1 << 30]float32)(ptr))[0:size]
}

func startup() {
	data := []float32{
		0.25, -0.25, 0.5, 1.0,
		-0.25, -0.25, 0.5, 1.0,
		0.25, 0.25, 0.5, 1.0,
	}
	size := len(data)
	program = utils.CompileShaders(utils.ShaderString, vss, fss)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(size*4), nil, gl.STATIC_DRAW)

	ptr := gl.MapBuffer(gl.ARRAY_BUFFER, gl.WRITE_ONLY)

	n := copy(ptr2Slice(ptr, size), data)
	log.Println("copy data", n)
	gl.UnmapBuffer(gl.ARRAY_BUFFER)
}

func render(currentTime float64) {
	if majorVersion > 3 || majorVersion == 3 && minorVersion >= 2 { // OpenGL version >= 3.2
		bgc := []gl.Float{0.0, 0.25, 0.0, 1.0}
		gl.ClearBufferfv(gl.COLOR, 0, &bgc[0])
	} else {
		gl.ClearColor(0, 0.25, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	}

	gl.GenBuffers(1, &buffer)
	loc := gl.GetAttribLocation(program, gl.GLString("position"))

	gl.VertexAttribPointer(gl.Uint(loc), 4, gl.FLOAT, gl.GLBool(false), 0, nil)
	gl.EnableVertexAttribArray(gl.Uint(loc))

	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.DisableVertexAttribArray(gl.Uint(loc))
}

func shutdown() {
	gl.DeleteProgram(program)
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &buffer)
}

func main() {
	utils.GlfwInit(width, height, title, majorVersion, minorVersion, debug)
	defer utils.GlfwDestroy()

	startup()
	defer shutdown()

	utils.GlfwMainLoop(render)
}
