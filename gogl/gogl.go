// gogl
package gogl

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
	window  *glfw.Window
)

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

// Main runs the main SDL service loop.
// The binary's main.main must call sdl.Main() to run this loop.
// Main does not return. If the binary needs to do other work, it
// must do it in separate goroutines.
func Main() {
	for f := range mainfunc {
		f()
	}
}

// queue of work to run in main thread.
var mainfunc = make(chan func())

// do runs f on the main thread.
func do(f func()) {
	done := make(chan bool, 1)
	mainfunc <- func() {
		f()
		done <- true
	}
	<-done
}

/* And then other functions you write in package sdl can be like
func Beep() {
    do(func() {
        // whatever must run in main thread
    })
}*/

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func GlfwInit() {
	do(func() {
		glfw.SetErrorCallback(errorCallback)

		if !glfw.Init() {
			panic("Can't init glfw!")
		}

		glfw.WindowHint(glfw.Resizable, glfw.False)
		var err error
		window, err = glfw.CreateWindow(640, 480, "Testing", nil, nil)
		if err != nil {
			panic(err)
		}

		window.MakeContextCurrent()
		gl.Init()
	})
}

func GlfwDestroy() {
	if window != nil {
		window.Destroy()
	}
	glfw.Terminate()
}

func CompileShaders() {
	do(func() {
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

		vertex := gl.CreateShader(gl.VERTEX_SHADER)
		vertex.Source(vss)
		vertex.Compile()
		defer vertex.Delete()

		frag := gl.CreateShader(gl.FRAGMENT_SHADER)
		frag.Source(fss)
		frag.Compile()
		defer frag.Delete()

		program = gl.CreateProgram()
		program.AttachShader(vertex)
		program.AttachShader(frag)
		program.Link()

		program.Use()
	})
}

func Startup() {
	do(func() {
		log.Println(gl.GetString(gl.VERSION))

		vao = gl.GenVertexArray()
		vao.Bind()
	})
}

func Shutdown() {
	do(func() {
		vao.Delete()
		program.Delete()
	})
}

func Render() {
	do(func() {
		if !window.ShouldClose() {

			gl.ClearColor(1, 0, 0, 1)
			gl.Clear(gl.COLOR_BUFFER_BIT)

			gl.PointSize(40)
			gl.DrawArrays(gl.POINTS, 0, 1)

			window.SwapBuffers()
			glfw.PollEvents()
		}
	})
}
