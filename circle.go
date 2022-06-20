//go:build ignore

package main

var Radius float
var Color vec4

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	alpha := 1.0
	center := vec2(Radius, Radius)
	if distance(texCoord, center) >= Radius {
		alpha = 0.0
	}
	return Color * alpha
}
