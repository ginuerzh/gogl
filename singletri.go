// singletri.go
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
	
			void main(void)
			{
				const vec4 vertices[3] = vec4[3](vec4( 0.25, -0.25, 0.5, 1.0),
				vec4(-0.25, -0.25, 0.5, 1.0),vec4( 0.25, 0.25, 0.5, 1.0));
				
				gl_Position = vertices[gl_VertexID];
			}`

	fss := `#version 130
	
			out vec4 color;
			
			void main(void)
			{
				color = vec4(0.0, 0.8, 1.0, 1.0);
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
	log.Println("program validate:", program.Get(gl.VALIDATE_STATUS))
	log.Println("program info log:", program.GetInfoLog())

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

	gl.PointSize(40)
	gl.DrawArrays(gl.POINTS, 0, 3)
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
