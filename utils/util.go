package utils

import (
	gl "github.com/chsc/gogl/gl42"
)

func ToGLFloat(s []float32) []gl.Float {
	data := make([]gl.Float, len(s))
	for i, _ := range data {
		data[i] = gl.Float(s[i])
	}
	return data
}
