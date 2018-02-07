// main
package main

import (
	"openglex/shader"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 500
	height = 500
)

var triangle = []float32{
	-.5, -.5, 0, 0, .5, 0, //position and color
	.5, -.5, 0, .5, 0, 1,
	0, .5, 0, 1, .5, .5,
}
var vao uint32

//initiliaze the glfw and then create a window by givng hints
func initGlfw() *glfw.Window {

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Hello world", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window
}
func glInit() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	prog := gl.CreateProgram()

	//load the vectorshader and fragment shader
	vShader, err := shader.LoadShader("vectorshader.txt", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fShader, err := shader.LoadShader("fragmentShader.txt", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	//attach the shaders and link the program
	gl.AttachShader(prog, vShader)
	gl.AttachShader(prog, fShader)
	gl.LinkProgram(prog)

	gl.DeleteShader(vShader)
	gl.DeleteShader(fShader)

	//create a vertex array object and bind it
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	//create the vertext buffer memory on the GPU, bind it(make it current) and then fill in with our triangle vertices
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(triangle), gl.Ptr(triangle), gl.STATIC_DRAW)

	//tell the GL on how to pass the vertex buffer memory to the shader program : first tell about position vetex and enable it
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(24), gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	//then tell baout color vertex and enable it
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(24), gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)

	//work done now unbind the buffer and the vertex array. We bind the vertex array again when we need ot draw the vertex buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return prog
}

//draws the vertex buffer pointed by the vertex array so bind the vertex array that we need drawn
func draw(window *glfw.Window, program uint32) {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)

	gl.BindVertexArray(vao)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

	window.SwapBuffers()
	glfw.PollEvents()

}

func main() {

	//lock the current thread as opengl needs to work on the main thread
	runtime.LockOSThread()
	window := initGlfw()
	defer glfw.Terminate()
	prog := glInit()

	//continuously draw on the screen
	for !window.ShouldClose() {

		draw(window, prog)
	}
}
