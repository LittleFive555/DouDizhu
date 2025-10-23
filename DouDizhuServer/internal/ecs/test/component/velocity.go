package component

import (
	"github.com/ungerik/go3d/vec3"
)

type Velocity struct {
	Value *vec3.T
}

func (v *Velocity) IsComponent() {}
