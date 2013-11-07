// shader
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"log"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func compileShaders() gl.Program {
	vss := `#version 130 core
			
			void main(void)
			{
				gl_Position = vec4(0.0, 0.0, 0.5, 1.0);
			}`
	fss := `#version 130 core
			out vec4 color;
			
			void main(void)
			{
				color = vec4(0.0, 0.8, 1.0, 1.0)
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

}

func render() {
	program := compileShaders()
	defer program.Delete()

	vao := gl.GenVertexArray()
	vao.Bind()
	defer vao.Delete()

	gl.ClearColor(1, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	program.Use()
	gl.DrawArrays(gl.POINTS, 0, 1)
}

func main() {
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

	for !window.ShouldClose() {
		//Do OpenGL stuff
		render()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
