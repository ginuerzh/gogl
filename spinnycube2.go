package main

import (
	gl "github.com/chsc/gogl/gl42"
	"github.com/ginuerzh/gogl/utils"
	"github.com/ginuerzh/math3d"
	"math"
)

const (
	width        = 640
	height       = 480
	title        = "OpenGL SuperBible - Spinny Cube"
	majorVersion = 3
	minorVersion = 0
	debug        = true
)

var (
	program  gl.Uint
	vao      gl.Uint
	vbuffer  gl.Uint
	ubuffer  gl.Uint
	mv_loc   gl.Int
	proj_loc gl.Int

	proj_matrix *math3d.Matrix4 = math3d.Perspective(50, 640.0/480.0, 0.1, 1000)
)

func startup() {
	vertex_positions := []float32{
		-0.25, 0.25, -0.25,
		-0.25, -0.25, -0.25,
		0.25, -0.25, -0.25,

		0.25, -0.25, -0.25,
		0.25, 0.25, -0.25,
		-0.25, 0.25, -0.25,

		0.25, -0.25, -0.25,
		0.25, -0.25, 0.25,
		0.25, 0.25, -0.25,

		0.25, -0.25, 0.25,
		0.25, 0.25, 0.25,
		0.25, 0.25, -0.25,

		0.25, -0.25, 0.25,
		-0.25, -0.25, 0.25,
		0.25, 0.25, 0.25,

		-0.25, -0.25, 0.25,
		-0.25, 0.25, 0.25,
		0.25, 0.25, 0.25,

		-0.25, -0.25, 0.25,
		-0.25, -0.25, -0.25,
		-0.25, 0.25, 0.25,

		-0.25, -0.25, -0.25,
		-0.25, 0.25, -0.25,
		-0.25, 0.25, 0.25,

		-0.25, -0.25, 0.25,
		0.25, -0.25, 0.25,
		0.25, -0.25, -0.25,

		0.25, -0.25, -0.25,
		-0.25, -0.25, -0.25,
		-0.25, -0.25, 0.25,

		-0.25, 0.25, -0.25,
		0.25, 0.25, -0.25,
		0.25, 0.25, 0.25,

		0.25, 0.25, 0.25,
		-0.25, 0.25, 0.25,
		-0.25, 0.25, -0.25,
	}
	size := len(vertex_positions)
	program = utils.CompileShaders(utils.ShaderFile, "./shaders/spinnycube_vs.glsl", "./shaders/spinnycube_fs.glsl")

	mv_loc = gl.GetUniformLocation(program, gl.GLString("mv_matrix"))
	proj_loc = gl.GetUniformLocation(program, gl.GLString("proj_matrix"))

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(size*4), gl.Pointer(&vertex_positions[0]), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.GLBool(false), 0, nil)
	gl.EnableVertexAttribArray(0)

	gl.Enable(gl.CULL_FACE)
	gl.FrontFace(gl.LEQUAL)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}

func render(currentTime float64) {
	if majorVersion > 3 || majorVersion == 3 && minorVersion >= 2 { // OpenGL version >= 3.2
		bgc := []gl.Float{0.0, 0.25, 0.0, 1.0}
		gl.ClearBufferfv(gl.COLOR, 0, &bgc[0])
	} else {
		gl.ClearColor(0, 0.25, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	}

	gl.ClearDepth(1.0)
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	data := utils.ToGLFloat(proj_matrix.ToSlice32())
	gl.UniformMatrix4fv(proj_loc, 1, gl.GLBool(false), &data[0])

	for i := 0; i < 24; i++ {
		f := float64(i) + currentTime*0.3

		mv_matrix := math3d.Translate(0, 0, -6)
		mv_matrix = mv_matrix.MultiM(math3d.Rotate(currentTime*45.0, 0, 1, 0))
		mv_matrix = mv_matrix.MultiM(math3d.Rotate(currentTime*21.0, 1.0, 0, 0))
		mv_matrix = mv_matrix.MultiM(math3d.Translate(
			math.Sin(2.1*f)*2.0,
			math.Cos(1.7*f)*2.0,
			math.Sin(1.3*f)*math.Cos(1.5*f)*2.0))

		data := utils.ToGLFloat(mv_matrix.ToSlice32())
		gl.UniformMatrix4fv(mv_loc, 1, gl.GLBool(false), &data[0])

		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}
}

func shutdown() {
	gl.DeleteBuffers(1, &vbuffer)
	gl.DeleteProgram(program)
	gl.DeleteVertexArrays(1, &vao)
}

func main() {
	utils.GlfwInit(width, height, title, majorVersion, minorVersion, debug)
	defer utils.GlfwDestroy()

	startup()
	defer shutdown()

	utils.GlfwMainLoop(render)
}
