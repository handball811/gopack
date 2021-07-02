package lockstep_test

import (
	"testing"

	"github.com/handball811/gopack/lib/lockstep"
	"github.com/stretchr/testify/assert"
)

func TestConditionAddAndActions(t *testing.T) {
	var ok bool
	var n interface{}
	target := lockstep.NewConditionLockStep(func(from interface{}) interface{} {
		return from
	})

	// when
	target.Add(1, 1)
	target.Add(2, 5)
	index := target.AddActions(map[uint32]interface{}{
		1: 10,
		2: -3,
	})
	target.Complete(index)

	//then
	n, ok = target.LoadBase(1)
	assert.Equal(t, 10, n.(int))
	assert.Equal(t, true, ok)

	n, ok = target.LoadExpected(1)
	assert.Equal(t, 10, n.(int))
	assert.Equal(t, true, ok)

	n, ok = target.LoadBase(2)
	assert.Equal(t, -3, n.(int))
	assert.Equal(t, true, ok)

	n, ok = target.LoadExpected(2)
	assert.Equal(t, -3, n.(int))
	assert.Equal(t, true, ok)
}
