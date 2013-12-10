// gogl
package utils

import (
	"fmt"
	gl "github.com/chsc/gogl/gl42"
	glfw "github.com/go-gl/glfw3"
	"io/ioutil"
	"log"
	"time"
)

const (
	ShaderString = 0
	ShaderFile   = 1
)

var (
	window *glfw.Window
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func getShaderInfoLog(s gl.Uint) string {
	var length gl.Int
	var status gl.Int

	gl.GetShaderiv(s, gl.COMPILE_STATUS, &status)
	log.Println("shader compile status:", status)

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
	var status gl.Int

	// GL_DELETE_STATUS, GL_LINK_STATUS, GL_INFO_LOG_LENGTH, GL_ATTACHED_SHADERS
	// GL_ACTIVE_ATTRIBUTES, GL_ACTIVE_UNIFORMS, GL_ACTIVE_UNIFORM_BLOCKS
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	log.Println("program link status:", status)

	gl.GetProgramiv(program, gl.ACTIVE_UNIFORMS, &status)
	log.Println("active uniforms:", status)

	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
	if length > 0 {
		info := gl.GLStringAlloc(gl.Sizei(length))
		defer gl.GLStringFree(info)
		gl.GetProgramInfoLog(program, gl.Sizei(length), nil, info)
		return gl.GoString(info)
	}

	return ""
}

func printGLParams() {
	log.Println("OpenGL:", gl.GoStringUb(gl.GetString(gl.VERSION)))
	log.Println("GLSL:", gl.GoStringUb(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))

	var v gl.Int
	gl.GetIntegerv(gl.MAX_UNIFORM_BLOCK_SIZE, &v)
	log.Println("MAX_UNIFORM_BLOCK_SIZE", v)
	gl.GetIntegerv(gl.MAX_UNIFORM_BUFFER_BINDINGS, &v)
	log.Println("MAX_UNIFORM_BUFFER_BINDINGS", v)
	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, &v)
	log.Println("MAX_VERTEX_ATTRIBS", v)

	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &v)
	log.Println("MAX_COMBINED_TEXTURE_IMAGE_UNITS", v)
}

func CompileShaders(shaderType int, vert, frag string) gl.Uint {
	vss := vert
	if shaderType == ShaderFile {
		vsf, err := ioutil.ReadFile(vert)
		if err != nil {
			log.Fatal(err)
		}
		vss = string(vsf)
	}

	fss := frag
	if shaderType == ShaderFile {
		fsf, err := ioutil.ReadFile(frag)
		if err != nil {
			log.Fatal(err)
		}
		fss = string(fsf)
	}

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

func GlfwInit(width, height int, title string, major, minor int, debug bool) {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, major)
	glfw.WindowHint(glfw.ContextVersionMinor, minor)
	if debug {
		glfw.WindowHint(glfw.OpenglDebugContext, 1)
	}

	if major > 3 || major == 3 && minor >= 2 { // OpenGL version >= 3.2
		glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	}

	var err error
	window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	r := gl.Init()
	log.Println("init opengl:", r)
	printGLParams()

	go func() {
		for {
			if e := gl.GetError(); e != gl.NO_ERROR {
				log.Println("Error detected:", e)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func GlfwMainLoop(render func(float64)) {
	for !window.ShouldClose() {
		//Do OpenGL stuff
		render(glfw.GetTime())

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func GlfwDestroy() {
	if window != nil {
		window.Destroy()
	}
	glfw.Terminate()
}
