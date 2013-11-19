// simpletexture.go
package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl43"
	glfw "github.com/go-gl/glfw3"
	"log"
)

var (
	program gl.Uint
	texture gl.Uint
	vao     gl.Uint
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func getShaderInfoLog(s gl.Uint) string {
	var length gl.Int
	gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &length)
	if length > 0 {
		info := make([]byte, length)
		gl.GetShaderInfoLog(s, gl.Sizei(length), nil, gl.GLString(string(info)))
		return string(info)
	}

	return ""
}

func getProgramInfoLog(program gl.Uint) string {
	var length gl.Int
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
	if length > 0 {
		info := make([]byte, length)
		gl.GetProgramInfoLog(program, gl.Sizei(length), nil, gl.GLString(string(info)))
		return string(info)
	}

	return ""
}

func compileShaders() gl.Uint {
	vss := `#version 430
	
			void main(void)
			{
				const vec4 vertices[] = vec4[](
					vec4(0.75, -0.75, 0.5, 1.0),
					vec4(-0.75, -0.75, 0.5, 1.0),
					vec4( 0.75,  0.75, 0.5, 1.0)
				);
				
				gl_Position = vertices[gl_VertexID];
			}`

	fss := `#version 430
	
			uniform sampler2D s;
			
			out vec4 color;
			
			void main(void)
			{
				color = textureFetch(s, ivec2(gl_FragCoord.xy), 0);
			}`

	vs := gl.CreateShader(gl.VERTEX_SHADER)
	cs := gl.GLString(vss)
	slen := gl.Int(len(vss))
	gl.ShaderSource(vs, 1, &cs, &slen)
	gl.CompileShader(vs)
	defer gl.DeleteShader(vs)
	log.Println("vertex shader info log:", getShaderInfoLog(vs))

	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	cs = gl.GLString(fss)
	slen = gl.Int(len(fss))
	gl.ShaderSource(fs, 1, &cs, &slen)
	gl.CompileShader(fs)
	defer gl.DeleteShader(fs)
	log.Println("frag shader info log:", getShaderInfoLog(fs))

	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)
	gl.UseProgram(program)

	gl.ValidateProgram(program)

	var status gl.Int
	gl.GetProgramiv(program, gl.VALIDATE_STATUS, &status)
	log.Println("validate status:", status)
	log.Println("program info log:", getProgramInfoLog(program))

	return program
}

func genTexture(width, height int) []float32 {
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
	var ver *gl.Ubyte = gl.GetString(gl.VERSION)

	log.Println(gl.GoStringUb(ver))

	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexStorage2D(gl.TEXTURE_2D, gl.Sizei(8), gl.RGBA32F, gl.Sizei(256), gl.Sizei(256))

	data := genTexture(256, 256)
	gl.TexSubImage2D(gl.TEXTURE_2D, gl.Int(0), gl.Int(0), gl.Int(0), gl.Sizei(256), gl.Sizei(256), gl.RGBA, gl.FLOAT, gl.Pointer(&data[0]))

	program = compileShaders()

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
}

func shutdown() {
	gl.DeleteProgram(program)
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteTextures(1, &texture)
}

func render() {
	green := []gl.Float{0.0, 0.25, 0.0, 1.0}
	gl.ClearBufferfv(gl.COLOR, 0, &green[0])

	gl.DrawArrays(gl.TRIANGLES, 0, 3)
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
