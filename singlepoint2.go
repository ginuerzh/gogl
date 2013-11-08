// singlepoint2
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"log"
	"runtime"
	"unsafe"
)

var (
	program gl.Program
	vao     gl.VertexArray
	buffer  gl.Buffer
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func compileShaders() gl.Program {
	vss := "#version 420 core\nlayout(location = 0) in vec4 position;\nvoid main(void)\n{\n		gl_Position = position;\n}"

	fss := "#version 420 core\nout vec4 color;\nvoid main(void)\n{\ncolor = vec4(0.0, 0.8, 1.0, 1.0);\n}"

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
	data := [4]float32{0.0, 0.0, 0.5, 1.0}
	program = compileShaders()

	vao = gl.GenVertexArray()
	vao.Bind()

	buffer = gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(data)), data, gl.STATIC_DRAW)
}

func render() {
	gl.ClearColor(1, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	program.Use()

	buffer.Bind(gl.ARRAY_BUFFER)
	var loc gl.AttribLocation = 0
	loc.AttribPointer(4, gl.FLOAT, false, 0, nil)
	loc.EnableArray()

	gl.DrawArrays(gl.POINTS, 0, 1)
	loc.DisableArray()
}

func shutdown() {
	buffer.Delete()
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
