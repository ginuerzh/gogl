// uniformblkstd
package main

import (
	"fmt"
	"github.com/ginuerzh/math3d"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"log"
	"math"
	"runtime"
	"unsafe"
)

var (
	program  gl.Program
	vao      gl.VertexArray
	vbuffer  gl.Buffer
	ubuffer  gl.Buffer
	mv_loc   gl.UniformLocation
	proj_loc gl.UniformLocation

	proj_matrix *math3d.Matrix4
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func resizeCallback(w *glfw.Window, width int, height int) {
	aspect := float64(width) / float64(height)
	log.Println("resizeCallback")
	w.SetSize(width, height)

	proj_matrix = math3d.Perspective(50, aspect, 0.1, 1000)
}

func compileShaders() gl.Program {
	vss := `#version 430

			layout(location=0) in vec4 position;

			out VS_OUT
			{
				vec4 color;
			} vs_out;

			uniform mat4 mv_matrix;
			uniform mat4 proj_matrix;

			void main(void)
			{
				gl_Position = proj_matrix * mv_matrix * position;
				vs_out.color = position * 2.0 + vec4(0.5, 0.5, 0.5, 0.0);
			}`
	fss := `#version 430

			out vec4 color;

			in VS_OUT
			{
				vec4 color;
			} fs_in;

			void main(void)
			{
				color = fs_in.color;
			}`

	vertex := gl.CreateShader(gl.VERTEX_SHADER)
	vertex.Source(vss)
	vertex.Compile()
	defer vertex.Delete()
	log.Println("vertex shader info log:", vertex.GetInfoLog())

	frag := gl.CreateShader(gl.FRAGMENT_SHADER)
	frag.Source(fss)
	frag.Compile()
	defer frag.Delete()
	log.Println("frag shader info log:", frag.GetInfoLog())

	program := gl.CreateProgram()
	program.AttachShader(vertex)
	program.AttachShader(frag)
	program.Link()

	program.Use()

	program.Validate()
	log.Println("validate status:", program.Get(gl.VALIDATE_STATUS))
	log.Println("program info log:", program.GetInfoLog())

	loc := program.GetAttribLocation("position")
	log.Println("position:", loc)

	max := make([]int32, 1)
	gl.GetIntegerv(gl.MAX_UNIFORM_BLOCK_SIZE, max)
	log.Println("MAX_UNIFORM_BLOCK_SIZE", max[0])
	gl.GetIntegerv(gl.MAX_UNIFORM_BUFFER_BINDINGS, max)
	log.Println("MAX_UNIFORM_BUFFER_BINDINGS", max[0])
	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, max)
	log.Println("MAX_VERTEX_ATTRIBS", max[0])

	return program
}

func ptr2Slice(ptr unsafe.Pointer, size int) []float32 {
	return ((*[1 << 30]float32)(ptr))[0:size]
}

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
	program = compileShaders()

	mv_loc = program.GetUniformLocation("mv_matrix")
	proj_loc = program.GetUniformLocation("proj_matrix")

	vao = gl.GenVertexArray()
	vao.Bind()

	vbuffer = gl.GenBuffer()
	vbuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, size*4, vertex_positions, gl.STATIC_DRAW)
	var vloc gl.AttribLocation = 0
	vloc.AttribPointer(3, gl.FLOAT, false, 0, nil)
	vloc.EnableArray()

	gl.Enable(gl.CULL_FACE)
	gl.FrontFace(gl.LEQUAL)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}

func render(currentTime float64) {

	gl.ClearColor(0.0, 0.25, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.ClearDepth(1.0)
	gl.Clear(gl.DEPTH_BUFFER_BIT)
	a := proj_matrix.ToArray32()
	proj_loc.UniformMatrix4f(false, &a)

	for i := 0; i < 24; i++ {
		f := float64(i) + currentTime*0.3

		mv_matrix := math3d.Translate(0, 0, -6)
		mv_matrix = mv_matrix.MultiM(math3d.Rotate(currentTime*45.0, 0, 1, 0))
		mv_matrix = mv_matrix.MultiM(math3d.Rotate(currentTime*21.0, 1.0, 0, 0))
		mv_matrix = mv_matrix.MultiM(math3d.Translate(
			math.Sin(2.1*f)*2.0,
			math.Cos(1.7*f)*2.0,
			math.Sin(1.3*f)*math.Cos(1.5*f)*2.0))

		a := mv_matrix.ToArray32()
		mv_loc.UniformMatrix4f(false, &a)

		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}
}

func shutdown() {
	vbuffer.Delete()
	vao.Delete()
	program.Delete()
}

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	//glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	proj_matrix = math3d.Perspective(50, 640.0/480.0, 0.1, 1000)

	window.SetSizeCallback(resizeCallback)
	window.MakeContextCurrent()

	gl.Init()

	startup()
	defer shutdown()

	for !window.ShouldClose() {
		//Do OpenGL stuff
		render(glfw.GetTime())

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
