// simpletexture.go
package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl42"
	//gogl "github.com/go-gl/gl"
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
		info := gl.GLStringAlloc(gl.Sizei(length))
		defer gl.GLStringFree(info)
		gl.GetShaderInfoLog(s, gl.Sizei(length), nil, info)
		return gl.GoString(info)
	}

	return ""
}

func getProgramInfoLog(program gl.Uint) string {
	var length gl.Int
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
	if length > 0 {
		info := gl.GLStringAlloc(gl.Sizei(length))
		defer gl.GLStringFree(info)
		gl.GetProgramInfoLog(program, gl.Sizei(length), nil, info)
		return gl.GoString(info)
	}

	return ""
}

func compileShaders() gl.Uint {
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
	/*
		vs := gogl.CreateShader(gogl.VERTEX_SHADER)
		vs.Source(vss)
		vs.Compile()
		defer vs.Delete()
		log.Println("vertex shader info log:", vs.GetInfoLog())

		fs := gogl.CreateShader(gogl.FRAGMENT_SHADER)
		fs.Source(fss)
		fs.Compile()
		defer fs.Delete()
		log.Println("frag shader info log:", fs.GetInfoLog())

		program := gogl.CreateProgram()
		program.AttachShader(vs)
		program.AttachShader(fs)
		program.Link()
		program.Use()

		program.Validate()
		log.Println("validate status:", program.Get(gl.VALIDATE_STATUS))
		log.Println("program info log:", program.GetInfoLog())
	*/
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
	/*
		size := 1024
		var length gl.Sizei = 0
		c := gl.GLStringAlloc(gl.Sizei(size))
		gl.GetShaderSource(fs, gl.Sizei(size), &length, c)
		log.Println(len(fss), length, gl.GoString(c))
		gl.GLStringFree(c)
	*/
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

	//Validation warning! - Sampler value s has not been set.
	log.Println("program info log:", getProgramInfoLog(program))

	return gl.Uint(program)
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
	log.Println(gl.GoStringUb(gl.GetString(gl.VERSION)))

	gl.ActiveTexture(gl.TEXTURE0 + 1)
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

	gl.UseProgram(program)

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
	glfw.SwapInterval(1)

	r := gl.Init()
	log.Println("init opengl:", r)

	//gogl.Init()

	startup()

	for !window.ShouldClose() {
		//Do OpenGL stuff
		render()

		window.SwapBuffers()
		glfw.PollEvents()
	}

	shutdown()
}
