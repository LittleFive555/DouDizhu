package test_test

import (
	"DouDizhuServer/internal/ecs/test"
	"testing"

	"github.com/ungerik/go3d/vec3"
)

func TestRunMoveSystem(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		runTimes int
		position *vec3.T
		velocity *vec3.T
		want     *vec3.T
	}{
		// TODO: Add test cases.
		{runTimes: 100, position: &vec3.T{0, 0, 0}, velocity: &vec3.T{1, 0, 0}, want: &vec3.T{100, 0, 0}},
		{runTimes: 100, position: &vec3.T{0, 0, 0}, velocity: &vec3.T{0, 1, 0}, want: &vec3.T{0, 100, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := test.RunMoveSystem(tt.runTimes, tt.position, tt.velocity)
			// TODO: update the condition below to compare got with tt.want.
			if *got != *tt.want {
				t.Errorf("Run() = %v, want %v", *got, *tt.want)
			}
		})
	}
}
