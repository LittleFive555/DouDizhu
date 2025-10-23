package component

import (
	"github.com/ungerik/go3d/vec3"
)

type Position struct {
	Value *vec3.T
}

func (p *Position) IsComponent() {}
