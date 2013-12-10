// simpletexture2.go
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"log"
	"time"
)

var (
	program gl.Program
	texture gl.Texture
	vao     gl.VertexArray
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func compileShaders() gl.Program {
	vss := `#version 430 core

			void main(void)
			{
				const vec4 vertices[] = vec4[](
					vec4(0.75, -0.75, 0.5, 1.0),
					vec4(-0.75, -0.75, 0.5, 1.0),
					vec4( 0.75,  0.75, 0.5, 1.0));

				gl_Position = vertices[gl_VertexID];
			}`

	fss := `#version 430 core

			uniform sampler2D s;

			out vec4 color;

			void main(void)
			{
				color = texture(s, gl_FragCoord.xy / textureSize(s, 0));
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
	log.Println("frag shader info log:", fs.GetInfoLog())

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

func (width, height int) []float32 {
	data := make([]float32, width*height*4)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			data[(y*width+x)*4+0] = float32((x&y)&0xFF) / 255.0
			data[(y*width+x)*4+1] = float32((x|y)&0xFF) / 255.0
			data[(y*width+x)*4+2] = float32((x^y)&0xFF) / 255.0
			data[(y*width+x)*4+3] = 1.0
		}
	}

	return data
}

func startup() {

	//gl.ActiveTexture(gl.TEXTURE0 + 1)
	texture = gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)
	//gl.TexStorage2D(gl.TEXTURE_2D, gl.Sizei(8), gl.RGBA32F, gl.Sizei(256), gl.Sizei(256))

	data := genTexture(256, 256)
	gl.TexImage2D(gl.TEXTURE_2D, 8, gl.RGBA32F, 256, 256, 0, gl.RGBA, gl.FLOAT, data)
	//gl.TexSubImage2D(gl.TEXTURE_2D, gl.Int(0), gl.Int(0), gl.Int(0), gl.Sizei(256), gl.Sizei(256), gl.RGBA, gl.FLOAT, gl.Pointer(&data[0]))

	program = compileShaders()

	vao = gl.GenVertexArray()
	vao.Bind()
}

func shutdown() {
	program.Delete()
	vao.Delete()
	texture.Delete()
}

func render() {
	gl.ClearColor(0.0, 0.25, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	program.Use()

	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

func PrintGLParams() {
	log.Println("OpenGL:", gl.GetString(gl.VERSION))
	log.Println("GLSL:", gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	var max [1]int32
	gl.GetIntegerv(gl.MAX_UNIFORM_BLOCK_SIZE, max[:])
	log.Println("MAX_UNIFORM_BLOCK_SIZE", max[0])
	gl.GetIntegerv(gl.MAX_UNIFORM_BUFFER_BINDINGS, max[:])
	log.Println("MAX_UNIFORM_BUFFER_BINDINGS", max[0])
	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, max[:])
	log.Println("MAX_VERTEX_ATTRIBS", max[0])

	var units [1]int32
	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, units[:])
	log.Println("MAX_COMBINED_TEXTURE_IMAGE_UNITS", units[0])
}

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)

	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	r := gl.Init()
	log.Println("init opengl:", r)

	go func() {
		for {
			if e := gl.GetError(); e != gl.NO_ERROR {
				log.Println("Error detected:", e)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	PrintGLParams()

	startup()

	for !window.ShouldClose() {
		//Do OpenGL stuff
		render()

		window.SwapBuffers()
		glfw.PollEvents()
	}

	shutdown()
}
