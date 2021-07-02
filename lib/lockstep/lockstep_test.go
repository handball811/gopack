package lockstep_test

import (
	"testing"

	"github.com/handball811/gopack/lib/lockstep"
	"github.com/stretchr/testify/assert"
)

const (
	id0 = 123
	id1 = 1234
	id2 = 15

	base0 = 500
	base1 = 5000
	base2 = 50000

	add0 = -3000
	add1 = 200
	add2 = -5600
)

func addNumStep(current *lockstep.Element, action *lockstep.Action) {
	current.Value = current.Value.(int) + action.Value.(int)
}

func cloneNum(from interface{}) interface{} {
	return from.(int)
}

// Add
func TestAdd(t *testing.T) {
	// setup
	var state interface{}
	var ok bool
	target := lockstep.NewLockStep(addNumStep, cloneNum)

	// when
	target.Add(id0, base0)
	target.Add(id1, base1)

	// then
	state, ok = target.LoadBase(id0)
	assert.Equal(t, base0, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id1)
	assert.Equal(t, base1, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id2)
	assert.Nil(t, state)
	assert.Equal(t, false, ok)
}

// Remove
func TestRemove(t *testing.T) {
	// setup
	var state interface{}
	var ok bool
	target := lockstep.NewLockStep(addNumStep, cloneNum)

	// when
	target.Add(id0, base0)
	target.Add(id1, base1)
	target.Add(id2, base2)
	target.Remove(id2)

	// then
	state, ok = target.LoadBase(id0)
	assert.Equal(t, base0, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id1)
	assert.Equal(t, base1, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id2)
	assert.Nil(t, state)
	assert.Equal(t, false, ok)
}

// Range
func TestRange(t *testing.T) {
	// setup
	baseMap := make(map[uint32]interface{})
	expectedMap := make(map[uint32]interface{})
	target := lockstep.NewLockStep(addNumStep, cloneNum)

	// when
	target.Add(id0, base0)
	target.Add(id1, base1)
	target.Add(id2, base2)
	target.Remove(id2)
	target.Range(func(id uint32, base, expected *lockstep.Element) bool {
		baseMap[id] = base.Value
		expectedMap[id] = expected.Value
		return true
	})

	// then
	assert.Equal(t, 2, len(baseMap))
	assert.Equal(t, 2, len(expectedMap))
	assert.Equal(t, base0, baseMap[id0])
	assert.Equal(t, base0, expectedMap[id0])
	assert.Equal(t, base1, baseMap[id1])
	assert.Equal(t, base1, expectedMap[id1])
}

// AddActions
func TestAddActions(t *testing.T) {
	// setup
	var state interface{}
	var ok bool
	target := lockstep.NewLockStep(addNumStep, cloneNum)

	// when
	target.Add(id0, base0)
	target.AddActions(map[uint32]interface{}{
		id0: add0,
		id1: add1,
		id2: add2,
	})
	target.Add(id1, base1)

	// then
	state, ok = target.LoadBase(id0)
	assert.Equal(t, base0, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadExpected(id0)
	assert.Equal(t, base0+add0, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id1)
	assert.Equal(t, base1, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadExpected(id1)
	assert.Equal(t, base1+add1, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id2)
	assert.Nil(t, state)
	assert.Equal(t, false, ok)

	state, ok = target.LoadExpected(id2)
	assert.Nil(t, state)
	assert.Equal(t, false, ok)
}

// AddActions+Complete
func TestAddActionsWithComplete(t *testing.T) {
	// setup
	var state interface{}
	var ok bool
	target := lockstep.NewLockStep(addNumStep, cloneNum)

	// when
	target.Add(id0, base0)
	index := target.AddActions(map[uint32]interface{}{
		id0: add0,
	})
	target.Complete(index)
	target.Add(id1, base1)

	// then
	state, ok = target.LoadBase(id0)
	assert.Equal(t, base0+add0, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadExpected(id0)
	assert.Equal(t, base0+add0, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id1)
	assert.Equal(t, base1, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadExpected(id1)
	assert.Equal(t, base1, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id2)
	assert.Nil(t, state)
	assert.Equal(t, false, ok)

	state, ok = target.LoadExpected(id2)
	assert.Nil(t, state)
	assert.Equal(t, false, ok)
}

// AddActions+Drop
func TestAddActionsWithDrop(t *testing.T) {
	// setup
	var state interface{}
	var ok bool
	target := lockstep.NewLockStep(addNumStep, cloneNum)

	// when
	target.Add(id0, base0)
	index := target.AddActions(map[uint32]interface{}{
		id0: add0,
	})
	target.Drop(index)
	target.Add(id1, base1)

	// then
	state, ok = target.LoadBase(id0)
	assert.Equal(t, base0, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadExpected(id0)
	assert.Equal(t, base0, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id1)
	assert.Equal(t, base1, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadExpected(id1)
	assert.Equal(t, base1, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id2)
	assert.Nil(t, state)
	assert.Equal(t, false, ok)

	state, ok = target.LoadExpected(id2)
	assert.Nil(t, state)
	assert.Equal(t, false, ok)
}

// AddAction+Complete+Drop
func TestAddActionsWithCompleteAndDrop(t *testing.T) {
	// setup
	var state interface{}
	var ok bool
	var index1, index2 int
	target := lockstep.NewLockStep(addNumStep, cloneNum)

	// when
	target.Add(id0, base0)
	target.AddActions(map[uint32]interface{}{
		id0: add0,
		id1: add1,
	})
	index1 = target.AddActions(map[uint32]interface{}{
		id0: add0,
		id1: add1,
	})
	index2 = target.AddActions(map[uint32]interface{}{
		id0: add0,
	})
	target.Drop(index1)
	target.Complete(index2)
	target.Add(id1, base1)

	// then
	state, ok = target.LoadBase(id0)
	assert.Equal(t, base0, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadExpected(id0)
	assert.Equal(t, base0+add0*2, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id1)
	assert.Equal(t, base1, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadExpected(id1)
	assert.Equal(t, base1+add1, state.(int))
	assert.Equal(t, true, ok)

	state, ok = target.LoadBase(id2)
	assert.Nil(t, state)
	assert.Equal(t, false, ok)

	state, ok = target.LoadExpected(id2)
	assert.Nil(t, state)
	assert.Equal(t, false, ok)
}
