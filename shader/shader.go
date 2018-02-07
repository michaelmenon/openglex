// shader
package shader

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func init() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
}

//loads vertex shader from a provided file and returns the shader id after compiling it
func LoadShader(filepath string, shaderType uint32) (uint32, error) {
	fileb, err := ioutil.ReadFile(filepath)
	if err != nil {
		return 1, err
	}

	str := string(fileb) + "\x00"

	shader := gl.CreateShader(shaderType)
	source, free := gl.Strs(str)
	free()
	gl.ShaderSource(shader, 1, source, nil)
	gl.CompileShader(shader)

	// Check shader status
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %s: %v", filepath, log)
	}
	return shader, nil
}
