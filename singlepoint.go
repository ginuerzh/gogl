// shader
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"log"
	"runtime"
)

var (
	program gl.Program
	vao     gl.VertexArray
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func compileShaders() gl.Program {
	vss := `#version 420 core
			
			void main(void)
			{
				gl_Position = vec4(0.0, 0.0, 0.5, 1.0);
			}`
	fss := `#version 420 core
			out vec4 color;
			
			void main(void)
			{
				color = vec4(0.0, 0.8, 1.0, 1.0);
			}`

	vs := gl.CreateShader(gl.VERTEX_SHADER)
	log.Println("create shader:", vs)
	vs.Source(vss)
	vs.Compile()
	defer vs.Delete()

	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	fs.Source(fss)
	fs.Compile()
	defer fs.Delete()

	program := gl.CreateProgram()
	program.AttachShader(vs)
	program.AttachShader(fs)
	program.Link()

	return program
}

func startup() {
	log.Println(gl.GetString(gl.VERSION))

	program = compileShaders()
	vao = gl.GenVertexArray()
	vao.Bind()
}

func render() {

	gl.ClearColor(1, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	program.Use()

	gl.PointSize(40)

	gl.DrawArrays(gl.POINTS, 0, 1)
}

func shutdown() {
	vao.Delete()
	program.Delete()
}

func main() {
	runtime.LockOSThread()

	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()

	startup()

	for !window.ShouldClose() {
		//Do OpenGL stuff
		render()

		window.SwapBuffers()
		glfw.PollEvents()
	}

	shutdown()
}
