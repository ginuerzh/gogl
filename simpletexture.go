// simpletexture.go
package main

import (
	gl "github.com/chsc/gogl/gl42"
	"github.com/ginuerzh/gogl/utils"
)

var (
	program gl.Uint
	texture gl.Uint
	vao     gl.Uint
)

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
	gl.ActiveTexture(gl.TEXTURE0 + 1)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexStorage2D(gl.TEXTURE_2D, gl.Sizei(8), gl.RGBA32F, gl.Sizei(256), gl.Sizei(256))

	data := genTexture(256, 256)
	gl.TexSubImage2D(gl.TEXTURE_2D, gl.Int(0), gl.Int(0), gl.Int(0), gl.Sizei(256), gl.Sizei(256), gl.RGBA, gl.FLOAT, gl.Pointer(&data[0]))

	program = utils.CompileShaders(utils.ShaderFile, "./shaders/simpletexture_vs.glsl", "./shaders/simpletexture_fs.glsl")

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
	utils.GlfwInit(640, 480, "OpenGL SuperBible - Simple Texturing", 3, 0)
	defer utils.GlfwDestroy()

	startup()
	defer shutdown()

	utils.GlfwMainLoop(render)
}
