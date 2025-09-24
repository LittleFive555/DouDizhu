package component

import "github.com/ungerik/go3d/vec3"

type Transform struct {
	Position vec3.T
	Rotation vec3.T
	Scale    vec3.T
}

func NewTransform() *Transform {
	return &Transform{
		Position: vec3.T{0, 0, 0},
		Rotation: vec3.T{0, 0, 0},
		Scale:    vec3.T{1, 1, 1},
	}
}
