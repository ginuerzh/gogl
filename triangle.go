// gogl project main.go
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func render(window *glfw.Window) {
	w, h := window.GetFramebufferSize()
	ratio := float64(w) / float64(h)
	gl.Viewport(0, 0, w, h)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(-ratio, ratio, -1, 1, 1, -1)
	gl.MatrixMode(gl.MODELVIEW)

	gl.LoadIdentity()
	gl.Rotatef(float32(glfw.GetTime())*50, 0, 0, 1)
	gl.Begin(gl.TRIANGLES)
	gl.Color3f(1, 0, 0)
	gl.Vertex3f(-0.6, -0.4, 0)
	gl.Color3f(0, 1, 0)
	gl.Vertex3f(0.6, -0.4, 0)
	gl.Color3f(0, 0, 1)
	gl.Vertex3f(0, 0.6, 0)
	gl.End()
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
		render(window)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
