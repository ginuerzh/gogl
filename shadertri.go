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

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func compileShaders() gl.Program {
	vss := `#version 130
			
			in vec4 position;
			
			void main(void)
			{
				gl_Position = position;
			}`
	fss := `#version 130
	
			out vec4 color;
			
			void main(void)
			{
				color = vec4(0.0, 0.8, 1.0, 1.0);
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
	log.Println(loc)

	return program
}

func ptr2Slice(ptr unsafe.Pointer, size int) []float32 {
	return ((*[1 << 30]float32)(ptr))[0:size]
}

func startup() {
	data := []float32{
		0.25, -0.25, 0.5, 1.0,
		-0.25, -0.25, 0.5, 1.0,
		0.25, 0.25, 0.5, 1.0,
	}
	size := len(data)
	program = compileShaders()

	vao = gl.GenVertexArray()
	vao.Bind()

	buffer = gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, size*4, nil, gl.STATIC_DRAW)

	ptr := gl.MapBuffer(gl.ARRAY_BUFFER, gl.WRITE_ONLY)

	n := copy(ptr2Slice(ptr, size), data)
	log.Println("copy data", n)
	gl.UnmapBuffer(gl.ARRAY_BUFFER)

}

func render() {
	gl.ClearColor(1, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	buffer.Bind(gl.ARRAY_BUFFER)
	loc := program.GetAttribLocation("position")
	loc.AttribPointer(4, gl.FLOAT, false, 0, nil)
	loc.EnableArray()

	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	loc.DisableArray()
}

func shutdown() {
	buffer.Delete()
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

	gl.Init()

	startup()
	defer shutdown()

	for !window.ShouldClose() {
		//Do OpenGL stuff
		render()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
