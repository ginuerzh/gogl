// shader
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"log"
)

var (
	program gl.Program
	vao     gl.VertexArray
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func compileShaders() gl.Program {
	vss := `#version 130
	
			in vec4 position;
			in vec4 color;
			
			out vec4 vs_color;
	
			void main(void)
			{
				gl_Position = position;
				vs_color = color;
			}`

	fss := `#version 130
	
			in vec4 vs_color;
			
			out vec4 color;
			
			void main(void)
			{
				color = vs_color;
			}`

	vs := gl.CreateShader(gl.VERTEX_SHADER)
	vs.Source(vss)
	vs.Compile()
	defer vs.Delete()
	log.Println("vertex shader info log:", vs.GetInfoLog())

	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	fs.Source(fss)
	fs.Compile()
	defer fs.Delete()
	log.Println("frag shader info log:", vs.GetInfoLog())

	program := gl.CreateProgram()
	program.AttachShader(vs)
	program.AttachShader(fs)
	program.Link()

	program.Use()

	program.Validate()
	log.Println("validate status:", program.Get(gl.VALIDATE_STATUS))
	log.Println("program info log:", program.GetInfoLog())

	return program
}

func startup() {
	log.Println(gl.GetString(gl.VERSION))
	program = compileShaders()
	vao = gl.GenVertexArray()
	vao.Bind()

	vpos := [4]float32{-0.5, 0.5, 0.5, 1.0}
	vcolor := [4]float32{0.5, 0.8, 1.0, 1.0}

	pos := program.GetAttribLocation("position")
	color := program.GetAttribLocation("color")
	log.Println(pos, color)

	pos.Attrib4fv(&vpos)
	color.Attrib4fv(&vcolor)
}

func render() {

	gl.ClearColor(1, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.PointSize(40)

	gl.DrawArrays(gl.POINTS, 0, 1)
}

func shutdown() {
	vao.Delete()
	program.Delete()
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

	r := gl.Init()
	log.Println("init opengl:", r)

	startup()

	for !window.ShouldClose() {
		//Do OpenGL stuff
		render()

		window.SwapBuffers()
		glfw.PollEvents()
	}

	shutdown()
}
