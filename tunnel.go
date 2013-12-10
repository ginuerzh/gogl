// simpletexture.go
package main

import (
	gl "github.com/chsc/gogl/gl42"
	"github.com/ginuerzh/gogl/utils"
	"github.com/ginuerzh/math3d"
)

const (
	width        = 640
	height       = 480
	title        = "OpenGL SuperBible - Tunnel"
	majorVersion = 3
	minorVersion = 0
	debug        = true
)

var (
	program gl.Uint
	vao     gl.Uint

	locMvp    gl.Int
	locOffset gl.Int

	texWall    gl.Uint
	texCeiling gl.Uint
	texFloor   gl.Uint
)

func startup() {

	program = utils.CompileShaders(utils.ShaderFile, "./shaders/tunnel_vs.glsl", "./shaders/tunnel_fs.glsl")

	locMvp = gl.GetUniformLocation(program, gl.GLString("mvp"))
	locOffset = gl.GetUniformLocation(program, gl.GLString("offset"))

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	texWall = utils.LoadKtx("./media/textures/brick.ktx", 0)
	texCeiling = utils.LoadKtx("./media/textures/ceiling.ktx", 0)
	texFloor = utils.LoadKtx("./media/textures/floor.ktx", 0)

	textures := []gl.Uint{texWall, texCeiling, texFloor}
	for _, v := range textures {
		gl.BindTexture(gl.TEXTURE_2D, v)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}
}

func shutdown() {
	gl.DeleteProgram(program)
	gl.DeleteVertexArrays(1, &vao)
}

func render(currentTime float64) {
	gl.Viewport(0, 0, width, height)

	if majorVersion > 3 || majorVersion == 3 && minorVersion >= 2 { // OpenGL version >= 3.2
		bgc := []gl.Float{0.0, 0.0, 0.0, 0.0}
		gl.ClearBufferfv(gl.COLOR, 0, &bgc[0])
	} else {
		gl.ClearColor(0, 0, 0, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	}

	gl.UseProgram(program)

	aspect := float64(width) / float64(height)
	proj_matrix := math3d.Perspective(60, aspect, 0.1, 100)

	gl.Uniform1f(locOffset, gl.Float(currentTime*0.003))

	textures := []gl.Uint{texWall, texFloor, texWall, texCeiling}
	for i, v := range textures {
		mv_matrix := math3d.RotateV(float64(90*i), math3d.NewVector3(0, 0, 1))
		mv_matrix = mv_matrix.MultiM(math3d.Translate(-0.5, 0, -10.0))
		mv_matrix = mv_matrix.MultiM(math3d.Rotate(90, 0, 1, 0))
		mv_matrix = mv_matrix.MultiM(math3d.Scale(30, 1, 1))

		data := utils.ToGLFloat(proj_matrix.MultiM(mv_matrix).ToSlice32())

		gl.UniformMatrix4fv(locMvp, 1, gl.GLBool(false), &data[0])

		gl.BindTexture(gl.TEXTURE_2D, v)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	}

}

func main() {
	utils.GlfwInit(width, height, title, majorVersion, minorVersion, debug)
	defer utils.GlfwDestroy()

	startup()
	defer shutdown()

	utils.GlfwMainLoop(render)
}
